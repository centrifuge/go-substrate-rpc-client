package types

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/centrifuge/go-substrate-rpc-client/scale"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
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

type hexAssert struct {
	input    interface{}
	expected string
}

func assertHex(t *testing.T, hexAsserts []hexAssert) {
	for _, test := range hexAsserts {
		result, err := Hex(test.input)
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
