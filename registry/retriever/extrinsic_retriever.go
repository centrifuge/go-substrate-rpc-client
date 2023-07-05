package retriever

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/registry"
	"github.com/centrifuge/go-substrate-rpc-client/v4/registry/exec"
	"github.com/centrifuge/go-substrate-rpc-client/v4/registry/parser"
	"github.com/centrifuge/go-substrate-rpc-client/v4/rpc/chain/generic"
	"github.com/centrifuge/go-substrate-rpc-client/v4/rpc/state"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

//nolint:lll
//go:generate mockery --name ExtrinsicRetriever --structname ExtrinsicRetrieverMock --filename extrinsic_retriever_mock.go --inpackage

// ExtrinsicRetriever is the interface used for retrieving and decoding extrinsic information.
//
// This interface is generic over types A, S, P, please check generic.GenericExtrinsicSignature for more
// information about these generic types.
type ExtrinsicRetriever[A, S, P any] interface {
	GetExtrinsics(blockHash types.Hash) ([]*parser.Extrinsic[A, S, P], error)
}

// extrinsicRetriever implements the ExtrinsicRetriever interface.
type extrinsicRetriever[
	A, S, P any,
	B generic.GenericSignedBlock[A, S, P],
] struct {
	extrinsicParser parser.ExtrinsicParser[A, S, P]

	genericChain generic.Chain[A, S, P, B]
	stateRPC     state.State

	registryFactory registry.Factory

	chainExecutor            exec.RetryableExecutor[B]
	extrinsicParsingExecutor exec.RetryableExecutor[[]*parser.Extrinsic[A, S, P]]

	callRegistry registry.CallRegistry
	meta         *types.Metadata
}

// NewExtrinsicRetriever creates a new ExtrinsicRetriever.
func NewExtrinsicRetriever[
	A, S, P any,
	B generic.GenericSignedBlock[A, S, P],
](
	extrinsicParser parser.ExtrinsicParser[A, S, P],
	genericChain generic.Chain[A, S, P, B],
	stateRPC state.State,
	registryFactory registry.Factory,
	chainExecutor exec.RetryableExecutor[B],
	extrinsicParsingExecutor exec.RetryableExecutor[[]*parser.Extrinsic[A, S, P]],
) (ExtrinsicRetriever[A, S, P], error) {
	retriever := &extrinsicRetriever[A, S, P, B]{
		extrinsicParser:          extrinsicParser,
		genericChain:             genericChain,
		stateRPC:                 stateRPC,
		registryFactory:          registryFactory,
		chainExecutor:            chainExecutor,
		extrinsicParsingExecutor: extrinsicParsingExecutor,
	}

	if err := retriever.updateInternalState(nil); err != nil {
		return nil, ErrInternalStateUpdate.Wrap(err)
	}

	return retriever, nil
}

// DefaultExtrinsicRetriever is the ExtrinsicRetriever interface with default for the generic types:
//
// Address - types.MultiAddress
// Signature - types.MultiSignature
// PaymentFields - generic.DefaultPaymentFields
type DefaultExtrinsicRetriever = ExtrinsicRetriever[
	types.MultiAddress,
	types.MultiSignature,
	generic.DefaultPaymentFields,
]

// NewDefaultExtrinsicRetriever returns a DefaultExtrinsicRetriever with defaults for the generic types:
//
// Address - types.MultiAddress
// Signature - types.MultiSignature
// PaymentFields - generic.DefaultPaymentFields
// Block - *generic.DefaultGenericSignedBlock
//
// Note that these generic defaults also apply to the args.
func NewDefaultExtrinsicRetriever(
	extrinsicParser parser.DefaultExtrinsicParser,
	genericChain generic.DefaultChain,
	stateRPC state.State,
	registryFactory registry.Factory,
	chainExecutor exec.RetryableExecutor[*generic.DefaultGenericSignedBlock],
	extrinsicParsingExecutor exec.RetryableExecutor[[]*parser.DefaultExtrinsic],
) (DefaultExtrinsicRetriever, error) {
	return NewExtrinsicRetriever[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
		*generic.DefaultGenericSignedBlock,
	](
		extrinsicParser,
		genericChain,
		stateRPC,
		registryFactory,
		chainExecutor,
		extrinsicParsingExecutor,
	)
}

// GetExtrinsics retrieves a generic.SignedBlock and then parses the extrinsics found in it.
//
// Both the block retrieval and the extrinsic parsing are handled via the exec.RetryableExecutor
// in order to ensure retries in case of network errors or parsing errors due to an outdated call registry.
func (e *extrinsicRetriever[A, S, P, B]) GetExtrinsics(blockHash types.Hash) ([]*parser.Extrinsic[A, S, P], error) {
	block, err := e.chainExecutor.ExecWithFallback(
		func() (B, error) {
			return e.genericChain.GetBlock(blockHash)
		},
		func() error {
			return nil
		},
	)

	if err != nil {
		return nil, ErrBlockRetrieval.Wrap(err)
	}

	calls, err := e.extrinsicParsingExecutor.ExecWithFallback(
		func() ([]*parser.Extrinsic[A, S, P], error) {
			return e.extrinsicParser.ParseExtrinsics(e.callRegistry, block)
		},
		func() error {
			return e.updateInternalState(&blockHash)
		},
	)

	if err != nil {
		return nil, ErrExtrinsicParsing.Wrap(err)
	}

	return calls, nil
}

// updateInternalState will retrieve the metadata at the provided blockHash, if provided,
// create a call registry based on this metadata and store both.
func (e *extrinsicRetriever[A, S, P, B]) updateInternalState(blockHash *types.Hash) error {
	var (
		meta *types.Metadata
		err  error
	)

	if blockHash == nil {
		meta, err = e.stateRPC.GetMetadataLatest()
	} else {
		meta, err = e.stateRPC.GetMetadata(*blockHash)
	}

	if err != nil {
		return ErrMetadataRetrieval.Wrap(err)
	}

	callRegistry, err := e.registryFactory.CreateCallRegistry(meta)

	if err != nil {
		return ErrCallRegistryCreation.Wrap(err)
	}

	e.meta = meta
	e.callRegistry = callRegistry

	return nil
}
