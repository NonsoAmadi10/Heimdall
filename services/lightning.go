package services

import (
	"context"
	"log"

	"github.com/NonsoAmadi10/p2p-analysis/lightning"
	"github.com/lncm/lnd-rpc/v0.10.0/lnrpc"
)

type LNodeMetrics struct {
	PubKey      string `json:"pub_key"`
	UserAgent   string `json:"user_agent"`
	Alias       string `json:"alias"`
	NetCapacity int    `json:"network_capacity"`
}

func GetNodeInfo() *LNodeMetrics {

	client := lightning.Client()

	infoReq := &lnrpc.GetInfoRequest{}

	info, err := client.GetInfo(context.Background(), infoReq)

	if err != nil {
		log.Fatalf("Error getting node info: %v", err)
	}

	moreInfo, _ := client.GetNetworkInfo(context.Background(), &lnrpc.NetworkInfoRequest{})

	result := &LNodeMetrics{
		UserAgent:   info.Version,
		Alias:       info.Alias,
		NetCapacity: int(moreInfo.TotalNetworkCapacity),
		PubKey:      info.IdentityPubkey,
	}

	return result
}
