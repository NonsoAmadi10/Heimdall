package services

import (
	"context"
	"fmt"

	"github.com/NonsoAmadi10/p2p-analysis/lightning"
	"github.com/lncm/lnd-rpc/v0.10.0/lnrpc"
)

type LNodeMetrics struct {
	PubKey      string `json:"pub_key"`
	UserAgent   string `json:"user_agent"`
	Alias       string `json:"alias"`
	NetCapacity int    `json:"network_capacity"`
}

func GetNodeInfo() (*LNodeMetrics, error) {

	client, err := lightning.Client()
	if err != nil {
		return nil, err
	}

	infoReq := &lnrpc.GetInfoRequest{}

	info, err := client.GetInfo(context.Background(), infoReq)

	if err != nil {
		return nil, fmt.Errorf("error getting node info: %w", err)
	}

	moreInfo, err := client.GetNetworkInfo(context.Background(), &lnrpc.NetworkInfoRequest{})
	if err != nil {
		return nil, fmt.Errorf("error getting network info: %w", err)
	}

	result := &LNodeMetrics{
		UserAgent:   info.Version,
		Alias:       info.Alias,
		NetCapacity: int(moreInfo.TotalNetworkCapacity),
		PubKey:      info.IdentityPubkey,
	}

	return result, nil
}
