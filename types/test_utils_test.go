// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
//
// Copyright 2019 Centrifuge GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types_test

import (
	"bytes"
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
	t.Errorf("Received %#v (type %v), expected %#v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
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
