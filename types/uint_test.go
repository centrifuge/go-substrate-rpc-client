// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
// Copyright (C) 2019  Philip Stanislaus, Philip Stehlik, Vimukthi Wickramasinghe
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

package types

import (
	"testing"
)

func TestU8_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewU8(0))
	assertRoundtrip(t, NewU8(12))
}

func TestU8_EncodedLength(t *testing.T) {
	assertEncodedLength(t, []encodedLengthAssert{{NewU8(13), 1}})
}

func TestU8_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{NewU8(29), mustDecodeHexString("0x1d")},
	})
}

func TestU8_Hash(t *testing.T) {
	assertHash(t, []hashAssert{
		{NewU8(29), mustDecodeHexString("0x6a9843ae0195ae1e6f95c7fbd34a42414c77e243aa18a959b5912a1f0f391b54")},
	})
}

func TestU8_Hex(t *testing.T) {
	assertEncodeToHex(t, []encodeToHexAssert{
		{NewU8(29), "0x1d"},
	})
}

func TestU8_String(t *testing.T) {
	assertString(t, []stringAssert{
		{NewU8(29), "29"},
	})
}

func TestU8_Eq(t *testing.T) {
	assertEq(t, []eqAssert{
		{NewU8(23), NewU8(23), true},
		{NewU8(23), NewBool(false), false},
	})
}

func TestU16_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewU16(0))
	assertRoundtrip(t, NewU16(12))
}

func TestU16_EncodedLength(t *testing.T) {
	assertEncodedLength(t, []encodedLengthAssert{{NewU16(13), 2}})
}

func TestU16_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{NewU16(29), mustDecodeHexString("0x1d00")},
	})
}

func TestU16_Hash(t *testing.T) {
	assertHash(t, []hashAssert{
		{NewU16(29), mustDecodeHexString("0x4e59f743a8e19ecb3022652bdef4343e62793d1f7378a688a82741b5d029d3d5")},
	})
}

func TestU16_Hex(t *testing.T) {
	assertEncodeToHex(t, []encodeToHexAssert{
		{NewU16(29), "0x1d00"},
	})
}

func TestU16_String(t *testing.T) {
	assertString(t, []stringAssert{
		{NewU16(29), "29"},
	})
}

func TestU16_Eq(t *testing.T) {
	assertEq(t, []eqAssert{
		{NewU16(23), NewU16(23), true},
		{NewU16(23), NewBool(false), false},
	})
}

func TestU32_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewU32(0))
	assertRoundtrip(t, NewU32(12))
}

func TestU32_EncodedLength(t *testing.T) {
	assertEncodedLength(t, []encodedLengthAssert{{NewU32(13), 4}})
}

func TestU32_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{NewU32(29), mustDecodeHexString("0x1d000000")},
	})
}

func TestU32_Hash(t *testing.T) {
	assertHash(t, []hashAssert{
		{NewU32(29), mustDecodeHexString("0x60ebb66f09bc7fdd21772ab1ed0efb1fd1208e3f5cd20d2d9a29a2a79b6f953f")},
	})
}

func TestU32_Hex(t *testing.T) {
	assertEncodeToHex(t, []encodeToHexAssert{
		{NewU32(29), "0x1d000000"},
	})
}

func TestU32_String(t *testing.T) {
	assertString(t, []stringAssert{
		{NewU32(29), "29"},
	})
}

func TestU32_Eq(t *testing.T) {
	assertEq(t, []eqAssert{
		{NewU32(23), NewU32(23), true},
		{NewU32(23), NewBool(false), false},
	})
}

func TestU64_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewU64(0))
	assertRoundtrip(t, NewU64(12))
}

func TestU64_EncodedLength(t *testing.T) {
	assertEncodedLength(t, []encodedLengthAssert{{NewU64(13), 8}})
}

func TestU64_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{NewU64(29), mustDecodeHexString("0x1d00000000000000")},
	})
}

func TestU64_Hash(t *testing.T) {
	assertHash(t, []hashAssert{
		{NewU64(29), mustDecodeHexString("0x83e168a13a013e6d47b0778f046aaa05d6c01d6857d044d9e9b658a6d85eb865")},
	})
}

func TestU64_Hex(t *testing.T) {
	assertEncodeToHex(t, []encodeToHexAssert{
		{NewU64(29), "0x1d00000000000000"},
	})
}

func TestU64_String(t *testing.T) {
	assertString(t, []stringAssert{
		{NewU64(29), "29"},
	})
}

func TestU64_Eq(t *testing.T) {
	assertEq(t, []eqAssert{
		{NewU64(23), NewU64(23), true},
		{NewU64(23), NewBool(false), false},
	})
}
