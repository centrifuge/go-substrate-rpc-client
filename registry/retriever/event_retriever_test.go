package retriever

import (
	"errors"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/registry"
	"github.com/centrifuge/go-substrate-rpc-client/v4/registry/exec"
	"github.com/centrifuge/go-substrate-rpc-client/v4/registry/parser"
	"github.com/centrifuge/go-substrate-rpc-client/v4/registry/state"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEventRetriever_New(t *testing.T) {
	eventParserMock := parser.NewEventParserMock(t)
	stateProviderMock := state.NewProviderMock(t)
	registryFactoryMock := registry.NewFactoryMock(t)
	storageExecMock := exec.NewRetryableExecutorMock[*types.StorageDataRaw](t)
	parsingExecMock := exec.NewRetryableExecutorMock[[]*parser.Event](t)

	latestMeta := &types.Metadata{}

	stateProviderMock.On("GetLatestMetadata").
		Return(latestMeta, nil).
		Once()

	eventRegistry := registry.EventRegistry(map[types.EventID]*registry.Type{})

	registryFactoryMock.On("CreateEventRegistry", latestMeta).
		Return(eventRegistry, nil).
		Once()

	res, err := NewEventRetriever(
		eventParserMock,
		stateProviderMock,
		registryFactoryMock,
		storageExecMock,
		parsingExecMock,
	)
	assert.NoError(t, err)
	assert.IsType(t, &eventRetriever{}, res)
}

func TestEventRetriever_New_InternalStateUpdateError(t *testing.T) {
	eventParserMock := parser.NewEventParserMock(t)
	stateProviderMock := state.NewProviderMock(t)
	registryFactoryMock := registry.NewFactoryMock(t)
	storageExecMock := exec.NewRetryableExecutorMock[*types.StorageDataRaw](t)
	parsingExecMock := exec.NewRetryableExecutorMock[[]*parser.Event](t)

	metadataRetrievalError := errors.New("error")

	stateProviderMock.On("GetLatestMetadata").
		Return(nil, metadataRetrievalError).
		Once()

	res, err := NewEventRetriever(
		eventParserMock,
		stateProviderMock,
		registryFactoryMock,
		storageExecMock,
		parsingExecMock,
	)
	assert.ErrorIs(t, err, metadataRetrievalError)
	assert.Nil(t, res)

	latestMeta := &types.Metadata{}

	stateProviderMock.On("GetLatestMetadata").
		Return(latestMeta, nil).
		Once()

	registryFactoryError := errors.New("error")

	registryFactoryMock.On("CreateEventRegistry", latestMeta).
		Return(nil, registryFactoryError).
		Once()

	res, err = NewEventRetriever(
		eventParserMock,
		stateProviderMock,
		registryFactoryMock,
		storageExecMock,
		parsingExecMock,
	)
	assert.ErrorIs(t, err, registryFactoryError)
	assert.Nil(t, res)
}

func TestEventRetriever_NewDefault(t *testing.T) {
	stateProviderMock := state.NewProviderMock(t)

	latestMeta := &types.Metadata{}

	stateProviderMock.On("GetLatestMetadata").
		Return(latestMeta, nil).
		Once()

	res, err := NewDefaultEventRetriever(stateProviderMock)
	assert.NoError(t, err)
	assert.IsType(t, &eventRetriever{}, res)

	retriever := res.(*eventRetriever)
	assert.IsType(t, parser.NewEventParser(), retriever.eventParser)
	assert.IsType(t, registry.NewFactory(), retriever.registryFactory)
	assert.IsType(t, exec.NewRetryableExecutor[*types.StorageDataRaw](), retriever.eventStorageExecutor)
	assert.IsType(t, exec.NewRetryableExecutor[[]*parser.Event](), retriever.eventParsingExecutor)
	assert.Equal(t, latestMeta, retriever.meta)
	assert.NotNil(t, retriever.eventRegistry)
}

func TestEventRetriever_NewDefault_MetadataRetrievalError(t *testing.T) {
	stateProviderMock := state.NewProviderMock(t)

	metadataRetrievalError := errors.New("error")

	stateProviderMock.On("GetLatestMetadata").
		Return(nil, metadataRetrievalError).
		Once()

	res, err := NewDefaultEventRetriever(stateProviderMock)
	assert.ErrorIs(t, err, metadataRetrievalError)
	assert.Nil(t, res)
}

func TestEventRetriever_GetEvents(t *testing.T) {
	eventParserMock := parser.NewEventParserMock(t)
	stateProviderMock := state.NewProviderMock(t)
	registryFactoryMock := registry.NewFactoryMock(t)
	storageExecMock := exec.NewRetryableExecutorMock[*types.StorageDataRaw](t)
	parsingExecMock := exec.NewRetryableExecutorMock[[]*parser.Event](t)

	eventRetriever := &eventRetriever{
		eventParser:          eventParserMock,
		stateProvider:        stateProviderMock,
		registryFactory:      registryFactoryMock,
		eventStorageExecutor: storageExecMock,
		eventParsingExecutor: parsingExecMock,
	}

	testMeta := &types.Metadata{}

	eventRetriever.meta = testMeta

	eventRegistry := registry.EventRegistry(map[types.EventID]*registry.Type{})

	eventRetriever.eventRegistry = eventRegistry

	blockHash := types.NewHash([]byte{0, 1, 2, 3})

	storageEvents := &types.StorageDataRaw{}

	stateProviderMock.On("GetStorageEvents", testMeta, blockHash).
		Return(storageEvents, nil).
		Once()

	storageExecMock.On("ExecWithFallback", mock.Anything, mock.Anything).
		Run(
			func(args mock.Arguments) {
				execFn, ok := args.Get(0).(func() (*types.StorageDataRaw, error))
				assert.True(t, ok)

				execFnRes, err := execFn()
				assert.NoError(t, err)
				assert.Equal(t, storageEvents, execFnRes)
			},
		).Return(storageEvents, nil)

	parsedEvents := []*parser.Event{}

	eventParserMock.On("ParseEvents", eventRegistry, storageEvents).
		Return(parsedEvents, nil).
		Once()

	parsingExecMock.On("ExecWithFallback", mock.Anything, mock.Anything).
		Run(
			func(args mock.Arguments) {
				execFn, ok := args.Get(0).(func() ([]*parser.Event, error))
				assert.True(t, ok)

				execFnRes, err := execFn()
				assert.NoError(t, err)
				assert.Equal(t, parsedEvents, execFnRes)
			},
		).Return(parsedEvents, nil)

	res, err := eventRetriever.GetEvents(blockHash)
	assert.NoError(t, err)
	assert.Equal(t, parsedEvents, res)
}

func TestEventRetriever_GetEvents_StorageRetrievalError(t *testing.T) {
	eventParserMock := parser.NewEventParserMock(t)
	stateProviderMock := state.NewProviderMock(t)
	registryFactoryMock := registry.NewFactoryMock(t)
	storageExecMock := exec.NewRetryableExecutorMock[*types.StorageDataRaw](t)
	parsingExecMock := exec.NewRetryableExecutorMock[[]*parser.Event](t)

	eventRetriever := &eventRetriever{
		eventParser:          eventParserMock,
		stateProvider:        stateProviderMock,
		registryFactory:      registryFactoryMock,
		eventStorageExecutor: storageExecMock,
		eventParsingExecutor: parsingExecMock,
	}

	testMeta := &types.Metadata{}

	eventRetriever.meta = testMeta

	eventRegistry := registry.EventRegistry(map[types.EventID]*registry.Type{})

	eventRetriever.eventRegistry = eventRegistry

	blockHash := types.NewHash([]byte{0, 1, 2, 3})

	storageRetrievalError := errors.New("error")

	stateProviderMock.On("GetStorageEvents", testMeta, blockHash).
		Return(nil, storageRetrievalError).
		Once()

	stateProviderMock.On("GetMetadata", blockHash).
		Return(testMeta, nil).
		Once()

	registryFactoryMock.On("CreateEventRegistry", testMeta).
		Return(eventRegistry, nil).
		Once()

	storageExecMock.On("ExecWithFallback", mock.Anything, mock.Anything).
		Run(
			func(args mock.Arguments) {
				execFn, ok := args.Get(0).(func() (*types.StorageDataRaw, error))
				assert.True(t, ok)

				execFnRes, err := execFn()
				assert.ErrorIs(t, err, storageRetrievalError)
				assert.Nil(t, execFnRes)

				fallbackFn, ok := args.Get(1).(func() error)
				assert.True(t, ok)

				err = fallbackFn()
				assert.NoError(t, err)
			},
		).Return(&types.StorageDataRaw{}, storageRetrievalError)

	res, err := eventRetriever.GetEvents(blockHash)
	assert.ErrorIs(t, err, storageRetrievalError)
	assert.Nil(t, res)
}

func TestEventRetriever_GetEvents_EventParsingError(t *testing.T) {
	eventParserMock := parser.NewEventParserMock(t)
	stateProviderMock := state.NewProviderMock(t)
	registryFactoryMock := registry.NewFactoryMock(t)
	storageExecMock := exec.NewRetryableExecutorMock[*types.StorageDataRaw](t)
	parsingExecMock := exec.NewRetryableExecutorMock[[]*parser.Event](t)

	eventRetriever := &eventRetriever{
		eventParser:          eventParserMock,
		stateProvider:        stateProviderMock,
		registryFactory:      registryFactoryMock,
		eventStorageExecutor: storageExecMock,
		eventParsingExecutor: parsingExecMock,
	}

	testMeta := &types.Metadata{}

	eventRetriever.meta = testMeta

	eventRegistry := registry.EventRegistry(map[types.EventID]*registry.Type{})

	eventRetriever.eventRegistry = eventRegistry

	blockHash := types.NewHash([]byte{0, 1, 2, 3})

	storageEvents := &types.StorageDataRaw{}

	stateProviderMock.On("GetStorageEvents", testMeta, blockHash).
		Return(storageEvents, nil).
		Once()

	storageExecMock.On("ExecWithFallback", mock.Anything, mock.Anything).
		Run(
			func(args mock.Arguments) {
				execFn, ok := args.Get(0).(func() (*types.StorageDataRaw, error))
				assert.True(t, ok)

				execFnRes, err := execFn()
				assert.NoError(t, err)
				assert.Equal(t, storageEvents, execFnRes)
			},
		).Return(storageEvents, nil)

	eventParsingError := errors.New("error")

	eventParserMock.On("ParseEvents", eventRegistry, storageEvents).
		Return(nil, eventParsingError).
		Once()

	stateProviderMock.On("GetMetadata", blockHash).
		Return(testMeta, nil).
		Once()

	registryFactoryMock.On("CreateEventRegistry", testMeta).
		Return(eventRegistry, nil).
		Once()

	parsingExecMock.On("ExecWithFallback", mock.Anything, mock.Anything).
		Run(
			func(args mock.Arguments) {
				execFn, ok := args.Get(0).(func() ([]*parser.Event, error))
				assert.True(t, ok)

				execFnRes, err := execFn()
				assert.ErrorIs(t, err, eventParsingError)
				assert.Nil(t, execFnRes)

				fallbackFn, ok := args.Get(1).(func() error)
				assert.True(t, ok)

				err = fallbackFn()
				assert.NoError(t, err)
			},
		).Return([]*parser.Event{}, eventParsingError)

	res, err := eventRetriever.GetEvents(blockHash)
	assert.ErrorIs(t, err, eventParsingError)
	assert.Nil(t, res)
}

func TestEventRetriever_updateInternalState(t *testing.T) {
	eventParserMock := parser.NewEventParserMock(t)
	stateProviderMock := state.NewProviderMock(t)
	registryFactoryMock := registry.NewFactoryMock(t)
	storageExecMock := exec.NewRetryableExecutorMock[*types.StorageDataRaw](t)
	parsingExecMock := exec.NewRetryableExecutorMock[[]*parser.Event](t)

	eventRetriever := &eventRetriever{
		eventParser:          eventParserMock,
		stateProvider:        stateProviderMock,
		registryFactory:      registryFactoryMock,
		eventStorageExecutor: storageExecMock,
		eventParsingExecutor: parsingExecMock,
	}

	testMeta := &types.Metadata{}

	eventRegistry := registry.EventRegistry(map[types.EventID]*registry.Type{})

	blockHash := types.NewHash([]byte{0, 1, 2, 3})

	stateProviderMock.On("GetMetadata", blockHash).
		Return(testMeta, nil).
		Once()

	registryFactoryMock.On("CreateEventRegistry", testMeta).
		Return(eventRegistry, nil).
		Once()

	err := eventRetriever.updateInternalState(&blockHash)
	assert.NoError(t, err)
	assert.Equal(t, testMeta, eventRetriever.meta)
	assert.Equal(t, eventRegistry, eventRetriever.eventRegistry)

	latestMeta := &types.Metadata{}

	stateProviderMock.On("GetLatestMetadata").
		Return(latestMeta, nil).
		Once()

	registryFactoryMock.On("CreateEventRegistry", latestMeta).
		Return(eventRegistry, nil).
		Once()

	err = eventRetriever.updateInternalState(nil)
	assert.NoError(t, err)
	assert.Equal(t, latestMeta, eventRetriever.meta)
	assert.Equal(t, eventRegistry, eventRetriever.eventRegistry)
}

func TestEventRetriever_updateInternalState_MetadataRetrievalError(t *testing.T) {
	eventParserMock := parser.NewEventParserMock(t)
	stateProviderMock := state.NewProviderMock(t)
	registryFactoryMock := registry.NewFactoryMock(t)
	storageExecMock := exec.NewRetryableExecutorMock[*types.StorageDataRaw](t)
	parsingExecMock := exec.NewRetryableExecutorMock[[]*parser.Event](t)

	eventRetriever := &eventRetriever{
		eventParser:          eventParserMock,
		stateProvider:        stateProviderMock,
		registryFactory:      registryFactoryMock,
		eventStorageExecutor: storageExecMock,
		eventParsingExecutor: parsingExecMock,
	}

	blockHash := types.NewHash([]byte{0, 1, 2, 3})

	metadataRetrievalError := errors.New("error")

	stateProviderMock.On("GetMetadata", blockHash).
		Return(nil, metadataRetrievalError).
		Once()

	err := eventRetriever.updateInternalState(&blockHash)
	assert.ErrorIs(t, err, metadataRetrievalError)

	stateProviderMock.On("GetLatestMetadata").
		Return(nil, metadataRetrievalError).
		Once()

	err = eventRetriever.updateInternalState(nil)
	assert.ErrorIs(t, err, metadataRetrievalError)
}

func TestEventRetriever_updateInternalState_RegistryFactoryError(t *testing.T) {
	eventParserMock := parser.NewEventParserMock(t)
	stateProviderMock := state.NewProviderMock(t)
	registryFactoryMock := registry.NewFactoryMock(t)
	storageExecMock := exec.NewRetryableExecutorMock[*types.StorageDataRaw](t)
	parsingExecMock := exec.NewRetryableExecutorMock[[]*parser.Event](t)

	eventRetriever := &eventRetriever{
		eventParser:          eventParserMock,
		stateProvider:        stateProviderMock,
		registryFactory:      registryFactoryMock,
		eventStorageExecutor: storageExecMock,
		eventParsingExecutor: parsingExecMock,
	}

	testMeta := &types.Metadata{}

	blockHash := types.NewHash([]byte{0, 1, 2, 3})

	stateProviderMock.On("GetMetadata", blockHash).
		Return(testMeta, nil).
		Once()

	registryFactoryError := errors.New("error")

	registryFactoryMock.On("CreateEventRegistry", testMeta).
		Return(nil, registryFactoryError).
		Once()

	err := eventRetriever.updateInternalState(&blockHash)
	assert.ErrorIs(t, err, registryFactoryError)

	latestMeta := &types.Metadata{}

	stateProviderMock.On("GetLatestMetadata").
		Return(latestMeta, nil).
		Once()

	registryFactoryMock.On("CreateEventRegistry", latestMeta).
		Return(nil, registryFactoryError).
		Once()

	err = eventRetriever.updateInternalState(nil)
	assert.ErrorIs(t, err, registryFactoryError)
}
