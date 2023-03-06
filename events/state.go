package events

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/rpc/state"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

type StateProvider interface {
	GetMetadata(blockHash types.Hash) (*types.Metadata, error)
	GetLatestMetadata() (*types.Metadata, error)

	GetStorageEvents(meta *types.Metadata, blockHash types.Hash) (*types.StorageDataRaw, error)
}

type stateProvider struct {
	stateRPC state.State
}

func NewStateProvider(stateRPC state.State) StateProvider {
	return &stateProvider{stateRPC: stateRPC}
}

func (s *stateProvider) GetLatestMetadata() (*types.Metadata, error) {
	return s.stateRPC.GetMetadataLatest()
}

func (s *stateProvider) GetMetadata(blockHash types.Hash) (*types.Metadata, error) {
	return s.stateRPC.GetMetadata(blockHash)
}

const (
	storagePrefix = "System"
	storageMethod = "Events"
)

func (s *stateProvider) GetStorageEvents(meta *types.Metadata, blockHash types.Hash) (*types.StorageDataRaw, error) {
	key, err := types.CreateStorageKey(meta, storagePrefix, storageMethod, nil)

	if err != nil {
		return nil, err
	}

	storageData, err := s.stateRPC.GetStorageRaw(key, blockHash)

	if err != nil {
		return nil, err
	}

	return storageData, nil
}
