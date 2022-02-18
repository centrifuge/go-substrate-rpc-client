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

	. "github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

var testMultiSig1 = MultiSignature{IsEd25519: true, AsEd25519: NewSignature(hash64)}
var testMultiSig2 = MultiSignature{IsSr25519: true, AsSr25519: NewSignature(hash64)}
var testMultiSig3 = MultiSignature{IsEcdsa: true, AsEcdsa: NewEcdsaSignature(hash65)}

func TestMultiSignature_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, testMultiSig1)
	assertRoundtrip(t, testMultiSig2)
	assertRoundtrip(t, testMultiSig3)
}

func TestMultiSignature_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{testMultiSig1, MustHexDecodeString("0x0001020304050607080900010203040506070809000102030405060708090001020304050607080900010203040506070809000102030405060708090001020304")},   //nolint:lll
		{testMultiSig2, MustHexDecodeString("0x0101020304050607080900010203040506070809000102030405060708090001020304050607080900010203040506070809000102030405060708090001020304")},   //nolint:lll
		{testMultiSig3, MustHexDecodeString("0x020102030405060708090001020304050607080900010203040506070809000102030405060708090001020304050607080900010203040506070809000102030405")}, //nolint:lll
	})
}

func TestMultiSignature_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{MustHexDecodeString("0x0001020304050607080900010203040506070809000102030405060708090001020304050607080900010203040506070809000102030405060708090001020304"), testMultiSig1},   //nolint:lll
		{MustHexDecodeString("0x0101020304050607080900010203040506070809000102030405060708090001020304050607080900010203040506070809000102030405060708090001020304"), testMultiSig2},   //nolint:lll
		{MustHexDecodeString("0x020102030405060708090001020304050607080900010203040506070809000102030405060708090001020304050607080900010203040506070809000102030405"), testMultiSig3}, //nolint:lll
	})
}
