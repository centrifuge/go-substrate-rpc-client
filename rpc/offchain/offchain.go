package offchain

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/client"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

type Offchain interface {
	LocalStorageGet(kind StorageKind, key []byte) (*types.StorageDataRaw, error)
	LocalStorageSet(kind StorageKind, key []byte, value []byte) error
}

// offchain exposes methods for retrieval of off-chain data
type offchain struct {
	client client.Client
}

// NewOffchain creates a new offchain struct
func NewOffchain(c client.Client) Offchain {
	return &offchain{client: c}
}
