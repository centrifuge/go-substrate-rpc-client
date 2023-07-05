package state

import (
	libErr "github.com/centrifuge/go-substrate-rpc-client/v4/error"
	"github.com/centrifuge/go-substrate-rpc-client/v4/rpc/state"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

const (
	ErrEventStorageKeyCreation = libErr.Error("event storage key creation")
	ErrEventStorageRetrieval   = libErr.Error("event storage retrieval")
)

//go:generate mockery --name EventProvider --structname EventProviderMock --filename event_provider_mock.go --inpackage

// EventProvider is the interface used for retrieving event data from the storage.
type EventProvider interface {
	GetStorageEvents(meta *types.Metadata, blockHash types.Hash) (*types.StorageDataRaw, error)
}

// eventProvider implements the EventProvider interface.
type eventProvider struct {
	stateRPC state.State
}

// NewEventProvider creates a new EventProvider.
func NewEventProvider(stateRPC state.State) EventProvider {
	return &eventProvider{stateRPC: stateRPC}
}

const (
	storagePrefix = "System"
	storageMethod = "Events"
)

// GetStorageEvents returns the event storage data found at the provided blockHash.
func (p *eventProvider) GetStorageEvents(meta *types.Metadata, blockHash types.Hash) (*types.StorageDataRaw, error) {
	key, err := types.CreateStorageKey(meta, storagePrefix, storageMethod, nil)

	if err != nil {
		return nil, ErrEventStorageKeyCreation.Wrap(err)
	}

	storageData, err := p.stateRPC.GetStorageRaw(key, blockHash)

	if err != nil {
		return nil, ErrEventStorageRetrieval.Wrap(err)
	}

	return storageData, nil
}
