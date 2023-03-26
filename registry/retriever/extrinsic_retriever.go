package retriever

import (
	"fmt"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/registry"
	"github.com/centrifuge/go-substrate-rpc-client/v4/registry/exec"
	"github.com/centrifuge/go-substrate-rpc-client/v4/registry/parser"
	"github.com/centrifuge/go-substrate-rpc-client/v4/rpc/chain"
	"github.com/centrifuge/go-substrate-rpc-client/v4/rpc/state"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

type ExtrinsicRetriever interface {
	GetExtrinsics(blockHash types.Hash) ([]*parser.Extrinsic, error)
}

type extrinsicRetriever struct {
	extrinsicParser parser.ExtrinsicParser

	chainRPC chain.Chain
	stateRPC state.State

	registryFactory registry.Factory

	chainExecutor            exec.RetryableExecutor[*types.SignedBlock]
	extrinsicParsingExecutor exec.RetryableExecutor[[]*parser.Extrinsic]

	callRegistry registry.CallRegistry
	meta         *types.Metadata
}

func NewExtrinsicRetriever(
	callParser parser.ExtrinsicParser,
	chainRPC chain.Chain,
	stateRPC state.State,
	registryFactory registry.Factory,
	chainExecutor exec.RetryableExecutor[*types.SignedBlock],
	callParsingExecutor exec.RetryableExecutor[[]*parser.Extrinsic],
) (ExtrinsicRetriever, error) {
	retriever := &extrinsicRetriever{
		extrinsicParser:          callParser,
		chainRPC:                 chainRPC,
		stateRPC:                 stateRPC,
		registryFactory:          registryFactory,
		chainExecutor:            chainExecutor,
		extrinsicParsingExecutor: callParsingExecutor,
	}

	if err := retriever.updateInternalState(nil); err != nil {
		return nil, err
	}

	return retriever, nil
}

func NewDefaultExtrinsicRetriever(
	chainRPC chain.Chain,
	stateRPC state.State,
) (ExtrinsicRetriever, error) {
	callParser := parser.NewExtrinsicParser()
	registryFactory := registry.NewFactory()

	chainExecutor := exec.NewRetryableExecutor[*types.SignedBlock](exec.WithRetryTimeout(1 * time.Second))
	callParsingExecutor := exec.NewRetryableExecutor[[]*parser.Extrinsic](exec.WithMaxRetryCount(1))

	return NewExtrinsicRetriever(callParser, chainRPC, stateRPC, registryFactory, chainExecutor, callParsingExecutor)
}

func (e *extrinsicRetriever) GetExtrinsics(blockHash types.Hash) ([]*parser.Extrinsic, error) {
	block, err := e.chainExecutor.ExecWithFallback(
		func() (*types.SignedBlock, error) {
			return e.chainRPC.GetBlock(blockHash)
		},
		func() error {
			return nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("couldn't retrieve block: %w", err)
	}

	calls, err := e.extrinsicParsingExecutor.ExecWithFallback(
		func() ([]*parser.Extrinsic, error) {
			return e.extrinsicParser.ParseExtrinsics(e.callRegistry, block)
		},
		func() error {
			return e.updateInternalState(&blockHash)
		},
	)

	if err != nil {
		return nil, fmt.Errorf("couldn't parse calls: %w", err)
	}

	return calls, nil
}

func (e *extrinsicRetriever) updateInternalState(blockHash *types.Hash) error {
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
		return fmt.Errorf("couldn't retrieve metadata: %w", err)
	}

	callRegistry, err := e.registryFactory.CreateCallRegistry(meta)

	if err != nil {
		return fmt.Errorf("couldn't create call registry: %w", err)
	}

	e.meta = meta
	e.callRegistry = callRegistry

	return nil
}
