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
)

func TestSignature_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewSignature(hash64))
}

func TestSignature_EncodedLength(t *testing.T) {
	assertEncodedLength(t, []encodedLengthAssert{{NewSignature(hash64), 66}})
}

func TestSignature_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{NewSignature(hash64), mustDecodeHexString("0x010101020304050607080900010203040506070809000102030405060708090001020304050607080900010203040506070809000102030405060708090001020304")}, //nolint:lll
	})
}

func TestSignature_Hash(t *testing.T) {
	assertHash(t, []hashAssert{
		{NewSignature(hash64), mustDecodeHexString("0x0926d23398a248b1c7723651a5ad05a5626cc8f9450512d6c3b5b2156615bcd5")},
	})
}

func TestSignature_Hex(t *testing.T) {
	assertEncodeToHex(t, []encodeToHexAssert{
		{NewSignature(hash64), "0x010101020304050607080900010203040506070809000102030405060708090001020304050607080900010203040506070809000102030405060708090001020304"}, //nolint:lll
	})
}

func TestSignature_String(t *testing.T) {
	assertString(t, []stringAssert{
		{NewSignature(hash64), "[1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4]"}, //nolint:lll
	})
}

func TestSignature_Eq(t *testing.T) {
	assertEq(t, []eqAssert{
		{NewSignature(hash64), NewSignature(hash64), true},
		{NewSignature(hash64), NewBytes(hash64[:]), false},
		{NewSignature(hash64), NewBool(false), false},
	})
}
