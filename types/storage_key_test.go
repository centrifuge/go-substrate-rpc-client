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
	"encoding/binary"
	"fmt"
	"strings"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/hash"
	. "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	. "github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	. "github.com/centrifuge/go-substrate-rpc-client/v4/types/test_utils"
	"github.com/centrifuge/go-substrate-rpc-client/v4/xxhash"
	"github.com/stretchr/testify/assert"
)

const (
	AlicePubKey = "0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d"
)

func TestCreateStorageKeyArgValidationForPlainKey(t *testing.T) {
	for _, m := range []*Metadata{ExamplaryMetadataV13, DecodedMetadataV14Example()} {
		fmt.Println("Testing against metadata v", m.Version)

		_, err := CreateStorageKey(m, "Timestamp", "Now")
		assert.NoError(t, err)

		_, err = CreateStorageKey(m, "Timestamp", "Now", nil)
		assert.NoError(t, err)

		_, err = CreateStorageKey(m, "Timestamp", "Now", nil, []byte{})
		assert.NoError(t, err)

		_, err = CreateStorageKey(m, "Timestamp", "Now", nil, []byte{0x01})
		assert.EqualError(t, err, "non-nil arguments cannot be preceded by nil arguments")

		_, err = CreateStorageKey(m, "Timestamp", "Now", []byte{0x01})
		assert.EqualError(t, err, "Timestamp:Now is a plain key, therefore requires no argument. received: 1")

		expectedKeyBuilder := strings.Builder{}
		hexStr, err := Hex(xxhash.New128([]byte("Timestamp")).Sum(nil))
		assert.NoError(t, err)
		expectedKeyBuilder.WriteString(hexStr)
		hexStr, err = Hex(xxhash.New128([]byte("Now")).Sum(nil))
		assert.NoError(t, err)
		expectedKeyBuilder.WriteString(strings.TrimPrefix(hexStr, "0x"))

		key, err := CreateStorageKey(m, "Timestamp", "Now")
		assert.NoError(t, err)
		hex, err := Hex(key)
		assert.NoError(t, err)
		assert.Equal(t, expectedKeyBuilder.String(), hex)
	}

}

func TestCreateStorageKeyArgValidationForMapKey(t *testing.T) {
	for _, m := range []*Metadata{ExamplaryMetadataV13, DecodedMetadataV14Example()} {
		fmt.Println("Testing against metadata v", m.Version)

		_, err := CreateStorageKey(m, "System", "Account")
		assert.EqualError(
			t,
			err,
			"System:Account is a map, therefore requires that number of arguments should exactly match number of hashers in metadata. Expected: 1, received: 0",
		)

		_, err = CreateStorageKey(m, "System", "Account", nil)
		assert.EqualError(
			t,
			err,
			"System:Account is a map, therefore requires that number of arguments should exactly match number of hashers in metadata. Expected: 1, received: 0",
		)

		_, err = CreateStorageKey(m, "System", "Account", nil, []byte{})
		assert.EqualError(
			t,
			err,
			"System:Account is a map, therefore requires that number of arguments should exactly match number of hashers in metadata. Expected: 1, received: 0",
		)

		_, err = CreateStorageKey(m, "System", "Account", nil, []byte{0x01})
		assert.EqualError(t, err, "non-nil arguments cannot be preceded by nil arguments")

		accountIdSerialized := MustHexDecodeString(AlicePubKey)

		// Build expected answer
		expectedKeyBuilder := strings.Builder{}
		hexStr, err := Hex(xxhash.New128([]byte("System")).Sum(nil))
		assert.NoError(t, err)
		expectedKeyBuilder.WriteString(hexStr)
		hexStr, err = Hex(xxhash.New128([]byte("Account")).Sum(nil))
		assert.NoError(t, err)
		expectedKeyBuilder.WriteString(strings.TrimPrefix(hexStr, "0x"))
		accountIdHasher, err := hash.NewBlake2b128Concat(nil)
		assert.NoError(t, err)
		_, err = accountIdHasher.Write(accountIdSerialized)
		assert.NoError(t, err)
		hexStr, err = Hex(accountIdHasher.Sum(nil))
		assert.NoError(t, err)
		expectedKeyBuilder.WriteString(strings.TrimPrefix(hexStr, "0x"))

		key, err := CreateStorageKey(m, "System", "Account", accountIdSerialized)
		assert.NoError(t, err)
		hex, err := Hex(key)
		assert.NoError(t, err)
		assert.Equal(t, expectedKeyBuilder.String(), hex)
	}
}

func TestCreateStorageKeyArgValidationForDoubleMapKey(t *testing.T) {
	for _, m := range []*Metadata{ExamplaryMetadataV13, DecodedMetadataV14Example()} {
		fmt.Println("Testing against metadata v", m.Version)

		_, err := CreateStorageKey(m, "Staking", "ErasStakers")
		assert.EqualError(t, err, "Staking:ErasStakers is a map, therefore requires that number of arguments should exactly match number of hashers in metadata. Expected: 2, received: 0")

		_, err = CreateStorageKey(m, "Staking", "ErasStakers", nil)
		assert.EqualError(t, err, "Staking:ErasStakers is a map, therefore requires that number of arguments should exactly match number of hashers in metadata. Expected: 2, received: 0")

		_, err = CreateStorageKey(m, "Staking", "ErasStakers", nil, []byte{})
		assert.EqualError(t, err, "Staking:ErasStakers is a map, therefore requires that number of arguments should exactly match number of hashers in metadata. Expected: 2, received: 0")

		_, err = CreateStorageKey(m, "Staking", "ErasStakers", nil, []byte{0x01})
		assert.EqualError(t, err, "non-nil arguments cannot be preceded by nil arguments")

		_, err = CreateStorageKey(m, "Staking", "ErasStakers", []byte{0x01})
		assert.EqualError(t, err, "Staking:ErasStakers is a map, therefore requires that number of arguments should exactly match number of hashers in metadata. Expected: 2, received: 1")

		// Serialize EraIndex and AccountId
		accountIdSerialized := MustHexDecodeString(AlicePubKey)
		var eraIndex uint32 = 3
		eraIndexSerialized := make([]byte, 4)
		binary.LittleEndian.PutUint32(eraIndexSerialized, eraIndex)

		// Build expected answer
		expectedKeyBuilder := strings.Builder{}
		hexStr, err := Hex(xxhash.New128([]byte("Staking")).Sum(nil))
		assert.NoError(t, err)
		expectedKeyBuilder.WriteString(hexStr)
		hexStr, err = Hex(xxhash.New128([]byte("ErasStakers")).Sum(nil))
		assert.NoError(t, err)
		expectedKeyBuilder.WriteString(strings.TrimPrefix(hexStr, "0x"))
		hexStr, err = Hex(xxhash.New64Concat(eraIndexSerialized).Sum(nil))
		assert.NoError(t, err)
		expectedKeyBuilder.WriteString(strings.TrimPrefix(hexStr, "0x"))
		hexStr, err = Hex(xxhash.New64Concat(accountIdSerialized).Sum(nil))
		assert.NoError(t, err)
		expectedKeyBuilder.WriteString(strings.TrimPrefix(hexStr, "0x"))

		key, err := CreateStorageKey(m, "Staking", "ErasStakers", eraIndexSerialized, accountIdSerialized)
		assert.NoError(t, err)
		hex, err := Hex(key)
		assert.NoError(t, err)

		assert.Equal(t, expectedKeyBuilder.String(), hex)
	}
}

func TestCreateStorageKeyArgValidationForNMapKey(t *testing.T) {
	m := ExamplaryMetadataV13
	//"Assets", "Approvals", "AssetId(u32)", "AccountId", "AccountId"

	_, err := CreateStorageKey(m, "Assets", "Approvals")
	assert.EqualError(t, err, "Assets:Approvals is a map, therefore requires that number of arguments "+
		"should exactly match number of hashers in metadata. Expected: 3, received: 0")

	_, err = CreateStorageKey(m, "Assets", "Approvals", nil)
	assert.EqualError(t, err, "Assets:Approvals is a map, therefore requires that number of arguments "+
		"should exactly match number of hashers in metadata. Expected: 3, received: 0")

	_, err = CreateStorageKey(m, "Assets", "Approvals", nil, []byte{})
	assert.EqualError(t, err, "Assets:Approvals is a map, therefore requires that number of arguments "+
		"should exactly match number of hashers in metadata. Expected: 3, received: 0")

	_, err = CreateStorageKey(m, "Assets", "Approvals", nil, []byte{0x01})
	assert.EqualError(t, err, "non-nil arguments cannot be preceded by nil arguments")

	_, err = CreateStorageKey(m, "Assets", "Approvals", []byte{0x01})
	assert.EqualError(t, err, "Assets:Approvals is a map, therefore requires that number of arguments "+
		"should exactly match number of hashers in metadata. Expected: 3, received: 1")

	// Serialize EraIndex and AccountId
	var assetId uint32 = 3
	assetIdSerialized := make([]byte, 4)
	binary.LittleEndian.PutUint32(assetIdSerialized, assetId)
	// Will be used both as owner as well as delegate
	accountIdSerialized := MustHexDecodeString(AlicePubKey)

	// Build expected answer
	expectedKeyBuilder := strings.Builder{}
	hexStr, err := Hex(xxhash.New128([]byte("Assets")).Sum(nil))
	assert.NoError(t, err)
	expectedKeyBuilder.WriteString(hexStr)
	hexStr, err = Hex(xxhash.New128([]byte("Approvals")).Sum(nil))
	assert.NoError(t, err)
	expectedKeyBuilder.WriteString(strings.TrimPrefix(hexStr, "0x"))
	// Hashing serialized asset id
	assetIdHasher, err := hash.NewBlake2b128Concat(nil)
	assert.NoError(t, err)
	_, err = assetIdHasher.Write(assetIdSerialized)
	assert.NoError(t, err)
	hexStr, err = Hex(assetIdHasher.Sum(nil))
	expectedKeyBuilder.WriteString(strings.TrimPrefix(hexStr, "0x"))
	// Hashing serialized account id
	accountIdHasher, err := hash.NewBlake2b128Concat(nil)
	assert.NoError(t, err)
	_, err = accountIdHasher.Write(accountIdSerialized)
	assert.NoError(t, err)
	hexStr, err = Hex(accountIdHasher.Sum(nil))
	// Writing it multiple times as both owner and delegate
	expectedKeyBuilder.WriteString(strings.TrimPrefix(hexStr, "0x"))
	expectedKeyBuilder.WriteString(strings.TrimPrefix(hexStr, "0x"))

	key, err := CreateStorageKey(m, "Assets", "Approvals", assetIdSerialized, accountIdSerialized,
		accountIdSerialized)
	assert.NoError(t, err)
	hex, err := Hex(key)
	assert.NoError(t, err)

	assert.Equal(t, expectedKeyBuilder.String(), hex)
}

func TestCreateStorageKeyPlainV14(t *testing.T) {
	m := DecodedMetadataV14Example()

	key, err := CreateStorageKey(m, "Timestamp", "Now")
	assert.NoError(t, err)
	hex, err := Hex(key)
	assert.NoError(t, err)
	assert.Equal(t, "0xf0c365c3cf59d671eb72da0e7a4113c49f1f0515f462cdcf84e0f1d6045dfcbb", hex)
}

func TestCreateStorageKeyPlainV13(t *testing.T) {
	m := ExamplaryMetadataV13

	key, err := CreateStorageKey(m, "Timestamp", "Now")
	assert.NoError(t, err)
	hex, err := Hex(key)
	assert.NoError(t, err)
	assert.Equal(t, "0xf0c365c3cf59d671eb72da0e7a4113c49f1f0515f462cdcf84e0f1d6045dfcbb", hex)
}

func TestCreateStorageKeyPlainV10(t *testing.T) {
	m := ExamplaryMetadataV10

	key, err := CreateStorageKey(m, "Timestamp", "Now")
	assert.NoError(t, err)
	hex, err := Hex(key)
	assert.NoError(t, err)
	assert.Equal(t, "0xf0c365c3cf59d671eb72da0e7a4113c49f1f0515f462cdcf84e0f1d6045dfcbb", hex)
}

func TestCreateStorageKeyPlainV9(t *testing.T) {
	m := ExamplaryMetadataV9

	key, err := CreateStorageKey(m, "Timestamp", "Now")
	assert.NoError(t, err)
	hex, err := Hex(key)
	assert.NoError(t, err)
	assert.Equal(t, "0xf0c365c3cf59d671eb72da0e7a4113c49f1f0515f462cdcf84e0f1d6045dfcbb", hex)
}

func TestCreateStorageKeyPlainV4(t *testing.T) {
	m := ExamplaryMetadataV4

	key, err := CreateStorageKey(m, "Timestamp", "Now")
	assert.NoError(t, err)
	hex, err := Hex(key)
	assert.NoError(t, err)
	assert.Equal(t, "0x0e4944cfd98d6f4cc374d16f5a4e3f9c", hex)
}

func TestCreateStorageKeyMapV10(t *testing.T) {
	b := MustHexDecodeString(AlicePubKey)
	m := ExamplaryMetadataV10
	key, err := CreateStorageKey(m, "System", "AccountNonce", b)
	assert.NoError(t, err)
	hex, err := Hex(key)
	assert.NoError(t, err)
	assert.Equal(t, "0x26aa394eea5630e07c48ae0c9558cef79c2f82b23e5fd031fb54c292794b4cc42e3fb4c297a84c5cebc0e78257d213d0927ccc7596044c6ba013dd05522aacba", hex) //nolint:lll
}

func TestCreateStorageKeyMapV9(t *testing.T) {
	b := MustHexDecodeString(AlicePubKey)
	m := ExamplaryMetadataV9
	key, err := CreateStorageKey(m, "System", "AccountNonce", b)
	assert.NoError(t, err)
	hex, err := Hex(key)
	assert.NoError(t, err)
	assert.Equal(t, "0x26aa394eea5630e07c48ae0c9558cef79c2f82b23e5fd031fb54c292794b4cc42e3fb4c297a84c5cebc0e78257d213d0927ccc7596044c6ba013dd05522aacba", hex) //nolint:lll
}

func TestCreateStorageKeyMapV13(t *testing.T) {
	m := ExamplaryMetadataV13

	b := MustHexDecodeString(AlicePubKey)
	key, err := CreateStorageKey(m, "System", "Account", b)
	assert.NoError(t, err)
	hex, err := Hex(key)
	assert.NoError(t, err)
	assert.Equal(t, "0x26aa394eea5630e07c48ae0c9558cef7b99d880ec681799c0cf30e8886371da9de1e86a9a8c739864cf3cc5ec2bea59fd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d", hex)
}

func TestCreateStorageKeyMapV14(t *testing.T) {
	m := DecodedMetadataV14Example()

	b := MustHexDecodeString(AlicePubKey)
	key, err := CreateStorageKey(m, "System", "Account", b)
	assert.NoError(t, err)
	hex, err := Hex(key)
	assert.NoError(t, err)
	assert.Equal(t, "0x26aa394eea5630e07c48ae0c9558cef7b99d880ec681799c0cf30e8886371da9de1e86a9a8c739864cf3cc5ec2bea59fd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d", hex)
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
	AssertEncodedLength(t, []EncodedLengthAssert{
		{NewStorageKey(MustHexDecodeString("0x00")), 1},
		{NewStorageKey(MustHexDecodeString("0xab1234")), 3},
		{NewStorageKey(MustHexDecodeString("0x0001")), 2},
	})
}

func TestStorageKey_Encode(t *testing.T) {
	AssertEncode(t, []EncodingAssert{
		{NewStorageKey([]byte{171, 18, 52}), MustHexDecodeString("0xab1234")},
		{NewStorageKey([]byte{}), MustHexDecodeString("0x")},
	})
}

func TestStorageKey_Decode(t *testing.T) {
	bz := []byte{12, 251, 42}
	decoded := make(StorageKey, len(bz))
	err := Decode(bz, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, StorageKey(bz), decoded)
}

func TestStorageKey_Hash(t *testing.T) {
	AssertHash(t, []HashAssert{
		{NewStorageKey([]byte{0, 42, 254}), MustHexDecodeString(
			"0x537db36f5b5970b679a28a3df8d219317d658014fb9c3d409c0c799d8ecf149d")},
		{NewStorageKey([]byte{0, 0}), MustHexDecodeString(
			"0x9ee6dfb61a2fb903df487c401663825643bb825d41695e63df8af6162ab145a6")},
	})
}

func TestStorageKey_Hex(t *testing.T) {
	AssertEncodeToHex(t, []EncodeToHexAssert{
		{NewStorageKey([]byte{0, 0, 0}), "0x000000"},
		{NewStorageKey([]byte{171, 18, 52}), "0xab1234"},
		{NewStorageKey([]byte{0, 1}), "0x0001"},
		{NewStorageKey([]byte{18, 52, 86}), "0x123456"},
	})
}

func TestStorageKey_String(t *testing.T) {
	AssertString(t, []StringAssert{
		{NewStorageKey([]byte{0, 0, 0}), "[0 0 0]"},
		{NewStorageKey([]byte{171, 18, 52}), "[171 18 52]"},
		{NewStorageKey([]byte{0, 1}), "[0 1]"},
		{NewStorageKey([]byte{18, 52, 86}), "[18 52 86]"},
	})
}

func TestStorageKey_Eq(t *testing.T) {
	AssertEq(t, []EqAssert{
		{NewStorageKey([]byte{1, 0, 0}), NewStorageKey([]byte{1, 0}), false},
		{NewStorageKey([]byte{0, 0, 1}), NewStorageKey([]byte{0, 1}), false},
		{NewStorageKey([]byte{0, 0, 0}), NewStorageKey([]byte{0, 0}), false},
		{NewStorageKey([]byte{12, 48, 255}), NewStorageKey([]byte{12, 48, 255}), true},
		{NewStorageKey([]byte{0}), NewStorageKey([]byte{0}), true},
		{NewStorageKey([]byte{1}), NewBool(true), false},
		{NewStorageKey([]byte{0}), NewBool(false), false},
	})
}

func DecodedMetadataV14Example() *Metadata {
	var metadata Metadata
	err := DecodeFromHex(MetadataV14Data, &metadata)
	if err != nil {
		panic("failed to decode the example metadata v14")
	}

	return &metadata
}
