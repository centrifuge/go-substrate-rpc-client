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

func TestAccountId_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewAccountId([32]byte{}))
	assertRoundtrip(t, NewAccountId([32]byte{0, 1, 2, 3, 4, 5, 6, 7}))
}

func TestAccountId_EncodedLength(t *testing.T) {
	assertEncodedLength(t, []encodedLengthAssert{
		{NewAccountId([32]byte{}), 33},
		{NewAccountId([32]byte{7, 6, 5, 4, 3, 2, 1, 0}), 33},
	})
}

func TestAccountId_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{NewAccountId([32]byte{0, 0, 0}), mustDecodeHexString("0x800000000000000000000000000000000000000000000000000000000000000000")},
		{NewAccountId([32]byte{171, 18, 52}), mustDecodeHexString("0x80ab12340000000000000000000000000000000000000000000000000000000000")},
	})
}

func TestAccountId_Hash(t *testing.T) {
	assertHash(t, []hashAssert{
		{NewAccountId([32]byte{0, 42, 254}), mustDecodeHexString(
			"0xa0824cc9ecc0a05ed2ed8974e3564e02f6a544d1c49eb1375f14e9830854eeed")},
		{NewAccountId([32]byte{0, 0}), mustDecodeHexString(
			"0xaf7bedde1fea222230b82d63d5b665ac75afbe4ad3f75999bb3386cf994a6963")},
	})
}

func TestAccountId_Hex(t *testing.T) {
	assertEncodeToHex(t, []encodeToHexAssert{
		{NewAccountId([32]byte{0, 0, 0}), "0x800000000000000000000000000000000000000000000000000000000000000000"},
		{NewAccountId([32]byte{171, 18, 52}), "0x80ab12340000000000000000000000000000000000000000000000000000000000"},
	})
}

func TestAccountId_String(t *testing.T) {
	assertString(t, []stringAssert{
		{NewAccountId([32]byte{0, 0, 0}), "[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]"},
		{NewAccountId([32]byte{171, 18, 52}), "[171 18 52 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]"},
	})
}

func TestAccountId_Eq(t *testing.T) {
	assertEq(t, []eqAssert{
		{NewAccountId([32]byte{1, 0, 0}), NewAccountId([32]byte{1, 0}), true},
		{NewAccountId([32]byte{0, 0, 1}), NewAccountId([32]byte{0, 1}), false},
		{NewAccountId([32]byte{0, 0, 0}), NewAccountId([32]byte{0, 0}), true},
		{NewAccountId([32]byte{12, 48, 255}), NewAccountId([32]byte{12, 48, 255}), true},
		{NewAccountId([32]byte{0}), NewAccountId([32]byte{0}), true},
		{NewAccountId([32]byte{1}), NewBool(true), false},
		{NewAccountId([32]byte{0}), NewBool(false), false},
	})
}
