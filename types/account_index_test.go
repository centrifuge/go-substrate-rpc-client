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

func TestAccountIndex_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewAccountIndex(336794129))
}

func TestAccountIndex_EncodedLength(t *testing.T) {
	assertEncodedLength(t, []encodedLengthAssert{
		{NewAccountIndex(336794129), 4},
	})
}

func TestAccountIndex_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{NewAccountIndex(336794129), MustHexDecodeString("0x11121314")},
	})
}

func TestAccountIndex_Hash(t *testing.T) {
	assertHash(t, []hashAssert{
		{NewAccountIndex(336794129), MustHexDecodeString(
			"0xa6730c0d3a95e0ff2068fa9a6ecf82c42c494c8c2cdd65379c898a4b88dd7138")},
	})
}

func TestAccountIndex_Hex(t *testing.T) {
	assertEncodeToHex(t, []encodeToHexAssert{
		{NewAccountIndex(336794129), "0x11121314"},
	})
}

func TestAccountIndex_String(t *testing.T) {
	assertString(t, []stringAssert{
		{NewAccountIndex(336794129), "336794129"},
	})
}

func TestAccountIndex_Eq(t *testing.T) {
	assertEq(t, []eqAssert{
		{NewAccountIndex(336794129), NewAccountIndex(336794129), true},
		{NewAccountIndex(336794129), NewAccountIndex(12), false},
		{NewAccountIndex(336794129), NewBool(false), false},
	})
}
