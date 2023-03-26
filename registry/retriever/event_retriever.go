package retriever

import (
	"fmt"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/registry"
	"github.com/centrifuge/go-substrate-rpc-client/v4/registry/exec"
	"github.com/centrifuge/go-substrate-rpc-client/v4/registry/parser"
	"github.com/centrifuge/go-substrate-rpc-client/v4/registry/state"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

//go:generate mockery --name EventRetriever --structname EventRetrieverMock --filename event_retriever_mock.go --inpackage

type EventRetriever interface {
	GetEvents(blockHash types.Hash) ([]*parser.Event, error)
}

type eventRetriever struct {
	eventParser parser.EventParser

	stateProvider   state.Provider
	registryFactory registry.Factory

	eventStorageExecutor exec.RetryableExecutor[*types.StorageDataRaw]
	eventParsingExecutor exec.RetryableExecutor[[]*parser.Event]

	eventRegistry registry.EventRegistry
	meta          *types.Metadata
}

func NewEventRetriever(
	eventParser parser.EventParser,
	stateProvider state.Provider,
	registryFactory registry.Factory,
	eventStorageExecutor exec.RetryableExecutor[*types.StorageDataRaw],
	eventParsingExecutor exec.RetryableExecutor[[]*parser.Event],
) (EventRetriever, error) {
	retriever := &eventRetriever{
		eventParser:          eventParser,
		stateProvider:        stateProvider,
		registryFactory:      registryFactory,
		eventStorageExecutor: eventStorageExecutor,
		eventParsingExecutor: eventParsingExecutor,
	}

	if err := retriever.updateInternalState(nil); err != nil {
		return nil, err
	}

	return retriever, nil
}

func NewDefaultEventRetriever(stateProvider state.Provider) (EventRetriever, error) {
	eventParser := parser.NewEventParser()
	registryFactory := registry.NewFactory()

	eventStorageExecutor := exec.NewRetryableExecutor[*types.StorageDataRaw](exec.WithRetryTimeout(1 * time.Second))
	eventParsingExecutor := exec.NewRetryableExecutor[[]*parser.Event](exec.WithMaxRetryCount(1))

	return NewEventRetriever(eventParser, stateProvider, registryFactory, eventStorageExecutor, eventParsingExecutor)
}

func (e *eventRetriever) GetEvents(blockHash types.Hash) ([]*parser.Event, error) {
	storageEvents, err := e.eventStorageExecutor.ExecWithFallback(
		func() (*types.StorageDataRaw, error) {
			return e.stateProvider.GetStorageEvents(e.meta, blockHash)
		},
		func() error {
			return e.updateInternalState(&blockHash)
		},
	)

	if err != nil {
		return nil, fmt.Errorf("couldn't retrieve raw events from storage: %w", err)
	}

	events, err := e.eventParsingExecutor.ExecWithFallback(
		func() ([]*parser.Event, error) {
			return e.eventParser.ParseEvents(e.eventRegistry, storageEvents)
		},
		func() error {
			return e.updateInternalState(&blockHash)
		},
	)

	if err != nil {
		return nil, fmt.Errorf("couldn't parse events: %w", err)
	}

	return events, nil
}

func (e *eventRetriever) updateInternalState(blockHash *types.Hash) error {
	var (
		meta *types.Metadata
		err  error
	)

	if blockHash == nil {
		meta, err = e.stateProvider.GetLatestMetadata()
	} else {
		meta, err = e.stateProvider.GetMetadata(*blockHash)
	}

	if err != nil {
		return fmt.Errorf("couldn't retrieve metadata: %w", err)
	}

	eventRegistry, err := e.registryFactory.CreateEventRegistry(meta)

	if err != nil {
		return fmt.Errorf("couldn't create event registry: %w", err)
	}

	e.meta = meta
	e.eventRegistry = eventRegistry

	return nil
}
