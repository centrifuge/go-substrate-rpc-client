package parser

import (
	"bytes"
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/registry"
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

//go:generate mockery --name EventParser --structname EventParserMock --filename event_parser_mock.go --inpackage

type EventParser interface {
	ParseEvents(eventRegistry registry.EventRegistry, sd *types.StorageDataRaw) ([]*Event, error)
}

type EventParserFn func(eventRegistry registry.EventRegistry, sd *types.StorageDataRaw) ([]*Event, error)

func (f EventParserFn) ParseEvents(eventRegistry registry.EventRegistry, sd *types.StorageDataRaw) ([]*Event, error) {
	return f(eventRegistry, sd)
}

func NewEventParser() EventParser {
	return EventParserFn(func(eventRegistry registry.EventRegistry, sd *types.StorageDataRaw) ([]*Event, error) {
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

			eventDecoder, ok := eventRegistry[eventID]

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
	})
}
