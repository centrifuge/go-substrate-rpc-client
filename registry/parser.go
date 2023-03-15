package registry

import (
	"bytes"
	"errors"
	"fmt"

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

type Parser interface {
	GetEvents(blockHash types.Hash) ([]*Event, error)
}

type parser struct {
	stateProvider   StateProvider ``
	registryFactory Factory

	eventRegistry EventRegistry
	meta          *types.Metadata
}

func NewParser(stateProvider StateProvider, registryFactory Factory) (Parser, error) {
	parser := &parser{
		stateProvider:   stateProvider,
		registryFactory: registryFactory,
	}

	if err := parser.updateInternalState(nil); err != nil {
		return nil, err
	}

	return parser, nil
}

func (p *parser) GetEvents(blockHash types.Hash) ([]*Event, error) {
	storageEvents, err := p.stateProvider.GetStorageEvents(p.meta, blockHash)

	if err != nil {
		return nil, fmt.Errorf("couldn't retrieve events from storage: %w", err)
	}

	events, err := p.parseEvents(storageEvents)

	if err == nil {
		return events, nil
	}

	if updateErr := p.updateInternalState(&blockHash); updateErr != nil {
		return nil, fmt.Errorf("couldn't update internal state: %w", err)
	}

	return p.parseEvents(storageEvents)
}

func (p *parser) parseEvents(sd *types.StorageDataRaw) ([]*Event, error) {
	decoder := scale.NewDecoder(bytes.NewReader(*sd))

	eventsCount, err := decoder.DecodeUintCompact()

	if err != nil {
		return nil, fmt.Errorf("couldn't get number of events: %w", err)
	}

	var events []*Event

	for i := uint64(0); i < eventsCount.Uint64(); i++ {
		var phase types.Phase

		if err := decoder.Decode(&phase); err != nil {
			return nil, fmt.Errorf("couldn't decode Phase for event #%v: %w", i, err)
		}

		var eventID types.EventID

		if err := decoder.Decode(&eventID); err != nil {
			return nil, fmt.Errorf("couldn't decode EventID for event #%v: %w", i, err)
		}

		eventDecoder, ok := p.eventRegistry[eventID]

		if !ok {
			return nil, errors.New("couldn't get event decoder")
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

func (p *parser) updateInternalState(blockHash *types.Hash) error {
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
