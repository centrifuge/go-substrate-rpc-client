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

func TestAccountID_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewAccountID([]byte{}))
	assertRoundtrip(t, NewAccountID([]byte{0, 1, 2, 3, 4, 5, 6, 7}))
}

func TestAccountID_EncodedLength(t *testing.T) {
	assertEncodedLength(t, []encodedLengthAssert{
		{NewAccountID([]byte{}), 32},
		{NewAccountID([]byte{7, 6, 5, 4, 3, 2, 1, 0}), 32},
	})
}

func TestAccountID_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{NewAccountID([]byte{0, 0, 0}), MustHexDecodeString("0x0000000000000000000000000000000000000000000000000000000000000000")},     //nolint:lll
		{NewAccountID([]byte{171, 18, 52}), MustHexDecodeString("0xab12340000000000000000000000000000000000000000000000000000000000")}, //nolint:lll
	})
}

func TestAccountID_Hash(t *testing.T) {
	assertHash(t, []hashAssert{
		{NewAccountID([]byte{0, 42, 254}), MustHexDecodeString(
			"0x7834db8eb04aefe8272c32d8160ce4fa3cb31fc95882e5bd53860715731c8198")},
		{NewAccountID([]byte{0, 0}), MustHexDecodeString(
			"0x89eb0d6a8a691dae2cd15ed0369931ce0a949ecafa5c3f93f8121833646e15c3")},
	})
}

func TestAccountID_Hex(t *testing.T) {
	assertEncodeToHex(t, []encodeToHexAssert{
		{NewAccountID([]byte{0, 0, 0}), "0x0000000000000000000000000000000000000000000000000000000000000000"},
		{NewAccountID([]byte{171, 18, 52}), "0xab12340000000000000000000000000000000000000000000000000000000000"},
	})
}

func TestAccountID_String(t *testing.T) {
	assertString(t, []stringAssert{
		{NewAccountID([]byte{0, 0, 0}), "[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]"},
		{NewAccountID([]byte{171, 18, 52}), "[171 18 52 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]"},
	})
}

func TestAccountID_Eq(t *testing.T) {
	assertEq(t, []eqAssert{
		{NewAccountID([]byte{1, 0, 0}), NewAccountID([]byte{1, 0}), true},
		{NewAccountID([]byte{0, 0, 1}), NewAccountID([]byte{0, 1}), false},
		{NewAccountID([]byte{0, 0, 0}), NewAccountID([]byte{0, 0}), true},
		{NewAccountID([]byte{12, 48, 255}), NewAccountID([]byte{12, 48, 255}), true},
		{NewAccountID([]byte{0}), NewAccountID([]byte{0}), true},
		{NewAccountID([]byte{1}), NewBool(true), false},
		{NewAccountID([]byte{0}), NewBool(false), false},
	})
}
