package services

import (
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/NonsoAmadi10/p2p-analysis/bitcoin"
	"github.com/NonsoAmadi10/p2p-analysis/db"
	"github.com/NonsoAmadi10/p2p-analysis/lightning"
	"github.com/NonsoAmadi10/p2p-analysis/utils"
	"github.com/lncm/lnd-rpc/v0.10.0/lnrpc"
)

func ConnectionMetrics() {
	database, err := db.DB()
	if err != nil {
		log.Printf("Failed to initialize database for metrics collection: %v", err)
		return
	}

	btcClient, err := bitcoin.Client()
	if err != nil {
		log.Printf("Failed to initialize Bitcoin client for metrics collection: %v", err)
		return
	}
	defer btcClient.Shutdown()

	lndClient, err := lightning.Client()
	if err != nil {
		log.Printf("Failed to initialize Lightning client for metrics collection: %v", err)
		return
	}

	writeMetrics := func() error {
		peerInfo, err := btcClient.GetPeerInfo()

		if err != nil {
			return fmt.Errorf("failed to fetch btcd peer info: %w", err)
		}

		infoReq := &lnrpc.GetInfoRequest{}

		lndInfo, err := lndClient.GetInfo(context.Background(), infoReq)

		if err != nil {
			return fmt.Errorf("failed to fetch lnd info: %w", err)
		}

		// Calculate the incoming and outgoing bandwidth for the btcd node
		var btcdBandwidthIn, btcdBandwidthOut uint64
		for _, peer := range peerInfo {
			btcdBandwidthIn += peer.BytesRecv
			btcdBandwidthOut += peer.BytesSent
		}

		metrics := &utils.ConnectionMetrics{
			Timestamp:           time.Now(),
			NumBTCPeers:         int32(len(peerInfo)),
			NumLNDPeers:         int32(lndInfo.NumPeers),
			NumActiveChannels:   int32(lndInfo.NumActiveChannels),
			NumPendingChannels:  int32(lndInfo.NumPendingChannels),
			NumInactiveChannels: int32(lndInfo.NumInactiveChannels),
			BtcdBandwidthIn:     btcdBandwidthIn,
			BtcdBandwidthOut:    btcdBandwidthOut,
			BlockHeight:         int64(lndInfo.BlockHeight),
			BlockHash:           lndInfo.BlockHash,
			BestHeaderAge:       lndInfo.BestHeaderTimestamp,
			SyncedToChain:       lndInfo.SyncedToChain,
		}

		if err := database.Create(&metrics).Error; err != nil {
			return fmt.Errorf("failed to persist connection metrics: %w", err)
		}

		if err := EvaluateAlerts(database, metrics); err != nil {
			return fmt.Errorf("failed to evaluate alerts: %w", err)
		}

		return nil
	}

	if err := writeMetrics(); err != nil {
		log.Printf("Connection metrics collection failed: %v", err)
	}

	ticker := time.NewTicker(3 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		if err := writeMetrics(); err != nil {
			log.Printf("Connection metrics collection failed: %v", err)
		}
	}

}

func FetchMetrics() ([]utils.ConnectionMetrics, error) {

	var allMetrics []utils.ConnectionMetrics

	database, err := db.DB()
	if err != nil {
		return nil, err
	}

	//fetch all metrics
	if err := database.Find(&allMetrics).Error; err != nil {
		return nil, err
	}

	return allMetrics, nil

}

type MetricsAnalyticsPoint struct {
	BucketStart       time.Time `json:"bucket_start"`
	BucketEnd         time.Time `json:"bucket_end"`
	Samples           int       `json:"samples"`
	AvgBTCPeers       float64   `json:"avg_btc_peers"`
	AvgLNDPeers       float64   `json:"avg_lnd_peers"`
	AvgBandwidthIn    float64   `json:"avg_bandwidth_in"`
	AvgBandwidthOut   float64   `json:"avg_bandwidth_out"`
	MaxBandwidthIn    uint64    `json:"max_bandwidth_in"`
	MaxBandwidthOut   uint64    `json:"max_bandwidth_out"`
	SyncHealthPercent float64   `json:"sync_health_percent"`
}

type MetricsAnalyticsSummary struct {
	TotalSamples      int     `json:"total_samples"`
	AvgBTCPeers       float64 `json:"avg_btc_peers"`
	AvgLNDPeers       float64 `json:"avg_lnd_peers"`
	AvgBandwidthIn    float64 `json:"avg_bandwidth_in"`
	AvgBandwidthOut   float64 `json:"avg_bandwidth_out"`
	SyncHealthPercent float64 `json:"sync_health_percent"`
}

type MetricsAnalyticsResponse struct {
	From            time.Time               `json:"from"`
	To              time.Time               `json:"to"`
	IntervalMinutes int                     `json:"interval_minutes"`
	Points          []MetricsAnalyticsPoint `json:"points"`
	Summary         MetricsAnalyticsSummary `json:"summary"`
}

type analyticsAccumulator struct {
	samples         int
	syncSamples     int
	sumBTCPeers     int64
	sumLNDPeers     int64
	sumBandwidthIn  uint64
	sumBandwidthOut uint64
	maxBandwidthIn  uint64
	maxBandwidthOut uint64
}

func FetchMetricsAnalytics(from, to time.Time, interval time.Duration) (*MetricsAnalyticsResponse, error) {
	if !to.After(from) {
		return nil, fmt.Errorf("to must be after from")
	}

	if interval <= 0 {
		return nil, fmt.Errorf("interval must be greater than 0")
	}

	database, err := db.DB()
	if err != nil {
		return nil, err
	}

	var metrics []utils.ConnectionMetrics
	if err := database.Order("timestamp ASC").Find(&metrics).Error; err != nil {
		return nil, err
	}

	response := &MetricsAnalyticsResponse{
		From:            from.UTC(),
		To:              to.UTC(),
		IntervalMinutes: int(interval.Minutes()),
		Points:          []MetricsAnalyticsPoint{},
	}

	if len(metrics) == 0 {
		return response, nil
	}

	buckets := make(map[int64]*analyticsAccumulator)
	for _, m := range metrics {
		if m.Timestamp.Before(from) || m.Timestamp.After(to) {
			continue
		}

		bucketIndex := int64(m.Timestamp.Sub(from) / interval)
		acc, ok := buckets[bucketIndex]
		if !ok {
			acc = &analyticsAccumulator{}
			buckets[bucketIndex] = acc
		}

		acc.samples++
		acc.sumBTCPeers += int64(m.NumBTCPeers)
		acc.sumLNDPeers += int64(m.NumLNDPeers)
		acc.sumBandwidthIn += m.BtcdBandwidthIn
		acc.sumBandwidthOut += m.BtcdBandwidthOut
		if m.BtcdBandwidthIn > acc.maxBandwidthIn {
			acc.maxBandwidthIn = m.BtcdBandwidthIn
		}
		if m.BtcdBandwidthOut > acc.maxBandwidthOut {
			acc.maxBandwidthOut = m.BtcdBandwidthOut
		}
		if m.SyncedToChain {
			acc.syncSamples++
		}
	}

	bucketIndexes := make([]int64, 0, len(buckets))
	for i := range buckets {
		bucketIndexes = append(bucketIndexes, i)
	}
	sort.Slice(bucketIndexes, func(i, j int) bool {
		return bucketIndexes[i] < bucketIndexes[j]
	})

	summary := analyticsAccumulator{}
	for _, idx := range bucketIndexes {
		acc := buckets[idx]
		bucketStart := from.Add(time.Duration(idx) * interval).UTC()
		bucketEnd := bucketStart.Add(interval).UTC()

		point := MetricsAnalyticsPoint{
			BucketStart:       bucketStart,
			BucketEnd:         bucketEnd,
			Samples:           acc.samples,
			AvgBTCPeers:       float64(acc.sumBTCPeers) / float64(acc.samples),
			AvgLNDPeers:       float64(acc.sumLNDPeers) / float64(acc.samples),
			AvgBandwidthIn:    float64(acc.sumBandwidthIn) / float64(acc.samples),
			AvgBandwidthOut:   float64(acc.sumBandwidthOut) / float64(acc.samples),
			MaxBandwidthIn:    acc.maxBandwidthIn,
			MaxBandwidthOut:   acc.maxBandwidthOut,
			SyncHealthPercent: (float64(acc.syncSamples) / float64(acc.samples)) * 100,
		}

		response.Points = append(response.Points, point)

		summary.samples += acc.samples
		summary.syncSamples += acc.syncSamples
		summary.sumBTCPeers += acc.sumBTCPeers
		summary.sumLNDPeers += acc.sumLNDPeers
		summary.sumBandwidthIn += acc.sumBandwidthIn
		summary.sumBandwidthOut += acc.sumBandwidthOut
	}

	response.Summary = MetricsAnalyticsSummary{
		TotalSamples:      summary.samples,
		AvgBTCPeers:       float64(summary.sumBTCPeers) / float64(summary.samples),
		AvgLNDPeers:       float64(summary.sumLNDPeers) / float64(summary.samples),
		AvgBandwidthIn:    float64(summary.sumBandwidthIn) / float64(summary.samples),
		AvgBandwidthOut:   float64(summary.sumBandwidthOut) / float64(summary.samples),
		SyncHealthPercent: (float64(summary.syncSamples) / float64(summary.samples)) * 100,
	}

	return response, nil
}
