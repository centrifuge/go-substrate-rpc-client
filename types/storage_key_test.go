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

func TestCreateStorageKeyPlainV10(t *testing.T) {
	m := ExamplaryMetadataV10

	key, err := CreateStorageKey(m, "Timestamp", "Now", nil, nil)
	assert.NoError(t, err)
	hex, err := Hex(key)
	assert.NoError(t, err)
	assert.Equal(t, "0xf0c365c3cf59d671eb72da0e7a4113c49f1f0515f462cdcf84e0f1d6045dfcbb", hex)
}

func TestCreateStorageKeyPlainV9(t *testing.T) {
	m := ExamplaryMetadataV9

	key, err := CreateStorageKey(m, "Timestamp", "Now", nil, nil)
	assert.NoError(t, err)
	hex, err := Hex(key)
	assert.NoError(t, err)
	assert.Equal(t, "0xf0c365c3cf59d671eb72da0e7a4113c49f1f0515f462cdcf84e0f1d6045dfcbb", hex)
}

func TestCreateStorageKeyPlainV4(t *testing.T) {
	m := ExamplaryMetadataV4

	key, err := CreateStorageKey(m, "Timestamp", "Now", nil, nil)
	assert.NoError(t, err)
	hex, err := Hex(key)
	assert.NoError(t, err)
	assert.Equal(t, "0x0e4944cfd98d6f4cc374d16f5a4e3f9c", hex)
}

func TestCreateStorageKeyMapV10(t *testing.T) {
	b := MustHexDecodeString(AlicePubKey)
	m := ExamplaryMetadataV10
	key, err := CreateStorageKey(m, "System", "AccountNonce", b, nil)
	assert.NoError(t, err)
	hex, err := Hex(key)
	assert.NoError(t, err)
	assert.Equal(t, "0x26aa394eea5630e07c48ae0c9558cef79c2f82b23e5fd031fb54c292794b4cc42e3fb4c297a84c5cebc0e78257d213d0927ccc7596044c6ba013dd05522aacba", hex) //nolint:lll
}

func TestCreateStorageKeyMapV9(t *testing.T) {
	b := MustHexDecodeString(AlicePubKey)
	m := ExamplaryMetadataV9
	key, err := CreateStorageKey(m, "System", "AccountNonce", b, nil)
	assert.NoError(t, err)
	hex, err := Hex(key)
	assert.NoError(t, err)
	assert.Equal(t, "0x26aa394eea5630e07c48ae0c9558cef79c2f82b23e5fd031fb54c292794b4cc42e3fb4c297a84c5cebc0e78257d213d0927ccc7596044c6ba013dd05522aacba", hex) //nolint:lll
}

func TestCreateStorageKeyMapV4(t *testing.T) {
	b := MustHexDecodeString(AlicePubKey)
	m := ExamplaryMetadataV4
	key, err := CreateStorageKey(m, "System", "AccountNonce", b, nil)
	assert.NoError(t, err)
	hex, err := Hex(key)
	assert.NoError(t, err)
	assert.Equal(t, "0x5c54163a1c72509b5250f0a30b9001fdee9d9b48388b06921f1b210e81e3a1f0", hex)
}

func TestCreateStorageKeyDoubleMapV10(t *testing.T) {
	m := ExamplaryMetadataV10
	key, err := CreateStorageKey(m, "Session", "NextKeys",
		[]byte{0x34, 0x3a, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x3a, 0x6b, 0x65, 0x79, 0x73},
		[]byte{0xbe, 0x5d, 0xdb, 0x15, 0x79, 0xb7, 0x2e, 0x84, 0x52, 0x4f, 0xc2, 0x9e, 0x78, 0x60, 0x9e, 0x3c,
			0xaf, 0x42, 0xe8, 0x5a, 0xa1, 0x18, 0xeb, 0xfe, 0x0b, 0x0a, 0xd4, 0x04, 0xb5, 0xbd, 0xd2, 0x5f},
	)
	assert.NoError(t, err)
	hex, err := Hex(key)
	assert.NoError(t, err)
	assert.Equal(t, "0x"+
		"cec5070d609dd3497f72bde07fc96ba0"+ // twox 128
		"4c014e6bf8b8c2c011e7290b85696bb3"+ // twox 128
		"9fe6329cc0b39e09"+ // twox 64
		"343a73657373696f6e3a6b657973"+ // twox 64 (concat, with length)
		"4724e5390fcf0d08afc9608ff4c45df257266ae599ac7a32baba26155dcf4402", // blake2
		hex) //nolint:lll
}

func TestCreateStorageKeyDoubleMapV9(t *testing.T) {
	m := ExamplaryMetadataV9
	key, err := CreateStorageKey(m, "Session", "NextKeys",
		[]byte{0x34, 0x3a, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x3a, 0x6b, 0x65, 0x79, 0x73},
		[]byte{0xbe, 0x5d, 0xdb, 0x15, 0x79, 0xb7, 0x2e, 0x84, 0x52, 0x4f, 0xc2, 0x9e, 0x78, 0x60, 0x9e, 0x3c,
			0xaf, 0x42, 0xe8, 0x5a, 0xa1, 0x18, 0xeb, 0xfe, 0x0b, 0x0a, 0xd4, 0x04, 0xb5, 0xbd, 0xd2, 0x5f},
	)
	assert.NoError(t, err)
	hex, err := Hex(key)
	assert.NoError(t, err)
	assert.Equal(t, "0x"+
		"cec5070d609dd3497f72bde07fc96ba0"+ // twox 128
		"4c014e6bf8b8c2c011e7290b85696bb3"+ // twox 128
		"9fe6329cc0b39e09"+ // twox 64
		"343a73657373696f6e3a6b657973"+ // twox 64 (concat, with length)
		"4724e5390fcf0d08afc9608ff4c45df257266ae599ac7a32baba26155dcf4402", // blake2
		hex) //nolint:lll
}

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
