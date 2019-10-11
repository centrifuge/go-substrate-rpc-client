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
)

func TestAccountInfo_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewAccountInfo([]byte{1, 2, 3}, 13))
}

func TestAccountInfo_EncodedLength(t *testing.T) {
	assertEncodedLength(t, []encodedLengthAssert{
		{NewAccountInfo([]byte{1, 2, 3}, 13), 12},
	})
}

func TestAccountInfo_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{NewAccountInfo([]byte{1, 2, 3}, 13), MustHexDecodeString("0x0c0102030d00000000000000")},
	})
}

func TestAccountInfo_Hash(t *testing.T) {
	assertHash(t, []hashAssert{
		{NewAccountInfo([]byte{1, 2, 3}, 13), MustHexDecodeString(
			"0x4fac0dfeb9b4efd2518c762e7d097fafaffaf8d56a2e784f9fc9919c22277804")},
	})
}

func TestAccountInfo_Hex(t *testing.T) {
	assertEncodeToHex(t, []encodeToHexAssert{
		{NewAccountInfo([]byte{1, 2, 3}, 13), "0x0c0102030d00000000000000"},
	})
}

func TestAccountInfo_String(t *testing.T) {
	assertString(t, []stringAssert{
		{NewAccountInfo([]byte{1, 2, 3}, 13), "{[1 2 3] 13}"},
	})
}

func TestAccountInfo_Eq(t *testing.T) {
	assertEq(t, []eqAssert{
		{NewAccountInfo([]byte{1, 2, 3}, 13), NewAccountInfo([]byte{1, 2, 3}, 13), true},
		{NewAccountInfo([]byte{1, 2, 3}, 13), NewAccountInfo([]byte{1, 2, 2}, 13), false},
		{NewAccountInfo([]byte{1, 2, 3}, 13), NewBool(false), false},
	})
}
