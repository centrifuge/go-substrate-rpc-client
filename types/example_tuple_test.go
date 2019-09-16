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
	"fmt"

	"golang.org/x/crypto/blake2b"

	. "github.com/centrifuge/go-substrate-rpc-client/types"
)

func ExampleExampleTuple() {
	// This represents a document tuple of types [uint64, hash]
	type Doc struct {
		ID   U64
		Hash Hash
	}

	doc := Doc{12, blake2b.Sum256([]byte("My document"))}

	encoded, err := EncodeToHexString(doc)
	if err != nil {
		panic(err)
	}
	fmt.Println(encoded)

	var decoded Doc
	err = DecodeFromHexString(encoded, &decoded)
	if err != nil {
		panic(err)
	}
	fmt.Println(decoded)
	// Output: 0x0c00000000000000809199a254aedc9d92a3157cd27bd21ceccc1e2ecee5760788663a3e523bc1a759
	// {12 [145 153 162 84 174 220 157 146 163 21 124 210 123 210 28 236 204 30 46 206 229 118 7 136 102 58 62 82 59 193 167 89]}
}
