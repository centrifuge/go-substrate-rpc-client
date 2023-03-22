package state

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/rpc/state"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

//go:generate mockery --name Provider --structname ProviderMock --filename provider_mock.go --inpackage

type Provider interface {
	GetMetadata(blockHash types.Hash) (*types.Metadata, error)
	GetLatestMetadata() (*types.Metadata, error)

	GetStorageEvents(meta *types.Metadata, blockHash types.Hash) (*types.StorageDataRaw, error)
}

type provider struct {
	stateRPC state.State
}

func NewProvider(stateRPC state.State) Provider {
	return &provider{stateRPC: stateRPC}
}

func (s *provider) GetLatestMetadata() (*types.Metadata, error) {
	return s.stateRPC.GetMetadataLatest()
}

func (s *provider) GetMetadata(blockHash types.Hash) (*types.Metadata, error) {
	return s.stateRPC.GetMetadata(blockHash)
}

const (
	storagePrefix = "System"
	storageMethod = "Events"
)

func (s *provider) GetStorageEvents(meta *types.Metadata, blockHash types.Hash) (*types.StorageDataRaw, error) {
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
