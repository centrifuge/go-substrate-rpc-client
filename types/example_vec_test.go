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
	"fmt"
	"reflect"
)

func ExampleExampleVec_simple() {
	ingredients := []string{"salt", "sugar"}

	encoded, err := EncodeToHexString(ingredients)
	if err != nil {
		panic(err)
	}
	fmt.Println(encoded)

	var decoded []string
	err = DecodeFromHexString(encoded, &decoded)
	if err != nil {
		panic(err)
	}
	fmt.Println(decoded)
	// Output: 0x081073616c74147375676172
	// [salt sugar]
}

func ExampleExampleVec_struct() {
	type Votes struct {
		Options     [2]string
		Yay         []string
		Nay         []string
		Outstanding []string
	}

	votes := Votes{
		Options:     [2]string{"no deal", "muddle through"},
		Yay:         []string{"Alice"},
		Nay:         nil,
		Outstanding: []string{"Bob", "Carol"},
	}

	encoded, err := EncodeToBytes(votes)
	if err != nil {
		panic(err)
	}
	var decoded Votes
	err = DecodeFromBytes(encoded, &decoded)
	if err != nil {
		panic(err)
	}

	fmt.Println(reflect.DeepEqual(votes, decoded))
	// Output: true
}

// type MyOption struct {
// 	Woohoo *bool
// }

// type MyOptionNoPoointer struct {
// 	Woohoo bool
// }

// func NewMyOption(b bool) MyOption {
// 	return MyOption{&b}
// }

// func Example2() {
// 	myopt := NewMyOption(true)
// 	// myopt := NewBool(true)

// 	encoded, err := EncodeToHexString(myopt)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(encoded)

// 	var decoded MyOption
// 	err = DecodeFromHexString(encoded, &decoded)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(decoded)
// 	// Output: 0x081073616c74147375676172
// 	// [salt sugar]
// }
