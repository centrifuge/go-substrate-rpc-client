// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
//
// Copyright 2019 Centrifuge GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types_test

import (
	"testing"

	. "github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/stretchr/testify/assert"
)

const (
	AlicePubKey = "0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d"
)

func TestCreateStorageKey(t *testing.T) {
	m := ExamplaryMetadataV4

	key, err := CreateStorageKey(m, "Timestamp", "Now", nil)
	assert.NoError(t, err)
	hex, err := Hex(key)
	assert.NoError(t, err)
	assert.Equal(t, "0x0e4944cfd98d6f4cc374d16f5a4e3f9c", hex)
}

func TestCreateStorageKey2(t *testing.T) {
	b := MustHexDecodeString(AlicePubKey)
	m := ExamplaryMetadataV4
	key, err := CreateStorageKey(m, "System", "AccountNonce", b)
	assert.NoError(t, err)
	hex, err := Hex(key)
	assert.NoError(t, err)
	assert.Equal(t, "0x5c54163a1c72509b5250f0a30b9001fdee9d9b48388b06921f1b210e81e3a1f0", hex)
}

func TestCreateStorageKey_MetadataV4(t *testing.T) {
	b := MustHexDecodeString(AlicePubKey)
	m := ExamplaryMetadataV4
	key, err := CreateStorageKey(m, "Balances", "FreeBalance", b)
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
// 	enc, err := EncodeToStorageKey(k)
// 	assert.NoError(t, err)
// 	key, err := CreateStorageKey(m, "System", "EventTopics", enc)
// 	assert.NoError(t, err)
// 	hex, err := Hex(key)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "0x7f864e18e3dd8b58386310d2fe0919eef27c6e558564b7f67f22d99d20f587bb", hex)
// }

func TestStorageKey_EncodedLength(t *testing.T) {
	assertEncodedLength(t, []encodedLengthAssert{
		{NewStorageKey(MustHexDecodeString("0x00")), 1},
		{NewStorageKey(MustHexDecodeString("0xab1234")), 3},
		{NewStorageKey(MustHexDecodeString("0x0001")), 2},
	})
}

func TestStorageKey_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{NewStorageKey([]byte{171, 18, 52}), MustHexDecodeString("0xab1234")},
		{NewStorageKey([]byte{}), MustHexDecodeString("0x")},
	})
}

func TestStorageKey_Decode(t *testing.T) {
	bz := []byte{12, 251, 42}
	decoded := make(StorageKey, len(bz))
	err := DecodeFromBytes(bz, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, StorageKey(bz), decoded)
}

func TestStorageKey_Hash(t *testing.T) {
	assertHash(t, []hashAssert{
		{NewStorageKey([]byte{0, 42, 254}), MustHexDecodeString(
			"0x537db36f5b5970b679a28a3df8d219317d658014fb9c3d409c0c799d8ecf149d")},
		{NewStorageKey([]byte{0, 0}), MustHexDecodeString(
			"0x9ee6dfb61a2fb903df487c401663825643bb825d41695e63df8af6162ab145a6")},
	})
}

func TestStorageKey_Hex(t *testing.T) {
	assertEncodeToHex(t, []encodeToHexAssert{
		{NewStorageKey([]byte{0, 0, 0}), "0x000000"},
		{NewStorageKey([]byte{171, 18, 52}), "0xab1234"},
		{NewStorageKey([]byte{0, 1}), "0x0001"},
		{NewStorageKey([]byte{18, 52, 86}), "0x123456"},
	})
}

func TestStorageKey_String(t *testing.T) {
	assertString(t, []stringAssert{
		{NewStorageKey([]byte{0, 0, 0}), "[0 0 0]"},
		{NewStorageKey([]byte{171, 18, 52}), "[171 18 52]"},
		{NewStorageKey([]byte{0, 1}), "[0 1]"},
		{NewStorageKey([]byte{18, 52, 86}), "[18 52 86]"},
	})
}

func TestStorageKey_Eq(t *testing.T) {
	assertEq(t, []eqAssert{
		{NewStorageKey([]byte{1, 0, 0}), NewStorageKey([]byte{1, 0}), false},
		{NewStorageKey([]byte{0, 0, 1}), NewStorageKey([]byte{0, 1}), false},
		{NewStorageKey([]byte{0, 0, 0}), NewStorageKey([]byte{0, 0}), false},
		{NewStorageKey([]byte{12, 48, 255}), NewStorageKey([]byte{12, 48, 255}), true},
		{NewStorageKey([]byte{0}), NewStorageKey([]byte{0}), true},
		{NewStorageKey([]byte{1}), NewBool(true), false},
		{NewStorageKey([]byte{0}), NewBool(false), false},
	})
}
