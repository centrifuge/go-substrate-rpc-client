package types

import (
	"testing"
)

func TestEncodedLength(t *testing.T) {
	testEncodedLength(t, []encodedLengthTest{
		{NewBool(true), 1},
		{NewBool(false), 1},
	})
}

func TestEncode(t *testing.T) {
	testEncode(t, []encodingTest{
		{NewBool(true), []byte{0x01}},
		{NewBool(false), []byte{0x00}},
	})
}

func TestHash(t *testing.T) {
	testHash(t, []hashTest{
		{NewBool(true), mustDecodeHexString("0xee155ace9c40292074cb6aff8c9ccdd273c81648ff1149ef36bcea6ebb8a3e25")},
		{NewBool(false), mustDecodeHexString("0x03170a2e7597b7b7e3d84c05391d139a62b157e78786d8c082f29dcf4c111314")},
	})
}

func TestHex(t *testing.T) {
	testHex(t, []hexTest{
		{NewBool(true), "0x01"},
		{NewBool(false), "0x00"},
	})
}

func TestString(t *testing.T) {
	testString(t, []stringTest{
		{NewBool(true), "true"},
		{NewBool(false), "false"},
	})
}

func TestEq(t *testing.T) {
	testEq(t, []eqTest{
		{NewBool(true), NewBool(true), true},
		{NewBool(false), NewBool(true), false},
		{NewBool(false), NewBool(false), true},
		{NewBool(true), NewBytes([]byte{0, 1, 2}), false},
		{NewBool(true), NewBytes([]byte{1}), false},
		{NewBool(false), NewBytes([]byte{0}), false},
	})
}
