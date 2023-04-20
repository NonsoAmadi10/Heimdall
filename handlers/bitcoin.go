package handlers

import (
	"fmt"

	"github.com/NonsoAmadi10/p2p-analysis/bitcoin"
)

type NodeMetrics struct {
	Difficulty    float64     `json:"difficulty"`
	Version       int32       `json:"version"`
	Chain         string      `json:"chain"`
	Blocks        int32       `json:"no_of_blocks"`
	BestBlockHash string      `json:"bestblockhash"`
	UserAgent     interface{} `json:"user_agent"`
}

func GetInfo() *NodeMetrics {

	client := bitcoin.Client()

	defer client.Shutdown()

	info, _ := client.GetInfo()

	moreInfo, err := client.GetBlockChainInfo()

	if err != nil {
		fmt.Println(err)
	}

	networkInfo, _ := client.GetNetworkInfo()

	metrics := &NodeMetrics{
		Version:       info.Version,
		Difficulty:    info.Difficulty,
		Chain:         moreInfo.Chain,
		Blocks:        moreInfo.Blocks,
		BestBlockHash: moreInfo.BestBlockHash,
		UserAgent:     networkInfo.SubVersion,
	}

	return metrics
}
