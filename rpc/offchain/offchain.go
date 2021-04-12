package offchain

import "github.com/centrifuge/go-substrate-rpc-client/v3/client"

// Offchain exposes methods for retrieval of off-chain data
type Offchain struct {
	client client.Client
}

// NewOffchain creates a new Offchain struct
func NewOffchain(c client.Client) *Offchain {
	return &Offchain{client: c}
}
