package services

import (
	"fmt"
	"math"
	"time"

	"github.com/NonsoAmadi10/p2p-analysis/bitcoin"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

type NodeMetrics struct {
	Difficulty       float64     `json:"difficulty"`
	Version          interface{} `json:"version"`
	Chain            string      `json:"chain"`
	Blocks           int32       `json:"no_of_blocks"`
	BestBlockHash    string      `json:"bestblockhash"`
	UserAgent        interface{} `json:"user_agent"`
	HashRate         float64     `json:"hash_rate"`
	BlockPropagation float64     `json:"block_propagation"`
}

func GetInfo() (*NodeMetrics, error) {

	client, err := bitcoin.Client()
	if err != nil {
		return nil, err
	}

	defer client.Shutdown()

	info, err := client.GetBlockChainInfo()

	if err != nil {
		return nil, fmt.Errorf("failed to fetch blockchain info: %w", err)
	}

	networkInfo, err := client.GetNetworkInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch network info: %w", err)
	}

	lastBlockHash, err := chainhash.NewHashFromStr(info.BestBlockHash)
	if err != nil {
		return nil, fmt.Errorf("failed to parse best block hash: %w", err)
	}

	lastBlock, err := client.GetBlock(lastBlockHash)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch best block: %w", err)
	}

	prevBlockHeader, err := client.GetBlockHeader(&lastBlock.Header.PrevBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch previous block header: %w", err)
	}

	timeToFindBlock := lastBlock.Header.Timestamp.Sub(prevBlockHeader.Timestamp).Seconds()
	hashrate := 0.0
	if timeToFindBlock > 0 {
		hashrate = float64(info.Difficulty) * math.Pow(2, 32) / timeToFindBlock
	}

	blockCount, err := client.GetBlockCount()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch block count: %w", err)
	}

	blockHash, err := client.GetBlockHash(blockCount)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch block hash: %w", err)
	}

	blockHeader, err := client.GetBlockHeader(blockHash)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch block header: %w", err)
	}

	propagationTime := time.Since(blockHeader.Timestamp).Minutes()

	metrics := &NodeMetrics{
		Difficulty:       info.Difficulty,
		Version:          networkInfo.Version,
		Chain:            info.Chain,
		Blocks:           info.Blocks,
		BestBlockHash:    info.BestBlockHash,
		UserAgent:        networkInfo.SubVersion,
		HashRate:         hashrate,
		BlockPropagation: propagationTime,
	}

	return metrics, nil

}
