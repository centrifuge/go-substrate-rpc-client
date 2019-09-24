// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
// Copyright (C) 2019  Centrifuge GmbH
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
	"encoding/hex"

	"github.com/centrifuge/go-substrate-rpc-client/types"
)

// GetStorage retreives the stored data and decodes them into the provided interface
func (s *State) GetStorage(key types.StorageKey, target interface{}, blockHash types.Hash) error {
	raw, err := s.getStorageRaw(key, &blockHash)
	if err != nil {
		return err
	}
	return types.DecodeFromBytes(*raw, target)
}

// GetStorageLatest retreives the stored data for the latest block height and decodes them into the provided interface
func (s *State) GetStorageLatest(key types.StorageKey, target interface{}) error {
	raw, err := s.getStorageRaw(key, nil)
	if err != nil {
		return err
	}
	return types.DecodeFromBytes(*raw, target)
}

// GetStorageRaw retreives the stored data as raw bytes, without decoding them
func (s *State) GetStorageRaw(key types.StorageKey, blockHash types.Hash) (*types.StorageDataRaw, error) {
	return s.getStorageRaw(key, &blockHash)
}

// GetStorageRawLatest retreives the stored data for the latest block height as raw bytes, without decoding them
func (s *State) GetStorageRawLatest(key types.StorageKey) (*types.StorageDataRaw, error) {
	return s.getStorageRaw(key, nil)
}

func (s *State) getStorageRaw(key types.StorageKey, blockHash *types.Hash) (*types.StorageDataRaw, error) {
	var res string
	var err error
	if blockHash == nil {
		err = (*s.client).Call(&res, "state_getStorage", key.Hex())
	} else {
		hexHash, err := types.Hex(*blockHash)
		if err != nil {
			return nil, err
		}
		err = (*s.client).Call(&res, "state_getStorage", key.Hex(), hexHash)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}

	bz, err := hex.DecodeString(res[2:])
	if err != nil {
		return nil, err
	}

	data := types.NewStorageDataRaw(bz)
	return &data, nil
}
