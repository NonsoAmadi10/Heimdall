package utils

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type ConnectionMetrics struct {
	ID                  uint      `json:"id" gorm:"primary_key"`
	Timestamp           time.Time `json:"timestamp" gorm:"not null"`
	BlockHeight         int64     `json:"block_height"`
	BlockHash           string    `json:"block_hash"`
	BestHeaderAge       int64     `json:"best_header_age"`
	SyncedToChain       bool      `json:"synched_to_chain"`
	SyncProgress        float64   `json:"sync_progress"`
	NumBTCPeers         int32     `json:"num_of_btc_peers"`
	NumLNDPeers         int32     `json:"num_lnd_peers"`
	NumPendingChannels  int32     `json:"num_pending_channels"`
	NumActiveChannels   int32     `json:"num_active_channels"`
	NumInactiveChannels int32     `json:"num_inactive_channels"`
	BtcdBandwidthIn     uint64    `json:"btc_bandwidth_in"`
	BtcdBandwidthOut    uint64    `json:"btc_bandwidth_out"`
}

func GetEnv(key string) string {

	err := godotenv.Load()

	if err != nil {
		log.Println("Error loading .env file, will use system environment variables.")
	}

	return os.Getenv(key)
}
