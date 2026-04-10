package lightning

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/user"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gopkg.in/macaroon.v2"

	"github.com/lncm/lnd-rpc/v0.10.0/lnrpc"
)

type rpcCreds map[string]string

func (m rpcCreds) RequireTransportSecurity() bool { return true }
func (m rpcCreds) GetRequestMetadata(_ context.Context, _ ...string) (map[string]string, error) {
	return m, nil
}
func newCreds(bytes []byte) rpcCreds {
	creds := make(map[string]string)
	creds["macaroon"] = hex.EncodeToString(bytes)
	return creds
}

func getClient(hostname string, port int, tlsFile, macaroonFile string) (lnrpc.LightningClient, error) {
	macaroonBytes, err := os.ReadFile(macaroonFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read macaroon file: %w", err)
	}

	mac := &macaroon.Macaroon{}
	if err = mac.UnmarshalBinary(macaroonBytes); err != nil {
		return nil, fmt.Errorf("cannot unmarshal macaroon: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	transportCredentials, err := credentials.NewClientTLSFromFile(tlsFile, hostname)
	if err != nil {
		return nil, fmt.Errorf("cannot load tls credentials: %w", err)
	}

	fullHostname := fmt.Sprintf("%s:%d", hostname, port)

	connection, err := grpc.DialContext(ctx, fullHostname, []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(transportCredentials),
		grpc.WithPerRPCCredentials(newCreds(macaroonBytes)),
	}...)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to %s: %w", fullHostname, err)
	}

	return lnrpc.NewLightningClient(connection), nil
}

func Client() (lnrpc.LightningClient, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}
	homeDir := usr.HomeDir
	lndDir := fmt.Sprintf("%s/app_container/lightning", homeDir)
	var (
		hostname     = "localhost"
		port         = 10009
		tlsFile      = fmt.Sprintf("%s/tls.cert", lndDir)
		macaroonFile = fmt.Sprintf("%s/data/chain/bitcoin/testnet/admin.macaroon", lndDir)
	)

	client, err := getClient(hostname, port, tlsFile, macaroonFile)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to lnd RPC")
	return client, nil
}
