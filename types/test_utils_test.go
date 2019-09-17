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
	"bytes"
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/scale"
	. "github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/stretchr/testify/assert"
)

type encodedLengthAssert struct {
	input    interface{}
	expected int
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if reflect.DeepEqual(a, b) {
		return
	}
	t.Errorf("Received %v (type %v), expected %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}

func assertRoundtrip(t *testing.T, value interface{}) {
	var buffer = bytes.Buffer{}
	err := scale.NewEncoder(&buffer).Encode(value)
	assert.NoError(t, err)
	target := reflect.New(reflect.TypeOf(value))
	err = scale.NewDecoder(&buffer).Decode(target.Interface())
	assert.NoError(t, err)
	assertEqual(t, target.Elem().Interface(), value)
}

func assertEncodedLength(t *testing.T, encodedLengthAsserts []encodedLengthAssert) {
	for _, test := range encodedLengthAsserts {
		result, err := EncodedLength(test.input)
		if err != nil {
			t.Errorf("Encoded length error for input %v: %v\n", test.input, err)
		}
		if result != test.expected {
			t.Errorf("Fail, input %v, expected %v, result %v\n", test.input, test.expected, result)
		}
	}
}

type encodingAssert struct {
	input    interface{}
	expected []byte
}

func assertEncode(t *testing.T, encodingAsserts []encodingAssert) {
	for _, test := range encodingAsserts {
		result, err := EncodeToBytes(test.input)
		if err != nil {
			t.Errorf("Encoding error for input %v: %v\n", test.input, err)
		}
		if !bytes.Equal(result, test.expected) {
			t.Errorf("Fail, input %v, expected %#x, result %#x\n", test.input, test.expected, result)
		}
	}
}

type decodingAssert struct {
	input    []byte
	expected interface{}
}

func assertDecode(t *testing.T, decodingAsserts []decodingAssert) {
	for _, test := range decodingAsserts {
		target := reflect.New(reflect.TypeOf(test.expected))
		err := DecodeFromBytes(test.input, target.Interface())
		if err != nil {
			t.Errorf("Encoding error for input %v: %v\n", test.input, err)
		}
		assertEqual(t, target.Elem().Interface(), test.expected)
	}
}

type hashAssert struct {
	input    interface{}
	expected []byte
}

func assertHash(t *testing.T, hashAsserts []hashAssert) {
	for _, test := range hashAsserts {
		result, err := GetHash(test.input)
		if err != nil {
			t.Errorf("Hash error for input %v: %v\n", test.input, err)
		}
		if !bytes.Equal(result[:], test.expected) {
			t.Errorf("Fail, input %v, expected %#x, result %#x\n", test.input, test.expected, result)
		}
	}
}

type encodeToHexAssert struct {
	input    interface{}
	expected string
}

func assertEncodeToHex(t *testing.T, encodeToHexAsserts []encodeToHexAssert) {
	for _, test := range encodeToHexAsserts {
		result, err := EncodeToHexString(test.input)
		if err != nil {
			t.Errorf("Hex error for input %v: %v\n", test.input, err)
		}
		if result != test.expected {
			t.Errorf("Fail, input %v, expected %v, result %v\n", test.input, test.expected, result)
		}
	}
}

type stringAssert struct {
	input    interface{}
	expected string
}

func assertString(t *testing.T, stringAsserts []stringAssert) {
	for _, test := range stringAsserts {
		result := fmt.Sprintf("%v", test.input)
		if result != test.expected {
			t.Errorf("Fail, input %v, expected %v, result %v\n", test.input, test.expected, result)
		}
	}
}

type eqAssert struct {
	input    interface{}
	other    interface{}
	expected bool
}

func assertEq(t *testing.T, eqAsserts []eqAssert) {
	for _, test := range eqAsserts {
		result := Eq(test.input, test.other)
		if result != test.expected {
			t.Errorf("Fail, input %v, other %v, expected %v, result %v\n", test.input, test.other, test.expected, result)
		}
	}
}

// mustDecodeHexString panics if str cannot be decoded. Param str is expected to start with "0x"
func mustDecodeHexString(str string) []byte {
	bz, err := hex.DecodeString(str[2:])
	if err != nil {
		panic(err)
	}
	return bz
}
