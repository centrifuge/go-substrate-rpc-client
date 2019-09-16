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
	"fmt"
)

func ExampleExampleStruct() {
	type Animal struct {
		Name     string
		Legs     U8
		Children []string
	}

	dog := Animal{Name: "Bello", Legs: 2, Children: []string{"Sam"}}

	encoded, err := EncodeToHexString(dog)
	if err != nil {
		panic(err)
	}
	fmt.Println(encoded)

	var decoded Animal
	err = DecodeFromHexString(encoded, &decoded)
	if err != nil {
		panic(err)
	}
	fmt.Println(decoded)
	// Output: 0x1442656c6c6f02040c53616d
	// {Bello 2 [Sam]}
}
