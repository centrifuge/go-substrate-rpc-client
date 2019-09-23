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

package types_test

import (
	"testing"

	. "github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
)

const (
	AlicePubKey = "0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d"
)

func TestStorageKey(t *testing.T) {
	m := ExamplaryMetadataV4

	key, err := NewStorageKey(m, "Timestamp", "Now", nil)
	assert.NoError(t, err)
	assert.Equal(t, StorageKey{0xe, 0x49, 0x44, 0xcf, 0xd9, 0x8d, 0x6f, 0x4c, 0xc3, 0x74, 0xd1, 0x6f, 0x5a, 0x4e, 0x3f, 0x9c}, key) //nolint:lll
}

func TestStorageKey2(t *testing.T) {
	b, _ := hexutil.Decode(AlicePubKey)
	m := ExamplaryMetadataV4
	key, err := NewStorageKey(m, "System", "AccountNonce", b)
	assert.NoError(t, err)
	assert.Equal(t, StorageKey{0x5c, 0x54, 0x16, 0x3a, 0x1c, 0x72, 0x50, 0x9b, 0x52, 0x50, 0xf0, 0xa3, 0xb, 0x90, 0x1, 0xfd, 0xee, 0x9d, 0x9b, 0x48, 0x38, 0x8b, 0x6, 0x92, 0x1f, 0x1b, 0x21, 0xe, 0x81, 0xe3, 0xa1, 0xf0}, key) //nolint:lll
}

func TestStorageKey_MetadataV4(t *testing.T) {
	b, _ := hexutil.Decode(AlicePubKey)
	m := ExamplaryMetadataV4
	key, err := NewStorageKey(m, "Balances", "FreeBalance", b)
	assert.NoError(t, err)
	hex, err := Hex(key)
	assert.NoError(t, err)
	assert.Equal(t, "0x7f864e18e3dd8b58386310d2fe0919eef27c6e558564b7f67f22d99d20f587bb", hex)
}

// TODO: add
// func TestStorageKey_MetadataV4_DoubleMap(t *testing.T) {
// 	// k := struct{
// 	// 	A string
// 	// 	b int[]
// 	// }{
// 	// 	"any",
// 	// 	[]int{0, 1, 2}
// 	// }

// 	k := struct{ A string }{"any "}
// 	m := ExamplaryMetadataV4
// 	enc, err := EncodeToBytes(k)
// 	assert.NoError(t, err)
// 	key, err := NewStorageKey(m, "System", "EventTopics", enc)
// 	assert.NoError(t, err)
// 	hex, err := Hex(key)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "0x7f864e18e3dd8b58386310d2fe0919eef27c6e558564b7f67f22d99d20f587bb", hex)
// }
