package codec

import (
	"bytes"
	"encoding/hex"
	"testing"
)

type encodedLengthTest struct {
	input    Codec
	expected int
}

func testEncodedLength(t *testing.T, encodedLengthTests []encodedLengthTest) {
	for _, test := range encodedLengthTests {
		result, err := test.input.EncodedLength()
		if err != nil {
			t.Errorf("Encoded length error for input %v: %v\n", test.input, err)
		}
		if result != test.expected {
			t.Errorf("Fail, input %v, expected %x, result %x\n", test.input, test.expected, result)
		}
	}
}

type encodingTest struct {
	input    Codec
	expected []byte
}

func testEncode(t *testing.T, encodingTests []encodingTest) {
	for _, test := range encodingTests {
		result, err := test.input.Encode()
		if err != nil {
			t.Errorf("Encoding error for input %v: %v\n", test.input, err)
		}
		if !bytes.Equal(result, test.expected) {
			t.Errorf("Fail, input %v, expected %x, result %x\n", test.input, test.expected, result)
		}
	}
}

type hashTest struct {
	input    Codec
	expected []byte
}

func testHash(t *testing.T, hashTests []hashTest) {
	for _, test := range hashTests {
		result, err := test.input.Hash()
		if err != nil {
			t.Errorf("Hash error for input %v: %v\n", test.input, err)
		}
		if !bytes.Equal(result[:], test.expected) {
			t.Errorf("Fail, input %v, expected %x, result %x\n", test.input, test.expected, result)
		}
	}
}

type hexTest struct {
	input    Codec
	expected string
}

func testHex(t *testing.T, encodingTests []hexTest) {
	for _, test := range encodingTests {
		result, err := test.input.Hex()
		if err != nil {
			t.Errorf("Hex error for input %v: %v\n", test.input, err)
		}
		if result != test.expected {
			t.Errorf("Fail, input %v, expected %v, result %v\n", test.input, test.expected, result)
		}
	}
}

type stringTest struct {
	input    Codec
	expected string
}

func testString(t *testing.T, stringTests []stringTest) {
	for _, test := range stringTests {
		result := test.input.String()
		if result != test.expected {
			t.Errorf("Fail, input %v, expected %v, result %v\n", test.input, test.expected, result)
		}
	}
}

type eqTest struct {
	input    Codec
	other    Codec
	expected bool
}

func testEq(t *testing.T, eqTests []eqTest) {
	for _, test := range eqTests {
		result := test.input.Eq(test.other)
		if result != test.expected {
			t.Errorf("Fail, input %v, other %v, expected %v, result %v\n", test.input, test.other, test.expected, result)
		}
	}
}

func mustDecodeHexString(str string) []byte {
	bytes, err := hex.DecodeString(str)
	if err != nil {
		panic(err)
	}
	return bytes
}
