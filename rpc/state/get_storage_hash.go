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

// GetStorageHash retreives the keys with the given key
func (s *State) GetStorageHash(key types.StorageKey, blockHash types.Hash) (types.Hash, error) {
	return s.getStorageHash(key, &blockHash)
}

// GetStorageHashLatest retreives the keys with the given key for the latest block height
func (s *State) GetStorageHashLatest(key types.StorageKey) (types.Hash, error) {
	return s.getStorageHash(key, nil)
}

func (s *State) getStorageHash(key types.StorageKey, blockHash *types.Hash) (types.Hash, error) {
	var res string
	err := client.CallWithBlockHash(*s.client, &res, "state_getStorageHash", blockHash, key.Hex())
	if err != nil {
		return types.Hash{}, err
	}

	return types.NewHashFromHexString(res)
}
