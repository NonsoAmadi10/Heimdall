package bitcoin

import (
	"fmt"
	"log"

	"github.com/NonsoAmadi10/p2p-analysis/utils"
	"github.com/btcsuite/btcd/rpcclient"
)

func Client() (*rpcclient.Client, error) {

	// Connect to a running Bitcoin Core node via RPC
	connCfg := &rpcclient.ConnConfig{
		Host:         utils.GetEnv("BTC_HOST"),
		User:         utils.GetEnv("BTC_USER"),
		Pass:         utils.GetEnv("BTC_PASS"),
		HTTPPostMode: true,
		DisableTLS:   true,
	}

	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return nil, fmt.Errorf("error connecting to bitcoind: %w", err)
	}

	if _, err := client.GetBlockCount(); err != nil {
		client.Shutdown()
		return nil, fmt.Errorf("error validating bitcoind connectivity: %w", err)
	}

	log.Println("Connected to bitcoind RPC")
	return client, nil
}
