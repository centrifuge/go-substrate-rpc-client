package ibc

import "github.com/ComposableFi/go-substrate-rpc-client/v4/client"

// IBC exposes methods for retrieval of chain data
type IBC struct {
	client client.Client
}

// NewIBC creates a new IBC struct
func NewIBC(cl client.Client) *IBC {
	return &IBC{cl}
}
