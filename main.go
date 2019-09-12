// Go Substrate RPC client implements all the basic RPC calls to communicate with a Substrate node.
package gsrpc

import (
	"github.com/centrifuge/go-substrate-rpc-client/client"
	"github.com/centrifuge/go-substrate-rpc-client/rpc"
)


type SubstrateApi struct {
	RPC *rpc.RPC
	Url string
	client *client.Client
}

func NewSubstrateApi(url string) (*SubstrateApi, error) {
	cl, err := client.Connect(url)
	if err != nil {
		return &SubstrateApi{}, err
	}

	clPtr := &cl

	return &SubstrateApi{
		RPC: rpc.NewRPC(clPtr),
		Url: url,
		client: clPtr,
	}, nil
}
