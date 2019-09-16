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

func TestOptionH160_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionH160(NewH160(hash20)))
	assertRoundtrip(t, NewOptionH160Empty())
}

func TestOptionH256_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionH256(NewH256(hash32)))
	assertRoundtrip(t, NewOptionH256Empty())
}

func TestOptionH512_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionH512(NewH512(hash64)))
	assertRoundtrip(t, NewOptionH512Empty())
}

func TestOptionHash_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionHash(NewHash(hash32)))
	assertRoundtrip(t, NewOptionHashEmpty())
}
