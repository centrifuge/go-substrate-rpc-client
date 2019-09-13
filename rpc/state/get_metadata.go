package state

import (
	"github.com/centrifuge/go-substrate-rpc-client/types"
)

func (s *State) GetMetadata(blockHash types.Hash) (*types.Metadata, error) {
	return s.getMetadata(&blockHash)
}

func (s *State) GetMetadataLatest() (*types.Metadata, error) {
	return s.getMetadata(nil)
}

func (s *State) getMetadata(blockHash *types.Hash) (*types.Metadata, error) {
	metadata := types.NewMetadata()

	var res string
	var err error
	if blockHash == nil {
		err = (*s.client).Call(&res, "state_getMetadata")
	} else {
		err = (*s.client).Call(&res, "state_getMetadata", *blockHash)
	}
	if err != nil {
		return metadata, err
	}

	err = types.DecodeFromHexString(res, metadata)
	return metadata, err
}
