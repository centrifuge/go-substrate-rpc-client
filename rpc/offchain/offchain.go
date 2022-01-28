package offchain

import "github.com/ComposableFi/go-substrate-rpc-client/v4/client"

// Offchain exposes methods for retrieval of off-chain data
type Offchain struct {
	client client.Client
}

// NewOffchain creates a new Offchain struct
func NewOffchain(c client.Client) *Offchain {
	return &Offchain{client: c}
}
