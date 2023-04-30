package services

import (
	"log"
	"math"

	"github.com/NonsoAmadi10/p2p-analysis/bitcoin"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

type NodeMetrics struct {
	Difficulty    float64     `json:"difficulty"`
	Version       interface{} `json:"version"`
	Chain         string      `json:"chain"`
	Blocks        int32       `json:"no_of_blocks"`
	BestBlockHash string      `json:"bestblockhash"`
	UserAgent     interface{} `json:"user_agent"`
	HashRate      float64     `json:"hash_rate"`
}

func GetInfo() *NodeMetrics {

	client := bitcoin.Client()

	defer client.Shutdown()

	info, err := client.GetBlockChainInfo()

	if err != nil {
		log.Println(err)
	}

	networkInfo, _ := client.GetNetworkInfo()

	lastBlockHash, err := chainhash.NewHashFromStr(info.BestBlockHash)
	if err != nil {
		log.Println(err)
	}

	lastBlock, err := client.GetBlock(lastBlockHash)
	if err != nil {
		log.Println(err)
	}

	timeToFindBlock := lastBlock.Header.Timestamp.Unix() - int64(lastBlock.Header.PrevBlock[len(lastBlock.Header.PrevBlock)-1])
	hashrate := float64(info.Difficulty) / (float64(timeToFindBlock) * math.Pow(2, 32))

	metrics := &NodeMetrics{
		Difficulty:    info.Difficulty,
		Version:       networkInfo.Version,
		Chain:         info.Chain,
		Blocks:        info.Blocks,
		BestBlockHash: info.BestBlockHash,
		UserAgent:     networkInfo.SubVersion,
		HashRate:      hashrate,
	}

	return metrics

}
