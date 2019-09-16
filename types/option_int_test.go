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

package types

import (
	"testing"
)

func TestOptionI8_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionI8(NewI8(7)))
	assertRoundtrip(t, NewOptionI8(NewI8(0)))
	assertRoundtrip(t, NewOptionI8Empty())
}

func TestOptionI16_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionI16(NewI16(14)))
	assertRoundtrip(t, NewOptionI16(NewI16(0)))
	assertRoundtrip(t, NewOptionI16Empty())
}

func TestOptionI32_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionI32(NewI32(21)))
	assertRoundtrip(t, NewOptionI32(NewI32(0)))
	assertRoundtrip(t, NewOptionI32Empty())
}

func TestOptionI64_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionI64(NewI64(28)))
	assertRoundtrip(t, NewOptionI64(NewI64(0)))
	assertRoundtrip(t, NewOptionI64Empty())
}
