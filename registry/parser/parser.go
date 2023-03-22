package parser

import (
	"bytes"
	"fmt"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/registry"
	"github.com/centrifuge/go-substrate-rpc-client/v4/registry/exec"
	"github.com/centrifuge/go-substrate-rpc-client/v4/registry/state"
	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

type Event struct {
	Name    string
	Fields  map[string]any
	EventID types.EventID
	Phase   *types.Phase
	Topics  []types.Hash
}

type EventParser interface {
	GetEvents(blockHash types.Hash) ([]*Event, error)
	ParseEvents(sd *types.StorageDataRaw) ([]*Event, error)
}

type eventParser struct {
	stateProvider   state.Provider
	registryFactory registry.Factory

	eventStorageExecutor exec.RetryableExecutor[*types.StorageDataRaw]
	eventParsingExecutor exec.RetryableExecutor[[]*Event]

	eventRegistry registry.EventRegistry
	meta          *types.Metadata
}

func NewEventParser(
	stateProvider state.Provider,
	registryFactory registry.Factory,
	eventStorageExecutor exec.RetryableExecutor[*types.StorageDataRaw],
	eventParsingExecutor exec.RetryableExecutor[[]*Event],
) (EventParser, error) {
	parser := &eventParser{
		stateProvider:        stateProvider,
		registryFactory:      registryFactory,
		eventStorageExecutor: eventStorageExecutor,
		eventParsingExecutor: eventParsingExecutor,
	}

	if err := parser.updateInternalState(nil); err != nil {
		return nil, err
	}

	return parser, nil
}

func NewDefaultEventParser(stateProvider state.Provider, registryFactory registry.Factory) (EventParser, error) {
	eventStorageExecutor := exec.NewRetryableExecutor[*types.StorageDataRaw](exec.WithErrTimeout(1 * time.Second))
	eventParsingExecutor := exec.NewRetryableExecutor[[]*Event](exec.WithMaxRetryCount(1))

	return NewEventParser(stateProvider, registryFactory, eventStorageExecutor, eventParsingExecutor)
}

func (p *eventParser) GetEvents(blockHash types.Hash) ([]*Event, error) {
	storageEvents, err := p.eventStorageExecutor.ExecWithFallback(
		func() (*types.StorageDataRaw, error) {
			return p.stateProvider.GetStorageEvents(p.meta, blockHash)
		},
		func() error {
			return p.updateInternalState(&blockHash)
		},
	)

	if err != nil {
		return nil, fmt.Errorf("couldn't retrieve raw events from storage: %w", err)
	}

	events, err := p.eventParsingExecutor.ExecWithFallback(
		func() ([]*Event, error) {
			return p.ParseEvents(storageEvents)
		},
		func() error {
			return p.updateInternalState(&blockHash)
		},
	)

	if err != nil {
		return nil, fmt.Errorf("couldn't parse events: %w", err)
	}

	return events, nil
}

func (p *eventParser) ParseEvents(sd *types.StorageDataRaw) ([]*Event, error) {
	decoder := scale.NewDecoder(bytes.NewReader(*sd))

	eventsCount, err := decoder.DecodeUintCompact()

	if err != nil {
		return nil, fmt.Errorf("couldn't get number of events: %w", err)
	}

	var events []*Event

	for i := uint64(0); i < eventsCount.Uint64(); i++ {
		var phase types.Phase

		if err := decoder.Decode(&phase); err != nil {
			return nil, fmt.Errorf("couldn't decode Phase for event #%d: %w", i, err)
		}

		var eventID types.EventID

		if err := decoder.Decode(&eventID); err != nil {
			return nil, fmt.Errorf("couldn't decode event ID for event #%d: %w", i, err)
		}

		eventDecoder, ok := p.eventRegistry[eventID]

		if !ok {
			return nil, fmt.Errorf("couldn't find decoder for event #%d with ID: %v", i, eventID)
		}

		eventFields, err := eventDecoder.Decode(decoder)

		if err != nil {
			return nil, fmt.Errorf("couldn't decode event fields: %w", err)
		}

		var topics []types.Hash

		if err := decoder.Decode(&topics); err != nil {
			return nil, fmt.Errorf("unable to decode topics for event #%v: %w", i, err)
		}

		event := &Event{
			Name:    eventDecoder.Name,
			Fields:  eventFields,
			EventID: eventID,
			Phase:   &phase,
			Topics:  topics,
		}

		events = append(events, event)
	}

	return events, nil
}

func (p *eventParser) updateInternalState(blockHash *types.Hash) error {
	var (
		meta *types.Metadata
		err  error
	)

	if blockHash == nil {
		meta, err = p.stateProvider.GetLatestMetadata()
	} else {
		meta, err = p.stateProvider.GetMetadata(*blockHash)
	}

	if err != nil {
		return fmt.Errorf("couldn't retrieve metadata: %w", err)
	}

	eventRegistry, err := p.registryFactory.CreateEventRegistry(meta)

	if err != nil {
		return fmt.Errorf("couldn't create event registry: %w", err)
	}

	p.meta = meta
	p.eventRegistry = eventRegistry

	return nil
}
