package generic

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/client"
	libErr "github.com/centrifuge/go-substrate-rpc-client/v4/error"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

const (
	ErrGetBlockCall = libErr.Error("get block call")
)

//go:generate mockery --name Chain --structname ChainMock --filename chain_mock.go --inpackage

// Chain defines an interface for a client that can return a generic block B.
//
// The generic block B allows for multiple implementation of the block that can differ
// for each substrate chain implementation.
type Chain[
	A, S, P any,
	B GenericSignedBlock[A, S, P],
] interface {
	GetBlock(blockHash types.Hash) (B, error)
	GetBlockLatest() (B, error)
}

// genericChain implements the Chain interface.
type genericChain[
	A, S, P any,
	B GenericSignedBlock[A, S, P],
] struct {
	client client.Client
}

// NewChain creates a new instance of Chain.
func NewChain[
	A, S, P any,
	B GenericSignedBlock[A, S, P],
](client client.Client) Chain[A, S, P, B] {
	return &genericChain[A, S, P, B]{
		client: client,
	}
}

// DefaultChain is the Chain interface with defaults for the generic types:
//
// Address - types.MultiAddress
// Signature - types.MultiSignature
// PaymentFields - DefaultPaymentFields
// Block - *DefaultGenericSignedBlock
type DefaultChain = Chain[
	types.MultiAddress,
	types.MultiSignature,
	DefaultPaymentFields,
	*DefaultGenericSignedBlock,
]

// NewDefaultChain creates a new DefaultChain.
func NewDefaultChain(client client.Client) DefaultChain {
	return NewChain[
		types.MultiAddress,
		types.MultiSignature,
		DefaultPaymentFields,
		*DefaultGenericSignedBlock,
	](client)
}

// GetBlock retrieves a generic block B found at blockHash.
func (g *genericChain[A, S, P, B]) GetBlock(blockHash types.Hash) (B, error) {
	return g.getBlock(&blockHash)
}

// GetBlockLatest returns the latest generic block B.
func (g *genericChain[A, S, P, B]) GetBlockLatest() (B, error) {
	return g.getBlock(nil)
}

const (
	getBlockMethod = "chain_getBlock"
)

// getBlock retrieves the generic block B.
func (g *genericChain[A, S, P, B]) getBlock(blockHash *types.Hash) (B, error) {
	block := new(B)

	if err := client.CallWithBlockHash(g.client, block, getBlockMethod, blockHash); err != nil {
		return *block, ErrGetBlockCall.Wrap(err)
	}

	return *block, nil
}
