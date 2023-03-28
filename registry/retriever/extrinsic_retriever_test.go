package retriever

import (
	"errors"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/registry"
	"github.com/centrifuge/go-substrate-rpc-client/v4/registry/exec"
	"github.com/centrifuge/go-substrate-rpc-client/v4/registry/parser"
	chainMocks "github.com/centrifuge/go-substrate-rpc-client/v4/rpc/chain/mocks"
	stateMocks "github.com/centrifuge/go-substrate-rpc-client/v4/rpc/state/mocks"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExtrinsicRetriever_New(t *testing.T) {
	extrinsicParserMock := parser.NewExtrinsicParserMock(t)
	chainRPCMock := chainMocks.NewChain(t)
	stateRPCMock := stateMocks.NewState(t)
	registryFactoryMock := registry.NewFactoryMock(t)
	chainExecMock := exec.NewRetryableExecutorMock[*types.SignedBlock](t)
	parsingExecMock := exec.NewRetryableExecutorMock[[]*parser.Extrinsic](t)

	latestMeta := &types.Metadata{}

	stateRPCMock.On("GetMetadataLatest").
		Return(latestMeta, nil).
		Once()

	callRegistry := registry.CallRegistry(map[types.CallIndex]*registry.Type{})

	registryFactoryMock.On("CreateCallRegistry", latestMeta).
		Return(callRegistry, nil).
		Once()

	res, err := NewExtrinsicRetriever(
		extrinsicParserMock,
		chainRPCMock,
		stateRPCMock,
		registryFactoryMock,
		chainExecMock,
		parsingExecMock,
	)
	assert.NoError(t, err)
	assert.IsType(t, &extrinsicRetriever{}, res)
}

func TestExtrinsicRetriever_New_InternalStateUpdateError(t *testing.T) {
	extrinsicParserMock := parser.NewExtrinsicParserMock(t)
	chainRPCMock := chainMocks.NewChain(t)
	stateRPCMock := stateMocks.NewState(t)
	registryFactoryMock := registry.NewFactoryMock(t)
	chainExecMock := exec.NewRetryableExecutorMock[*types.SignedBlock](t)
	parsingExecMock := exec.NewRetryableExecutorMock[[]*parser.Extrinsic](t)

	metadataRetrievalError := errors.New("error")

	stateRPCMock.On("GetMetadataLatest").
		Return(nil, metadataRetrievalError).
		Once()

	res, err := NewExtrinsicRetriever(
		extrinsicParserMock,
		chainRPCMock,
		stateRPCMock,
		registryFactoryMock,
		chainExecMock,
		parsingExecMock,
	)
	assert.ErrorIs(t, err, metadataRetrievalError)
	assert.Nil(t, res)

	latestMeta := &types.Metadata{}

	stateRPCMock.On("GetMetadataLatest").
		Return(latestMeta, nil).
		Once()

	registryFactoryError := errors.New("error")

	registryFactoryMock.On("CreateCallRegistry", latestMeta).
		Return(nil, registryFactoryError).
		Once()

	res, err = NewExtrinsicRetriever(
		extrinsicParserMock,
		chainRPCMock,
		stateRPCMock,
		registryFactoryMock,
		chainExecMock,
		parsingExecMock,
	)
	assert.ErrorIs(t, err, registryFactoryError)
	assert.Nil(t, res)
}

func TestExtrinsicRetriever_NewDefault(t *testing.T) {
	chainRPCMock := chainMocks.NewChain(t)
	stateRPCMock := stateMocks.NewState(t)

	latestMeta := &types.Metadata{}

	stateRPCMock.On("GetMetadataLatest").
		Return(latestMeta, nil).
		Once()

	res, err := NewDefaultExtrinsicRetriever(chainRPCMock, stateRPCMock)
	assert.NoError(t, err)
	assert.IsType(t, &extrinsicRetriever{}, res)

	retriever := res.(*extrinsicRetriever)
	assert.IsType(t, parser.NewExtrinsicParser(), retriever.extrinsicParser)
	assert.IsType(t, registry.NewFactory(), retriever.registryFactory)
	assert.IsType(t, exec.NewRetryableExecutor[*types.SignedBlock](), retriever.chainExecutor)
	assert.IsType(t, exec.NewRetryableExecutor[[]*parser.Extrinsic](), retriever.extrinsicParsingExecutor)
	assert.Equal(t, latestMeta, retriever.meta)
	assert.NotNil(t, retriever.callRegistry)
}

func TestExtrinsicRetriever_NewDefault_MetadataRetrievalError(t *testing.T) {
	chainRPCMock := chainMocks.NewChain(t)
	stateRPCMock := stateMocks.NewState(t)

	metadataRetrievalError := errors.New("error")

	stateRPCMock.On("GetMetadataLatest").
		Return(nil, metadataRetrievalError).
		Once()

	res, err := NewDefaultExtrinsicRetriever(chainRPCMock, stateRPCMock)
	assert.ErrorIs(t, err, metadataRetrievalError)
	assert.Nil(t, res)
}

func TestExtrinsicRetriever_GetExtrinsics(t *testing.T) {
	extrinsicParserMock := parser.NewExtrinsicParserMock(t)
	chainRPCMock := chainMocks.NewChain(t)
	stateRPCMock := stateMocks.NewState(t)
	registryFactoryMock := registry.NewFactoryMock(t)
	chainExecMock := exec.NewRetryableExecutorMock[*types.SignedBlock](t)
	parsingExecMock := exec.NewRetryableExecutorMock[[]*parser.Extrinsic](t)

	extrinsicRetriever := &extrinsicRetriever{
		extrinsicParser:          extrinsicParserMock,
		chainRPC:                 chainRPCMock,
		stateRPC:                 stateRPCMock,
		registryFactory:          registryFactoryMock,
		chainExecutor:            chainExecMock,
		extrinsicParsingExecutor: parsingExecMock,
	}

	testMeta := &types.Metadata{}

	extrinsicRetriever.meta = testMeta

	callRegistry := registry.CallRegistry(map[types.CallIndex]*registry.Type{})

	extrinsicRetriever.callRegistry = callRegistry

	blockHash := types.NewHash([]byte{0, 1, 2, 3})

	signedBlock := &types.SignedBlock{}

	chainRPCMock.On("GetBlock", blockHash).
		Return(signedBlock, nil).
		Once()

	chainExecMock.On("ExecWithFallback", mock.Anything, mock.Anything).
		Run(
			func(args mock.Arguments) {
				execFn, ok := args.Get(0).(func() (*types.SignedBlock, error))
				assert.True(t, ok)

				execFnRes, err := execFn()
				assert.NoError(t, err)
				assert.Equal(t, signedBlock, execFnRes)
			},
		).Return(signedBlock, nil)

	parsedExtrinsics := []*parser.Extrinsic{}

	extrinsicParserMock.On("ParseExtrinsics", callRegistry, signedBlock).
		Return(parsedExtrinsics, nil).
		Once()

	parsingExecMock.On("ExecWithFallback", mock.Anything, mock.Anything).
		Run(
			func(args mock.Arguments) {
				execFn, ok := args.Get(0).(func() ([]*parser.Extrinsic, error))
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
	extrinsicParserMock := parser.NewExtrinsicParserMock(t)
	chainRPCMock := chainMocks.NewChain(t)
	stateRPCMock := stateMocks.NewState(t)
	registryFactoryMock := registry.NewFactoryMock(t)
	chainExecMock := exec.NewRetryableExecutorMock[*types.SignedBlock](t)
	parsingExecMock := exec.NewRetryableExecutorMock[[]*parser.Extrinsic](t)

	extrinsicRetriever := &extrinsicRetriever{
		extrinsicParser:          extrinsicParserMock,
		chainRPC:                 chainRPCMock,
		stateRPC:                 stateRPCMock,
		registryFactory:          registryFactoryMock,
		chainExecutor:            chainExecMock,
		extrinsicParsingExecutor: parsingExecMock,
	}

	testMeta := &types.Metadata{}

	extrinsicRetriever.meta = testMeta

	callRegistry := registry.CallRegistry(map[types.CallIndex]*registry.Type{})

	extrinsicRetriever.callRegistry = callRegistry

	blockHash := types.NewHash([]byte{0, 1, 2, 3})

	blockRetrievalError := errors.New("error")

	chainRPCMock.On("GetBlock", blockHash).
		Return(nil, blockRetrievalError).
		Once()

	chainExecMock.On("ExecWithFallback", mock.Anything, mock.Anything).
		Run(
			func(args mock.Arguments) {
				execFn, ok := args.Get(0).(func() (*types.SignedBlock, error))
				assert.True(t, ok)

				execFnRes, err := execFn()
				assert.ErrorIs(t, err, blockRetrievalError)
				assert.Nil(t, execFnRes)

				fallbackFn, ok := args.Get(1).(func() error)
				assert.True(t, ok)

				assert.NoError(t, fallbackFn())
			},
		).Return(&types.SignedBlock{}, blockRetrievalError)

	res, err := extrinsicRetriever.GetExtrinsics(blockHash)
	assert.ErrorIs(t, err, blockRetrievalError)
	assert.Nil(t, res)
}

func TestExtrinsicRetriever_GetExtrinsics_ExtrinsicParsingError(t *testing.T) {
	extrinsicParserMock := parser.NewExtrinsicParserMock(t)
	chainRPCMock := chainMocks.NewChain(t)
	stateRPCMock := stateMocks.NewState(t)
	registryFactoryMock := registry.NewFactoryMock(t)
	chainExecMock := exec.NewRetryableExecutorMock[*types.SignedBlock](t)
	parsingExecMock := exec.NewRetryableExecutorMock[[]*parser.Extrinsic](t)

	extrinsicRetriever := &extrinsicRetriever{
		extrinsicParser:          extrinsicParserMock,
		chainRPC:                 chainRPCMock,
		stateRPC:                 stateRPCMock,
		registryFactory:          registryFactoryMock,
		chainExecutor:            chainExecMock,
		extrinsicParsingExecutor: parsingExecMock,
	}

	testMeta := &types.Metadata{}

	extrinsicRetriever.meta = testMeta

	callRegistry := registry.CallRegistry(map[types.CallIndex]*registry.Type{})

	extrinsicRetriever.callRegistry = callRegistry

	blockHash := types.NewHash([]byte{0, 1, 2, 3})

	signedBlock := &types.SignedBlock{}

	chainRPCMock.On("GetBlock", blockHash).
		Return(signedBlock, nil).
		Once()

	chainExecMock.On("ExecWithFallback", mock.Anything, mock.Anything).
		Run(
			func(args mock.Arguments) {
				execFn, ok := args.Get(0).(func() (*types.SignedBlock, error))
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
				execFn, ok := args.Get(0).(func() ([]*parser.Extrinsic, error))
				assert.True(t, ok)

				execFnRes, err := execFn()
				assert.ErrorIs(t, err, extrinsicParsingError)
				assert.Nil(t, execFnRes)

				fallbackFn, ok := args.Get(1).(func() error)
				assert.True(t, ok)

				err = fallbackFn()
				assert.NoError(t, err)
			},
		).Return([]*parser.Extrinsic{}, extrinsicParsingError)

	res, err := extrinsicRetriever.GetExtrinsics(blockHash)
	assert.ErrorIs(t, err, extrinsicParsingError)
	assert.Nil(t, res)
}

func TestExtrinsicRetriever_updateInternalState(t *testing.T) {
	extrinsicParserMock := parser.NewExtrinsicParserMock(t)
	chainRPCMock := chainMocks.NewChain(t)
	stateRPCMock := stateMocks.NewState(t)
	registryFactoryMock := registry.NewFactoryMock(t)
	chainExecMock := exec.NewRetryableExecutorMock[*types.SignedBlock](t)
	parsingExecMock := exec.NewRetryableExecutorMock[[]*parser.Extrinsic](t)

	extrinsicRetriever := &extrinsicRetriever{
		extrinsicParser:          extrinsicParserMock,
		chainRPC:                 chainRPCMock,
		stateRPC:                 stateRPCMock,
		registryFactory:          registryFactoryMock,
		chainExecutor:            chainExecMock,
		extrinsicParsingExecutor: parsingExecMock,
	}

	testMeta := &types.Metadata{}

	callRegistry := registry.CallRegistry(map[types.CallIndex]*registry.Type{})

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
	extrinsicParserMock := parser.NewExtrinsicParserMock(t)
	chainRPCMock := chainMocks.NewChain(t)
	stateRPCMock := stateMocks.NewState(t)
	registryFactoryMock := registry.NewFactoryMock(t)
	chainExecMock := exec.NewRetryableExecutorMock[*types.SignedBlock](t)
	parsingExecMock := exec.NewRetryableExecutorMock[[]*parser.Extrinsic](t)

	extrinsicRetriever := &extrinsicRetriever{
		extrinsicParser:          extrinsicParserMock,
		chainRPC:                 chainRPCMock,
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
	assert.ErrorIs(t, err, metadataRetrievalError)

	stateRPCMock.On("GetMetadataLatest").
		Return(nil, metadataRetrievalError).
		Once()

	err = extrinsicRetriever.updateInternalState(nil)
	assert.ErrorIs(t, err, metadataRetrievalError)
}

func TestExtrinsicRetriever_updateInternalState_RegistryFactoryError(t *testing.T) {
	extrinsicParserMock := parser.NewExtrinsicParserMock(t)
	chainRPCMock := chainMocks.NewChain(t)
	stateRPCMock := stateMocks.NewState(t)
	registryFactoryMock := registry.NewFactoryMock(t)
	chainExecMock := exec.NewRetryableExecutorMock[*types.SignedBlock](t)
	parsingExecMock := exec.NewRetryableExecutorMock[[]*parser.Extrinsic](t)

	extrinsicRetriever := &extrinsicRetriever{
		extrinsicParser:          extrinsicParserMock,
		chainRPC:                 chainRPCMock,
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
	assert.ErrorIs(t, err, registryFactoryError)

	latestMeta := &types.Metadata{}

	stateRPCMock.On("GetMetadataLatest").
		Return(latestMeta, nil).
		Once()

	registryFactoryMock.On("CreateCallRegistry", latestMeta).
		Return(nil, registryFactoryError).
		Once()

	err = extrinsicRetriever.updateInternalState(nil)
	assert.ErrorIs(t, err, registryFactoryError)
}
