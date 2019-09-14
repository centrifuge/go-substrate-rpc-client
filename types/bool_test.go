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

func TestBool_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewBool(true))
	assertRoundtrip(t, NewBool(false))
}

func TestBool_EncodedLength(t *testing.T) {
	assertEncodedLength(t, []encodedLengthAssert{
		{NewBool(true), 1},
		{NewBool(false), 1},
	})
}

func TestBool_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{NewBool(true), []byte{0x01}},
		{NewBool(false), []byte{0x00}},
	})
}

func TestBool_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{[]byte{0x01}, NewBool(true)},
		{[]byte{0x00}, NewBool(false)},
	})
}

func TestBool_Hash(t *testing.T) {
	assertHash(t, []hashAssert{
		{NewBool(true), mustDecodeHexString("0xee155ace9c40292074cb6aff8c9ccdd273c81648ff1149ef36bcea6ebb8a3e25")},
		{NewBool(false), mustDecodeHexString("0x03170a2e7597b7b7e3d84c05391d139a62b157e78786d8c082f29dcf4c111314")},
	})
}

func TestBool_Hex(t *testing.T) {
	assertHex(t, []hexAssert{
		{NewBool(true), "0x01"},
		{NewBool(false), "0x00"},
	})
}

func TestBool_String(t *testing.T) {
	assertString(t, []stringAssert{
		{NewBool(true), "true"},
		{NewBool(false), "false"},
	})
}

func TestBool_Eq(t *testing.T) {
	assertEq(t, []eqAssert{
		{NewBool(true), NewBool(true), true},
		{NewBool(false), NewBool(true), false},
		{NewBool(false), NewBool(false), true},
		{NewBool(true), NewBytes([]byte{0, 1, 2}), false},
		{NewBool(true), NewBytes([]byte{1}), false},
		{NewBool(false), NewBytes([]byte{0}), false},
	})
}
