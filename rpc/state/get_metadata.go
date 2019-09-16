// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
// Copyright (C) 2019  Philip Stanislaus, Philip Stehlik, Vimukthi Wickramasinghe
//
// This file is part of Go Substrate RPC Client (GSRPC).
//
// GSRPC is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// GSRPC is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
		hexHash, err := types.Hex(*blockHash)
		if err != nil {
			return metadata, err
		}
		err = (*s.client).Call(&res, "state_getMetadata", hexHash)
		if err != nil {
			return metadata, err
		}
	}
	if err != nil {
		return metadata, err
	}

	err = types.DecodeFromHexString(res, metadata)
	return metadata, err
}
