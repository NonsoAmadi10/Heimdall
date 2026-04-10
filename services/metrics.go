package services

import (
	"context"
	"fmt"
	"log"
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
