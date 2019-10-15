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
	"github.com/centrifuge/go-substrate-rpc-client/client"
	"github.com/centrifuge/go-substrate-rpc-client/types"
)

// GetChildStorage retreives the child storage for a key and decodes them into the provided interface
func (s *State) GetChildStorage(childStorageKey, key types.StorageKey, target interface{}, blockHash types.Hash) error {
	raw, err := s.getChildStorageRaw(childStorageKey, key, &blockHash)
	if err != nil {
		return err
	}
	return types.DecodeFromBytes(*raw, target)
}

// GetChildStorageLatest retreives the child storage for a key for the latest block height and decodes them into the
// provided interface
func (s *State) GetChildStorageLatest(childStorageKey, key types.StorageKey, target interface{}) error {
	raw, err := s.getChildStorageRaw(childStorageKey, key, nil)
	if err != nil {
		return err
	}
	return types.DecodeFromBytes(*raw, target)
}

// GetChildStorageRaw retreives the child storage for a key as raw bytes, without decoding them
func (s *State) GetChildStorageRaw(childStorageKey, key types.StorageKey, blockHash types.Hash) (
	*types.StorageDataRaw, error) {
	return s.getChildStorageRaw(childStorageKey, key, &blockHash)
}

// GetChildStorageRawLatest retreives the child storage for a key for the latest block height as raw bytes,
// without decoding them
func (s *State) GetChildStorageRawLatest(childStorageKey, key types.StorageKey) (*types.StorageDataRaw, error) {
	return s.getChildStorageRaw(childStorageKey, key, nil)
}

func (s *State) getChildStorageRaw(childStorageKey, key types.StorageKey, blockHash *types.Hash) (
	*types.StorageDataRaw, error) {
	var res string
	err := client.CallWithBlockHash(*s.client, &res, "state_getChildStorage", blockHash, childStorageKey.Hex(),
		key.Hex())
	if err != nil {
		return nil, err
	}

	bz, err := types.HexDecodeString(res)
	if err != nil {
		return nil, err
	}

	data := types.NewStorageDataRaw(bz)
	return &data, nil
}
