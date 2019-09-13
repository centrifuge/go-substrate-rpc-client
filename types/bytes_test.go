package types

import (
	"testing"
)

func TestBytesEncodedLength(t *testing.T) {
	testEncodedLength(t, []encodedLengthTest{
		{NewBytes(mustDecodeHexString("0x00")), 2},
		{NewBytes(mustDecodeHexString("0xab1234")), 4},
		{NewBytes(mustDecodeHexString("0x0001")), 3},
	})
}

func TestBytesEncode(t *testing.T) {
	testEncode(t, []encodingTest{
		{NewBytes([]byte{0, 0, 0}), mustDecodeHexString("0x0c000000")},
		{NewBytes([]byte{171, 18, 52}), mustDecodeHexString("0x0cab1234")},
		{NewBytes([]byte{0, 1}), mustDecodeHexString("0x080001")},
		{NewBytes([]byte{18, 52, 86}), mustDecodeHexString("0x0c123456")},
	})
}

func TestBytesHash(t *testing.T) {
	testHash(t, []hashTest{
		{NewBytes([]byte{0, 42, 254}), mustDecodeHexString(
			"0xabf7fe6eb94e0816bf2db57abb296d012f7cb9ddfe59ebf52f9c2770f49a0a46")},
		{NewBytes([]byte{0, 0}), mustDecodeHexString(
			"0xd1200120e01c48b4bbf7e1cd7ebab20087b34ea11e1e9e4ebc2f207aea77139d")},
	})
}

func TestBytesHex(t *testing.T) {
	testHex(t, []hexTest{
		{NewBytes([]byte{0, 0, 0}), "0x0c000000"},
		{NewBytes([]byte{171, 18, 52}), "0x0cab1234"},
		{NewBytes([]byte{0, 1}), "0x080001"},
		{NewBytes([]byte{18, 52, 86}), "0x0c123456"},
	})
}

func TestBytesString(t *testing.T) {
	testString(t, []stringTest{
		{NewBytes([]byte{0, 0, 0}), "0x000000"},
		{NewBytes([]byte{171, 18, 52}), "0xab1234"},
		{NewBytes([]byte{0, 1}), "0x0001"},
		{NewBytes([]byte{18, 52, 86}), "0x123456"},
	})
}

func TestBytesEq(t *testing.T) {
	testEq(t, []eqTest{
		{NewBytes([]byte{1, 0, 0}), NewBytes([]byte{1, 0}), false},
		{NewBytes([]byte{0, 0, 1}), NewBytes([]byte{0, 1}), false},
		{NewBytes([]byte{0, 0, 0}), NewBytes([]byte{0, 0}), false},
		{NewBytes([]byte{12, 48, 255}), NewBytes([]byte{12, 48, 255}), true},
		{NewBytes([]byte{0}), NewBytes([]byte{0}), true},
		{NewBytes([]byte{1}), NewBool(true), false},
		{NewBytes([]byte{0}), NewBool(false), false},
	})
}
