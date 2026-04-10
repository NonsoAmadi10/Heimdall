package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/NonsoAmadi10/p2p-analysis/bitcoin"
	"github.com/btcsuite/btcd/rpcclient"
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

type blockchainInfoResult struct {
	Difficulty    float64 `json:"difficulty"`
	Chain         string  `json:"chain"`
	Blocks        int64   `json:"blocks"`
	BestBlockHash string  `json:"bestblockhash"`
}

type networkInfoResult struct {
	Version    int64  `json:"version"`
	SubVersion string `json:"subversion"`
}

type blockHeaderVerboseResult struct {
	Time int64 `json:"time"`
}

func rawRequest[T any](client *rpcclient.Client, method string, params ...interface{}) (*T, error) {
	rawParams := make([]json.RawMessage, 0, len(params))
	for _, param := range params {
		raw, err := json.Marshal(param)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal %s params: %w", method, err)
		}
		rawParams = append(rawParams, raw)
	}

	rawResponse, err := client.RawRequest(method, rawParams)
	if err != nil {
		return nil, fmt.Errorf("failed to call %s: %w", method, err)
	}

	var result T
	if err := json.Unmarshal(rawResponse, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s response: %w", method, err)
	}

	return &result, nil
}

func GetInfo() (*NodeMetrics, error) {
	client, err := bitcoin.Client()
	if err != nil {
		return nil, err
	}

	defer client.Shutdown()

	info, err := rawRequest[blockchainInfoResult](client, "getblockchaininfo")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch blockchain info: %w", err)
	}

	networkInfo, err := rawRequest[networkInfoResult](client, "getnetworkinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch network info: %w", err)
	}

	hashRate, err := rawRequest[float64](client, "getnetworkhashps")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch network hash rate: %w", err)
	}

	blockHeader, err := rawRequest[blockHeaderVerboseResult](client, "getblockheader", info.BestBlockHash, true)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch block header: %w", err)
	}

	propagationTime := time.Since(time.Unix(blockHeader.Time, 0)).Minutes()

	metrics := &NodeMetrics{
		Difficulty:       info.Difficulty,
		Version:          networkInfo.Version,
		Chain:            info.Chain,
		Blocks:           int32(info.Blocks),
		BestBlockHash:    info.BestBlockHash,
		UserAgent:        networkInfo.SubVersion,
		HashRate:         *hashRate,
		BlockPropagation: propagationTime,
	}

	return metrics, nil

}
