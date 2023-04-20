package services

import (
	"log"

	"github.com/NonsoAmadi10/p2p-analysis/bitcoin"
)

type NodeMetrics struct {
	Difficulty    float64     `json:"difficulty"`
	Version       interface{} `json:"version"`
	Chain         string      `json:"chain"`
	Blocks        int32       `json:"no_of_blocks"`
	BestBlockHash string      `json:"bestblockhash"`
	UserAgent     interface{} `json:"user_agent"`
}

func GetInfo() *NodeMetrics {

	client := bitcoin.Client()

	defer client.Shutdown()

	info, err := client.GetBlockChainInfo()

	if err != nil {
		log.Println(err)
	}

	networkInfo, _ := client.GetNetworkInfo()

	metrics := &NodeMetrics{
		Difficulty:    info.Difficulty,
		Version:       networkInfo.Version,
		Chain:         info.Chain,
		Blocks:        info.Blocks,
		BestBlockHash: info.BestBlockHash,
		UserAgent:     networkInfo.SubVersion,
	}

	return metrics

}
