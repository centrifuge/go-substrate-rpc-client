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

var hash20 = []byte{
	1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
}

func TestH160_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewH160(hash20))
}

func TestH160_EncodedLength(t *testing.T) {
	assertEncodedLength(t, []encodedLengthAssert{{NewH160(hash20), 20}})
}

func TestH160_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{NewH160(hash20), MustHexDecodeString("0x0102030405060708090001020304050607080900")},
	})
}

func TestH160_Hash(t *testing.T) {
	assertHash(t, []hashAssert{
		{NewH160(hash20), MustHexDecodeString("0xdb34a42220dae1c95e29ee2e97d09995887790554cb1ac302680934ffc636b82")},
	})
}

func TestH160_Hex(t *testing.T) {
	assertEncodeToHex(t, []encodeToHexAssert{
		{NewH160(hash20), "0x0102030405060708090001020304050607080900"},
	})
}

func TestH160_String(t *testing.T) {
	assertString(t, []stringAssert{
		{NewH160(hash20), "[1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0]"},
	})
}

func TestH160_Eq(t *testing.T) {
	assertEq(t, []eqAssert{
		{NewH160(hash20), NewH160(hash20), true},
		{NewH160(hash20), NewBytes(hash20), false},
		{NewH160(hash20), NewBool(false), false},
	})
}

var hash32 = []byte{
	1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
	1, 2,
}

func TestH256_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewH256(hash32))
}

func TestH256_EncodedLength(t *testing.T) {
	assertEncodedLength(t, []encodedLengthAssert{{NewH256(hash32), 32}})
}

func TestH256_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{NewH256(hash32), MustHexDecodeString("0x0102030405060708090001020304050607080900010203040506070809000102")},
	})
}

func TestH256_Hash(t *testing.T) {
	assertHash(t, []hashAssert{
		{NewH256(hash32), MustHexDecodeString("0x95248da71ae3de701a70523b32a1e5a982de1e49dcf9891e188748e8f9189a2c")},
	})
}

func TestH256_Hex(t *testing.T) {
	assertEncodeToHex(t, []encodeToHexAssert{
		{NewH256(hash32), "0x0102030405060708090001020304050607080900010203040506070809000102"},
	})
}

func TestH256_String(t *testing.T) {
	assertString(t, []stringAssert{
		{NewH256(hash32), "[1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2]"},
	})
}

func TestH256_Eq(t *testing.T) {
	assertEq(t, []eqAssert{
		{NewH256(hash32), NewH256(hash32), true},
		{NewH256(hash32), NewBytes(hash32), false},
		{NewH256(hash32), NewBool(false), false},
	})
}

func TestHash_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewHash(hash32))
}

func TestHash_EncodedLength(t *testing.T) {
	assertEncodedLength(t, []encodedLengthAssert{{NewHash(hash32), 32}})
}

func TestHash_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{NewHash(hash32), MustHexDecodeString("0x0102030405060708090001020304050607080900010203040506070809000102")},
	})
}

func TestHash_Hash(t *testing.T) {
	assertHash(t, []hashAssert{
		{NewHash(hash32), MustHexDecodeString("0x95248da71ae3de701a70523b32a1e5a982de1e49dcf9891e188748e8f9189a2c")},
	})
}

func TestHash_Hex(t *testing.T) {
	assertEncodeToHex(t, []encodeToHexAssert{
		{NewHash(hash32), "0x0102030405060708090001020304050607080900010203040506070809000102"},
	})
}

func TestHash_String(t *testing.T) {
	assertString(t, []stringAssert{
		{NewHash(hash32), "[1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2]"},
	})
}

func TestHash_Eq(t *testing.T) {
	assertEq(t, []eqAssert{
		{NewHash(hash32), NewHash(hash32), true},
		{NewHash(hash32), NewBytes(hash32), false},
		{NewHash(hash32), NewBool(false), false},
	})
}

var hash64 = []byte{
	1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
	1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
	1, 2, 3, 4,
}

var hash65 = []byte{
	1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
	1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
	1, 2, 3, 4, 5,
}

func TestH512_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewH512(hash64))
}

func TestH512_EncodedLength(t *testing.T) {
	assertEncodedLength(t, []encodedLengthAssert{{NewH512(hash64), 64}})
}

func TestH512_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{NewH512(hash64), MustHexDecodeString("0x01020304050607080900010203040506070809000102030405060708090001020304050607080900010203040506070809000102030405060708090001020304")}, //nolint:lll
	})
}

func TestH512_Hash(t *testing.T) {
	assertHash(t, []hashAssert{
		{NewH512(hash64), MustHexDecodeString("0x893a41fa8d4e6447fe2d74a3ae529b1f1a13f3ac5a194907bf19e78e084a0ef6")},
	})
}

func TestH512_Hex(t *testing.T) {
	assertEncodeToHex(t, []encodeToHexAssert{
		{NewH512(hash64), "0x01020304050607080900010203040506070809000102030405060708090001020304050607080900010203040506070809000102030405060708090001020304"}, //nolint:lll
	})
}

func TestH512_String(t *testing.T) {
	assertString(t, []stringAssert{
		{NewH512(hash64), "[1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4]"}, //nolint:lll
	})
}

func TestH512_Eq(t *testing.T) {
	assertEq(t, []eqAssert{
		{NewH512(hash64), NewH512(hash64), true},
		{NewH512(hash64), NewBytes(hash64), false},
		{NewH512(hash64), NewBool(false), false},
	})
}
