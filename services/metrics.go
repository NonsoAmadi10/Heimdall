package services

import (
	"context"
	"log"
	"time"

	"github.com/NonsoAmadi10/p2p-analysis/bitcoin"
	"github.com/NonsoAmadi10/p2p-analysis/db"
	"github.com/NonsoAmadi10/p2p-analysis/lightning"
	"github.com/NonsoAmadi10/p2p-analysis/utils"
	"github.com/lncm/lnd-rpc/v0.10.0/lnrpc"
)

func ConnectionMetrics() {

	db := db.DB()

	// Get Bitcoin Client
	bitcoin := bitcoin.Client()

	defer bitcoin.Shutdown()

	// Get Lightning Client
	lnd := lightning.Client()

	for {
		//Get Bitcoin Peer Info

		peerInfo, err := bitcoin.GetPeerInfo()

		if err != nil {
			log.Printf("Failed to fetch btcd peer info: %v", err)
			continue
		}

		infoReq := &lnrpc.GetInfoRequest{}

		lndInfo, err := lnd.GetInfo(context.Background(), infoReq)

		if err != nil {
			log.Printf("Failed to fetch lnd info: %v", err)
			continue
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

		db.Create(&metrics)

		// Wait for 1 minute before fetching the next set of connection metrics
		time.Sleep(time.Minute * 3)
	}

}

func FetchMetrics() []utils.ConnectionMetrics {

	var allMetrics []utils.ConnectionMetrics

	db := db.DB()

	//fetch all metrics
	if err := db.Find(&allMetrics).Error; err != nil {
		log.Fatal(err)
	}

	return allMetrics

}
