package state

import (
	"bytes"
	"encoding/hex"
	"github.com/centrifuge/go-substrate-rpc-client/client"
	"github.com/centrifuge/go-substrate-rpc-client/scale"
	"github.com/centrifuge/go-substrate-rpc-client/types"
)

type State struct {
	client *client.Client
}

func NewState(c *client.Client) *State {
	return &State{c}
}

func (s *State) GetMetadata(blockHash types.Hash) (*types.Metadata, error) {
	return s.getMetadata(&blockHash)
}

func (s *State) 	GetMetadataLatest() (*types.Metadata, error) {
	return s.getMetadata(nil)
}


func (s *State) getMetadata(blockHash *types.Hash) (*types.Metadata, error) {
	var res string
	err := (*s.client).Call(&res, "state_getMetadata")
	if err != nil {
		return types.NewMetadata(), err
	}

	bz, err := hex.DecodeString(res[2:])
	if err != nil {
		return types.NewMetadata(), err
	}

	decoder := scale.NewDecoder(bytes.NewReader(bz))

	metadata := types.NewMetadata()

	err = decoder.Decode(metadata)
	if err != nil {
		return types.NewMetadata(), err
	}

	return metadata, err
}
