// Go Substrate RPC client implements all the basic RPC calls to communicate with a Substrate node.
package gsrpc

import (
	"github.com/centrifuge/go-substrate-rpc-client/client"
	"github.com/centrifuge/go-substrate-rpc-client/rpc"
)

type SubstrateAPI struct {
	RPC    *rpc.RPC
	URL    string
	client *client.Client
}

func NewSubstrateAPI(url string) (*SubstrateAPI, error) {
	cl, err := client.Connect(url)
	if err != nil {
		return &SubstrateAPI{}, err
	}

	clPtr := &cl

	return &SubstrateAPI{
		RPC:    rpc.NewRPC(clPtr),
		URL:    url,
		client: clPtr,
	}, nil
}
