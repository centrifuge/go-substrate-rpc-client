package rpc

import (
	"github.com/centrifuge/go-substrate-rpc-client/client"
	"github.com/centrifuge/go-substrate-rpc-client/rpc/chain"
	"github.com/centrifuge/go-substrate-rpc-client/rpc/state"
)

type RPC struct {
	Chain *chain.Chain
	State *state.State
	client *client.Client
}

func NewRPC(cl *client.Client) *RPC {
	return &RPC{
		Chain: chain.NewChain(cl),
		State: state.NewState(cl),
		client: cl,
	}
}
