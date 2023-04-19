package lightning

import (
	"context"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
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

func getClient(hostname string, port int, tlsFile, macaroonFile string) lnrpc.LightningClient {
	macaroonBytes, err := ioutil.ReadFile(macaroonFile)
	if err != nil {
		panic(fmt.Sprintln("Cannot read macaroon file", err))
	}

	mac := &macaroon.Macaroon{}
	if err = mac.UnmarshalBinary(macaroonBytes); err != nil {
		panic(fmt.Sprintln("Cannot unmarshal macaroon", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	transportCredentials, err := credentials.NewClientTLSFromFile(tlsFile, hostname)
	if err != nil {
		panic(err)
	}

	fullHostname := fmt.Sprintf("%s:%d", hostname, port)

	connection, err := grpc.DialContext(ctx, fullHostname, []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(transportCredentials),
		grpc.WithPerRPCCredentials(newCreds(macaroonBytes)),
	}...)
	if err != nil {
		panic(fmt.Errorf("unable to connect to %s: %w", fullHostname, err))
	}

	return lnrpc.NewLightningClient(connection)
}

func Client() lnrpc.LightningClient {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	homeDir := usr.HomeDir
	lndDir := fmt.Sprintf("%s/app_container/lightning", homeDir)
	var (
		hostname     = "localhost"
		port         = 10009
		tlsFile      = fmt.Sprintf("%s/tls.cert", lndDir)
		macaroonFile = fmt.Sprintf("%s/data/chain/bitcoin/testnet/admin.macaroon", lndDir)
	)

	client := getClient(hostname, port, tlsFile, macaroonFile)

	return client
}
