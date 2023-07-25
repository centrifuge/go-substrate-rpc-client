package retriever

import (
	"errors"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/registry"
	"github.com/centrifuge/go-substrate-rpc-client/v4/registry/exec"
	"github.com/centrifuge/go-substrate-rpc-client/v4/registry/parser"
	"github.com/centrifuge/go-substrate-rpc-client/v4/rpc/chain/generic"
	stateMocks "github.com/centrifuge/go-substrate-rpc-client/v4/rpc/state/mocks"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExtrinsicRetriever_New(t *testing.T) {
	extrinsicParserMock := parser.NewExtrinsicParserMock[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	](t)

	genericChainMock := generic.NewChainMock[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
		*generic.SignedBlock[
			types.MultiAddress,
			types.MultiSignature,
			generic.DefaultPaymentFields,
		],
	](t)

	stateRPCMock := stateMocks.NewState(t)
	registryFactoryMock := registry.NewFactoryMock(t)
	chainExecMock := exec.NewRetryableExecutorMock[*generic.SignedBlock[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	]](t)
	parsingExecMock := exec.NewRetryableExecutorMock[[]*parser.Extrinsic[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	]](t)

	latestMeta := &types.Metadata{}

	stateRPCMock.On("GetMetadataLatest").
		Return(latestMeta, nil).
		Once()

	callRegistry := registry.CallRegistry(map[types.CallIndex]*registry.TypeDecoder{})

	registryFactoryMock.On("CreateCallRegistry", latestMeta).
		Return(callRegistry, nil).
		Once()

	res, err := NewExtrinsicRetriever[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
		*generic.SignedBlock[
			types.MultiAddress,
			types.MultiSignature,
			generic.DefaultPaymentFields,
		],
	](
		extrinsicParserMock,
		genericChainMock,
		stateRPCMock,
		registryFactoryMock,
		chainExecMock,
		parsingExecMock,
	)
	assert.NoError(t, err)
	assert.IsType(t, &extrinsicRetriever[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
		*generic.SignedBlock[
			types.MultiAddress,
			types.MultiSignature,
			generic.DefaultPaymentFields,
		],
	]{}, res)
}

func TestExtrinsicRetriever_New_InternalStateUpdateError(t *testing.T) {
	extrinsicParserMock := parser.NewExtrinsicParserMock[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	](t)

	genericChainMock := generic.NewChainMock[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
		*generic.SignedBlock[
			types.MultiAddress,
			types.MultiSignature,
			generic.DefaultPaymentFields,
		],
	](t)

	stateRPCMock := stateMocks.NewState(t)
	registryFactoryMock := registry.NewFactoryMock(t)
	chainExecMock := exec.NewRetryableExecutorMock[*generic.SignedBlock[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	]](t)
	parsingExecMock := exec.NewRetryableExecutorMock[[]*parser.Extrinsic[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	]](t)
	metadataRetrievalError := errors.New("error")

	stateRPCMock.On("GetMetadataLatest").
		Return(nil, metadataRetrievalError).
		Once()

	res, err := NewExtrinsicRetriever[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
		*generic.SignedBlock[
			types.MultiAddress,
			types.MultiSignature,
			generic.DefaultPaymentFields,
		],
	](
		extrinsicParserMock,
		genericChainMock,
		stateRPCMock,
		registryFactoryMock,
		chainExecMock,
		parsingExecMock,
	)
	assert.ErrorIs(t, err, ErrInternalStateUpdate)
	assert.Nil(t, res)

	latestMeta := &types.Metadata{}

	stateRPCMock.On("GetMetadataLatest").
		Return(latestMeta, nil).
		Once()

	registryFactoryError := errors.New("error")

	registryFactoryMock.On("CreateCallRegistry", latestMeta).
		Return(nil, registryFactoryError).
		Once()

	res, err = NewExtrinsicRetriever[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
		*generic.SignedBlock[
			types.MultiAddress,
			types.MultiSignature,
			generic.DefaultPaymentFields,
		],
	](
		extrinsicParserMock,
		genericChainMock,
		stateRPCMock,
		registryFactoryMock,
		chainExecMock,
		parsingExecMock,
	)
	assert.ErrorIs(t, err, ErrInternalStateUpdate)
	assert.Nil(t, res)
}

func TestExtrinsicRetriever_NewDefault(t *testing.T) {
	genericChainMock := generic.NewChainMock[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
		*generic.DefaultGenericSignedBlock,
	](t)

	stateRPCMock := stateMocks.NewState(t)

	latestMeta := &types.Metadata{}

	stateRPCMock.On("GetMetadataLatest").
		Return(latestMeta, nil).
		Once()

	extrinsicParser := parser.NewDefaultExtrinsicParser()

	registryFactoryMock := registry.NewFactoryMock(t)

	callRegistry := registry.CallRegistry(map[types.CallIndex]*registry.TypeDecoder{})

	registryFactoryMock.On("CreateCallRegistry", latestMeta).
		Return(callRegistry, nil).
		Once()

	chainExecMock := exec.NewRetryableExecutorMock[*generic.DefaultGenericSignedBlock](t)
	parsingExecMock := exec.NewRetryableExecutorMock[[]*parser.DefaultExtrinsic](t)

	res, err := NewDefaultExtrinsicRetriever(
		extrinsicParser,
		genericChainMock,
		stateRPCMock,
		registryFactoryMock,
		chainExecMock,
		parsingExecMock,
	)

	assert.NoError(t, err)
	assert.IsType(t, &extrinsicRetriever[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
		*generic.DefaultGenericSignedBlock,
	]{}, res)

	retriever := res.(*extrinsicRetriever[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
		*generic.DefaultGenericSignedBlock,
	])
	assert.Equal(t, latestMeta, retriever.meta)
	assert.Equal(t, callRegistry, retriever.callRegistry)
}

func TestExtrinsicRetriever_GetExtrinsics(t *testing.T) {
	extrinsicParserMock := parser.NewExtrinsicParserMock[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	](t)

	genericChainMock := generic.NewChainMock[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
		*generic.SignedBlock[
			types.MultiAddress,
			types.MultiSignature,
			generic.DefaultPaymentFields,
		],
	](t)

	stateRPCMock := stateMocks.NewState(t)
	registryFactoryMock := registry.NewFactoryMock(t)
	chainExecMock := exec.NewRetryableExecutorMock[*generic.SignedBlock[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	]](t)
	parsingExecMock := exec.NewRetryableExecutorMock[[]*parser.Extrinsic[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	]](t)

	extrinsicRetriever := &extrinsicRetriever[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
		*generic.SignedBlock[
			types.MultiAddress,
			types.MultiSignature,
			generic.DefaultPaymentFields,
		],
	]{
		extrinsicParser:          extrinsicParserMock,
		genericChain:             genericChainMock,
		stateRPC:                 stateRPCMock,
		registryFactory:          registryFactoryMock,
		chainExecutor:            chainExecMock,
		extrinsicParsingExecutor: parsingExecMock,
	}

	testMeta := &types.Metadata{}

	extrinsicRetriever.meta = testMeta

	callRegistry := registry.CallRegistry(map[types.CallIndex]*registry.TypeDecoder{})

	extrinsicRetriever.callRegistry = callRegistry

	blockHash := types.NewHash([]byte{0, 1, 2, 3})

	signedBlock := &generic.SignedBlock[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	]{}

	genericChainMock.On("GetBlock", blockHash).
		Return(signedBlock, nil).
		Once()

	chainExecMock.On("ExecWithFallback", mock.Anything, mock.Anything).
		Run(
			func(args mock.Arguments) {
				execFn, ok := args.Get(0).(func() (*generic.SignedBlock[
					types.MultiAddress,
					types.MultiSignature,
					generic.DefaultPaymentFields,
				], error))
				assert.True(t, ok)

				execFnRes, err := execFn()
				assert.NoError(t, err)
				assert.Equal(t, signedBlock, execFnRes)
			},
		).Return(signedBlock, nil)

	var parsedExtrinsics []*parser.Extrinsic[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	]

	extrinsicParserMock.On("ParseExtrinsics", callRegistry, signedBlock).
		Return(parsedExtrinsics, nil).
		Once()

	parsingExecMock.On("ExecWithFallback", mock.Anything, mock.Anything).
		Run(
			func(args mock.Arguments) {
				execFn, ok := args.Get(0).(func() ([]*parser.Extrinsic[
					types.MultiAddress,
					types.MultiSignature,
					generic.DefaultPaymentFields,
				], error))
				assert.True(t, ok)

				execFnRes, err := execFn()
				assert.NoError(t, err)
				assert.Equal(t, parsedExtrinsics, execFnRes)
			},
		).Return(parsedExtrinsics, nil)

	res, err := extrinsicRetriever.GetExtrinsics(blockHash)
	assert.NoError(t, err)
	assert.Equal(t, parsedExtrinsics, res)
}

func TestExtrinsicRetriever_GetExtrinsics_BlockRetrievalError(t *testing.T) {
	extrinsicParserMock := parser.NewExtrinsicParserMock[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	](t)

	genericChainMock := generic.NewChainMock[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
		*generic.SignedBlock[
			types.MultiAddress,
			types.MultiSignature,
			generic.DefaultPaymentFields,
		],
	](t)

	stateRPCMock := stateMocks.NewState(t)
	registryFactoryMock := registry.NewFactoryMock(t)
	chainExecMock := exec.NewRetryableExecutorMock[*generic.SignedBlock[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	]](t)
	parsingExecMock := exec.NewRetryableExecutorMock[[]*parser.Extrinsic[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	]](t)

	extrinsicRetriever := &extrinsicRetriever[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
		*generic.SignedBlock[
			types.MultiAddress,
			types.MultiSignature,
			generic.DefaultPaymentFields,
		],
	]{
		extrinsicParser:          extrinsicParserMock,
		genericChain:             genericChainMock,
		stateRPC:                 stateRPCMock,
		registryFactory:          registryFactoryMock,
		chainExecutor:            chainExecMock,
		extrinsicParsingExecutor: parsingExecMock,
	}

	testMeta := &types.Metadata{}

	extrinsicRetriever.meta = testMeta

	callRegistry := registry.CallRegistry(map[types.CallIndex]*registry.TypeDecoder{})

	extrinsicRetriever.callRegistry = callRegistry

	blockHash := types.NewHash([]byte{0, 1, 2, 3})

	blockRetrievalError := errors.New("error")

	signedBlock := &generic.SignedBlock[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	]{}

	genericChainMock.On("GetBlock", blockHash).
		Return(signedBlock, blockRetrievalError).
		Once()

	chainExecMock.On("ExecWithFallback", mock.Anything, mock.Anything).
		Run(
			func(args mock.Arguments) {
				execFn, ok := args.Get(0).(func() (*generic.SignedBlock[
					types.MultiAddress,
					types.MultiSignature,
					generic.DefaultPaymentFields,
				], error))
				assert.True(t, ok)

				execFnRes, err := execFn()
				assert.ErrorIs(t, err, blockRetrievalError)
				assert.Equal(t, signedBlock, execFnRes)

				fallbackFn, ok := args.Get(1).(func() error)
				assert.True(t, ok)

				assert.NoError(t, fallbackFn())
			},
		).Return(signedBlock, blockRetrievalError)

	res, err := extrinsicRetriever.GetExtrinsics(blockHash)
	assert.ErrorIs(t, err, ErrBlockRetrieval)
	assert.Nil(t, res)
}

func TestExtrinsicRetriever_GetExtrinsics_ExtrinsicParsingError(t *testing.T) {
	extrinsicParserMock := parser.NewExtrinsicParserMock[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	](t)

	genericChainMock := generic.NewChainMock[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
		*generic.SignedBlock[
			types.MultiAddress,
			types.MultiSignature,
			generic.DefaultPaymentFields,
		],
	](t)

	stateRPCMock := stateMocks.NewState(t)
	registryFactoryMock := registry.NewFactoryMock(t)
	chainExecMock := exec.NewRetryableExecutorMock[*generic.SignedBlock[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	]](t)
	parsingExecMock := exec.NewRetryableExecutorMock[[]*parser.Extrinsic[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	]](t)

	extrinsicRetriever := &extrinsicRetriever[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
		*generic.SignedBlock[
			types.MultiAddress,
			types.MultiSignature,
			generic.DefaultPaymentFields,
		],
	]{
		extrinsicParser:          extrinsicParserMock,
		genericChain:             genericChainMock,
		stateRPC:                 stateRPCMock,
		registryFactory:          registryFactoryMock,
		chainExecutor:            chainExecMock,
		extrinsicParsingExecutor: parsingExecMock,
	}

	testMeta := &types.Metadata{}

	extrinsicRetriever.meta = testMeta

	callRegistry := registry.CallRegistry(map[types.CallIndex]*registry.TypeDecoder{})

	extrinsicRetriever.callRegistry = callRegistry

	blockHash := types.NewHash([]byte{0, 1, 2, 3})

	signedBlock := &generic.SignedBlock[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	]{}

	genericChainMock.On("GetBlock", blockHash).
		Return(signedBlock, nil).
		Once()

	chainExecMock.On("ExecWithFallback", mock.Anything, mock.Anything).
		Run(
			func(args mock.Arguments) {
				execFn, ok := args.Get(0).(func() (*generic.SignedBlock[
					types.MultiAddress,
					types.MultiSignature,
					generic.DefaultPaymentFields,
				], error))
				assert.True(t, ok)

				execFnRes, err := execFn()
				assert.NoError(t, err)
				assert.Equal(t, signedBlock, execFnRes)
			},
		).Return(signedBlock, nil)

	extrinsicParsingError := errors.New("error")

	extrinsicParserMock.On("ParseExtrinsics", callRegistry, signedBlock).
		Return(nil, extrinsicParsingError).
		Once()

	stateRPCMock.On("GetMetadata", blockHash).
		Return(testMeta, nil).
		Once()

	registryFactoryMock.On("CreateCallRegistry", testMeta).
		Return(callRegistry, nil).
		Once()

	parsingExecMock.On("ExecWithFallback", mock.Anything, mock.Anything).
		Run(
			func(args mock.Arguments) {
				execFn, ok := args.Get(0).(func() ([]*parser.Extrinsic[
					types.MultiAddress,
					types.MultiSignature,
					generic.DefaultPaymentFields,
				], error))
				assert.True(t, ok)

				execFnRes, err := execFn()
				assert.ErrorIs(t, err, extrinsicParsingError)
				assert.Nil(t, execFnRes)

				fallbackFn, ok := args.Get(1).(func() error)
				assert.True(t, ok)

				err = fallbackFn()
				assert.NoError(t, err)
			},
		).Return([]*parser.Extrinsic[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	]{}, extrinsicParsingError)

	res, err := extrinsicRetriever.GetExtrinsics(blockHash)
	assert.ErrorIs(t, err, ErrExtrinsicParsing)
	assert.Nil(t, res)
}

func TestExtrinsicRetriever_updateInternalState(t *testing.T) {
	extrinsicParserMock := parser.NewExtrinsicParserMock[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	](t)

	genericChainMock := generic.NewChainMock[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
		*generic.SignedBlock[
			types.MultiAddress,
			types.MultiSignature,
			generic.DefaultPaymentFields,
		],
	](t)

	stateRPCMock := stateMocks.NewState(t)
	registryFactoryMock := registry.NewFactoryMock(t)
	chainExecMock := exec.NewRetryableExecutorMock[*generic.SignedBlock[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	]](t)
	parsingExecMock := exec.NewRetryableExecutorMock[[]*parser.Extrinsic[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	]](t)

	extrinsicRetriever := &extrinsicRetriever[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
		*generic.SignedBlock[
			types.MultiAddress,
			types.MultiSignature,
			generic.DefaultPaymentFields,
		],
	]{
		extrinsicParser:          extrinsicParserMock,
		genericChain:             genericChainMock,
		stateRPC:                 stateRPCMock,
		registryFactory:          registryFactoryMock,
		chainExecutor:            chainExecMock,
		extrinsicParsingExecutor: parsingExecMock,
	}

	testMeta := &types.Metadata{}

	callRegistry := registry.CallRegistry(map[types.CallIndex]*registry.TypeDecoder{})

	blockHash := types.NewHash([]byte{0, 1, 2, 3})

	stateRPCMock.On("GetMetadata", blockHash).
		Return(testMeta, nil).
		Once()

	registryFactoryMock.On("CreateCallRegistry", testMeta).
		Return(callRegistry, nil).
		Once()

	err := extrinsicRetriever.updateInternalState(&blockHash)
	assert.NoError(t, err)

	latestMeta := &types.Metadata{}

	stateRPCMock.On("GetMetadataLatest").
		Return(latestMeta, nil).
		Once()

	registryFactoryMock.On("CreateCallRegistry", latestMeta).
		Return(callRegistry, nil).
		Once()

	err = extrinsicRetriever.updateInternalState(nil)
	assert.NoError(t, err)
}

func TestExtrinsicRetriever_updateInternalState_MetadataRetrievalError(t *testing.T) {
	extrinsicParserMock := parser.NewExtrinsicParserMock[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	](t)

	genericChainMock := generic.NewChainMock[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
		*generic.SignedBlock[
			types.MultiAddress,
			types.MultiSignature,
			generic.DefaultPaymentFields,
		],
	](t)

	stateRPCMock := stateMocks.NewState(t)
	registryFactoryMock := registry.NewFactoryMock(t)
	chainExecMock := exec.NewRetryableExecutorMock[*generic.SignedBlock[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	]](t)
	parsingExecMock := exec.NewRetryableExecutorMock[[]*parser.Extrinsic[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	]](t)

	extrinsicRetriever := &extrinsicRetriever[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
		*generic.SignedBlock[
			types.MultiAddress,
			types.MultiSignature,
			generic.DefaultPaymentFields,
		],
	]{
		extrinsicParser:          extrinsicParserMock,
		genericChain:             genericChainMock,
		stateRPC:                 stateRPCMock,
		registryFactory:          registryFactoryMock,
		chainExecutor:            chainExecMock,
		extrinsicParsingExecutor: parsingExecMock,
	}

	metadataRetrievalError := errors.New("error")

	blockHash := types.NewHash([]byte{0, 1, 2, 3})

	stateRPCMock.On("GetMetadata", blockHash).
		Return(nil, metadataRetrievalError).
		Once()

	err := extrinsicRetriever.updateInternalState(&blockHash)
	assert.ErrorIs(t, err, ErrMetadataRetrieval)

	stateRPCMock.On("GetMetadataLatest").
		Return(nil, metadataRetrievalError).
		Once()

	err = extrinsicRetriever.updateInternalState(nil)
	assert.ErrorIs(t, err, ErrMetadataRetrieval)
}

func TestExtrinsicRetriever_updateInternalState_RegistryFactoryError(t *testing.T) {
	extrinsicParserMock := parser.NewExtrinsicParserMock[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	](t)

	genericChainMock := generic.NewChainMock[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
		*generic.SignedBlock[
			types.MultiAddress,
			types.MultiSignature,
			generic.DefaultPaymentFields,
		],
	](t)

	stateRPCMock := stateMocks.NewState(t)
	registryFactoryMock := registry.NewFactoryMock(t)
	chainExecMock := exec.NewRetryableExecutorMock[*generic.SignedBlock[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	]](t)
	parsingExecMock := exec.NewRetryableExecutorMock[[]*parser.Extrinsic[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
	]](t)

	extrinsicRetriever := &extrinsicRetriever[
		types.MultiAddress,
		types.MultiSignature,
		generic.DefaultPaymentFields,
		*generic.SignedBlock[
			types.MultiAddress,
			types.MultiSignature,
			generic.DefaultPaymentFields,
		],
	]{
		extrinsicParser:          extrinsicParserMock,
		genericChain:             genericChainMock,
		stateRPC:                 stateRPCMock,
		registryFactory:          registryFactoryMock,
		chainExecutor:            chainExecMock,
		extrinsicParsingExecutor: parsingExecMock,
	}

	testMeta := &types.Metadata{}

	registryFactoryError := errors.New("error")

	blockHash := types.NewHash([]byte{0, 1, 2, 3})

	stateRPCMock.On("GetMetadata", blockHash).
		Return(testMeta, nil).
		Once()

	registryFactoryMock.On("CreateCallRegistry", testMeta).
		Return(nil, registryFactoryError).
		Once()

	err := extrinsicRetriever.updateInternalState(&blockHash)
	assert.ErrorIs(t, err, ErrCallRegistryCreation)

	latestMeta := &types.Metadata{}

	stateRPCMock.On("GetMetadataLatest").
		Return(latestMeta, nil).
		Once()

	registryFactoryMock.On("CreateCallRegistry", latestMeta).
		Return(nil, registryFactoryError).
		Once()

	err = extrinsicRetriever.updateInternalState(nil)
	assert.ErrorIs(t, err, ErrCallRegistryCreation)
}
