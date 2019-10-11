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
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/stretchr/testify/assert"
)

func TestState_GetStorageLatest(t *testing.T) {
	var decoded types.U64
	err := state.GetStorageLatest(types.MustHexDecodeString(mockSrv.storageKeyHex), &decoded)
	assert.NoError(t, err)
	assert.Equal(t, types.U64(0x5d892db8), decoded)
}

func TestState_GetStorage(t *testing.T) {
	var decoded types.U64
	err := state.GetStorage(types.MustHexDecodeString(mockSrv.storageKeyHex), &decoded, mockSrv.blockHashLatest)
	assert.NoError(t, err)
	assert.Equal(t, types.U64(0x5d892db8), decoded)
}

func TestState_GetStorageRawLatest(t *testing.T) {
	data, err := state.GetStorageRawLatest(types.MustHexDecodeString(mockSrv.storageKeyHex))
	assert.NoError(t, err)
	assert.Equal(t, mockSrv.storageDataHex, data.Hex())
}

func TestState_GetStorageRaw(t *testing.T) {
	data, err := state.GetStorageRaw(types.MustHexDecodeString(mockSrv.storageKeyHex), mockSrv.blockHashLatest)
	assert.NoError(t, err)
	assert.Equal(t, mockSrv.storageDataHex, data.Hex())
}
