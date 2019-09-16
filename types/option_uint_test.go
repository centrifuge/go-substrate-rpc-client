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

func TestOptionU8_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionU8(NewU8(7)))
	assertRoundtrip(t, NewOptionU8(NewU8(0)))
	assertRoundtrip(t, NewOptionU8Empty())
}

func TestOptionU16_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionU16(NewU16(14)))
	assertRoundtrip(t, NewOptionU16(NewU16(0)))
	assertRoundtrip(t, NewOptionU16Empty())
}

func TestOptionU32_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionU32(NewU32(21)))
	assertRoundtrip(t, NewOptionU32(NewU32(0)))
	assertRoundtrip(t, NewOptionU32Empty())
}

func TestOptionU64_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionU64(NewU64(28)))
	assertRoundtrip(t, NewOptionU64(NewU64(0)))
	assertRoundtrip(t, NewOptionU64Empty())
}
