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

	"github.com/centrifuge/go-substrate-rpc-client/signature"
	. "github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/stretchr/testify/assert"
)

var examplaryExtrinsicPayload = ExtrinsicPayloadV4{ExtrinsicPayloadV3: ExtrinsicPayloadV3{Method: BytesBare{0x6, 0x0, 0xff, 0xd7, 0x56, 0x8e, 0x5f, 0xa, 0x7e, 0xda, 0x67, 0xa8, 0x26, 0x91, 0xff, 0x37, 0x9a, 0xc4, 0xbb, 0xa4, 0xf9, 0xc9, 0xb8, 0x59, 0xfe, 0x77, 0x9b, 0x5d, 0x46, 0x36, 0x3b, 0x61, 0xad, 0x2d, 0xb9, 0xe5, 0x6c}, Era: ExtrinsicEra{IsImmortalEra: false, IsMortalEra: true, AsMortalEra: MortalEra{First: 0x7, Second: 0x3}}, Nonce: NewUCompactFromUInt(0x1234), Tip: NewUCompactFromUInt(0x5678), SpecVersion: 0x7b, GenesisHash: Hash{0xdc, 0xd1, 0x34, 0x67, 0x1, 0xca, 0x83, 0x96, 0x49, 0x6e, 0x52, 0xaa, 0x27, 0x85, 0xb1, 0x74, 0x8d, 0xeb, 0x6d, 0xb0, 0x95, 0x51, 0xb7, 0x21, 0x59, 0xdc, 0xb3, 0xe0, 0x89, 0x91, 0x2, 0x5b}, BlockHash: Hash{0xde, 0x8f, 0x69, 0xee, 0xb5, 0xe0, 0x65, 0xe1, 0x8c, 0x69, 0x50, 0xff, 0x70, 0x8d, 0x7e, 0x55, 0x1f, 0x68, 0xdc, 0x9b, 0xf5, 0x9a, 0x7, 0xc5, 0x23, 0x67, 0xc0, 0x28, 0xf, 0x80, 0x5e, 0xc7}}, TransactionVersion: 1}  //nolint:lll

func TestExtrinsicPayload(t *testing.T) {
	var era ExtrinsicEra
	err := DecodeFromHexString("0x0703", &era)
	assert.NoError(t, err)

	p := ExtrinsicPayloadV4{
		ExtrinsicPayloadV3: ExtrinsicPayloadV3{
			Method: MustHexDecodeString(
				"0x0600ffd7568e5f0a7eda67a82691ff379ac4bba4f9c9b859fe779b5d46363b61ad2db9e56c"),
			Era:         era,
			Nonce:       NewUCompactFromUInt(4660),
			Tip:         NewUCompactFromUInt(22136),
			SpecVersion: 123,
			GenesisHash: NewHash(MustHexDecodeString("0xdcd1346701ca8396496e52aa2785b1748deb6db09551b72159dcb3e08991025b")),
			BlockHash:   NewHash(MustHexDecodeString("0xde8f69eeb5e065e18c6950ff708d7e551f68dc9bf59a07c52367c0280f805ec7")),
		},
		TransactionVersion: 1,
	}

	assert.Equal(t, examplaryExtrinsicPayload, p)

	enc, err := EncodeToHexString(examplaryExtrinsicPayload)
	assert.NoError(t, err)

	assert.Equal(t, "0x"+
		"0600ffd7568e5f0a7eda67a82691ff379ac4bba4f9c9b859fe779b5d46363b61ad2db9e56c"+ // Method
		"0703"+ // Era
		"d148"+ // Nonce
		"e2590100"+ // Tip
		"7b000000"+ // Spec version
		"01000000"+ // Tx version
		"dcd1346701ca8396496e52aa2785b1748deb6db09551b72159dcb3e08991025b"+ // Genesis Hash
		"de8f69eeb5e065e18c6950ff708d7e551f68dc9bf59a07c52367c0280f805ec7", // BlockHash
		enc)

	// b := bytes.NewBuffer(MustHexDecodeString())

	var dec ExtrinsicPayloadV4
	err = DecodeFromHexString(enc, &dec)
	assert.Error(t, err)
}

func TestExtrinsicPayload_Sign(t *testing.T) {
	sig, err := examplaryExtrinsicPayload.Sign(signature.TestKeyringPairAlice)
	assert.NoError(t, err)

	// verify sig
	b, err := EncodeToBytes(examplaryExtrinsicPayload)
	assert.NoError(t, err)
	ok, err := signature.Verify(b, sig[:], HexEncodeToString(signature.TestKeyringPairAlice.PublicKey))
	assert.NoError(t, err)
	assert.True(t, ok)
}
