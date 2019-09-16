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

	"github.com/centrifuge/go-substrate-rpc-client/scale"
)

// MyVal is a custom type that is used to hold arbitrarily encoded data. In this example, we encode uint8s with a 0x00
// and strings with 0x01 as the first byte.
type MyVal struct {
	Value interface{}
}

func (a *MyVal) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()

	if err != nil {
		return err
	}

	if b == 0 {
		var u uint8
		err = decoder.Decode(&u)
		a.Value = u
	} else if b == 1 {
		var s string
		err = decoder.Decode(&s)
		a.Value = s
	}

	if err != nil {
		return err
	}

	return nil
}

func (a MyVal) Encode(encoder scale.Encoder) error {
	var err1, err2 error

	switch v := a.Value.(type) {
	case uint8:
		err1 = encoder.PushByte(0)
		err2 = encoder.Encode(v)
	case string:
		err1 = encoder.PushByte(1)
		err2 = encoder.Encode(v)
	default:
		return fmt.Errorf("unknown type %T", v)
	}

	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}

	return nil
}

func ExampleExampleVecAny() {
	myValSlice := []MyVal{{uint8(12)}, {"Abc"}}

	encoded, err := EncodeToBytes(myValSlice)
	if err != nil {
		panic(err)
	}
	fmt.Println(encoded)

	var decoded []MyVal
	err = DecodeFromBytes(encoded, &decoded)
	if err != nil {
		panic(err)
	}

	fmt.Println(reflect.DeepEqual(myValSlice, decoded))
	// Output: [8 0 12 1 12 65 98 99]
	// true
}
