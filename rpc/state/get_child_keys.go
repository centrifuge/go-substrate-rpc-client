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

// GetChildKeys retreives the keys with the given prefix of a specific child storage
func (s *State) GetChildKeys(childStorageKey, prefix types.StorageKey, blockHash types.Hash) (
	[]types.StorageKey, error) {
	return s.getChildKeys(childStorageKey, prefix, &blockHash)
}

// GetChildKeysLatest retreives the keys with the given prefix of a specific child storage for the latest block height
func (s *State) GetChildKeysLatest(childStorageKey, prefix types.StorageKey) ([]types.StorageKey, error) {
	return s.getChildKeys(childStorageKey, prefix, nil)
}

func (s *State) getChildKeys(childStorageKey, prefix types.StorageKey, blockHash *types.Hash) (
	[]types.StorageKey, error) {
	var res []string
	err := client.CallWithBlockHash(*s.client, &res, "state_getChildKeys", blockHash, childStorageKey.Hex(), prefix.Hex())
	if err != nil {
		return nil, err
	}

	keys := make([]types.StorageKey, len(res))
	for i, r := range res {
		err = types.DecodeFromHexString(r, &keys[i])
		if err != nil {
			return nil, err
		}
	}
	return keys, err
}
