package types

import (
	"testing"
)

func TestBytes_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewBytes(mustDecodeHexString("0x00")))
	assertRoundtrip(t, NewBytes(mustDecodeHexString("0xab1234")))
	assertRoundtrip(t, NewBytes(mustDecodeHexString("0x0001")))
}

func TestBytes_EncodedLength(t *testing.T) {
	assertEncodedLength(t, []encodedLengthAssert{
		{NewBytes(mustDecodeHexString("0x00")), 2},
		{NewBytes(mustDecodeHexString("0xab1234")), 4},
		{NewBytes(mustDecodeHexString("0x0001")), 3},
	})
}

func TestBytes_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{NewBytes([]byte{0, 0, 0}), mustDecodeHexString("0x0c000000")},
		{NewBytes([]byte{171, 18, 52}), mustDecodeHexString("0x0cab1234")},
		{NewBytes([]byte{0, 1}), mustDecodeHexString("0x080001")},
		{NewBytes([]byte{18, 52, 86}), mustDecodeHexString("0x0c123456")},
	})
}

func TestBytes_Hash(t *testing.T) {
	assertHash(t, []hashAssert{
		{NewBytes([]byte{0, 42, 254}), mustDecodeHexString(
			"0xabf7fe6eb94e0816bf2db57abb296d012f7cb9ddfe59ebf52f9c2770f49a0a46")},
		{NewBytes([]byte{0, 0}), mustDecodeHexString(
			"0xd1200120e01c48b4bbf7e1cd7ebab20087b34ea11e1e9e4ebc2f207aea77139d")},
	})
}

func TestBytes_Hex(t *testing.T) {
	assertHex(t, []hexAssert{
		{NewBytes([]byte{0, 0, 0}), "0x0c000000"},
		{NewBytes([]byte{171, 18, 52}), "0x0cab1234"},
		{NewBytes([]byte{0, 1}), "0x080001"},
		{NewBytes([]byte{18, 52, 86}), "0x0c123456"},
	})
}

func TestBytes_String(t *testing.T) {
	assertString(t, []stringAssert{
		{NewBytes([]byte{0, 0, 0}), "[0 0 0]"},
		{NewBytes([]byte{171, 18, 52}), "[171 18 52]"},
		{NewBytes([]byte{0, 1}), "[0 1]"},
		{NewBytes([]byte{18, 52, 86}), "[18 52 86]"},
	})
}

func TestBytes_Eq(t *testing.T) {
	assertEq(t, []eqAssert{
		{NewBytes([]byte{1, 0, 0}), NewBytes([]byte{1, 0}), false},
		{NewBytes([]byte{0, 0, 1}), NewBytes([]byte{0, 1}), false},
		{NewBytes([]byte{0, 0, 0}), NewBytes([]byte{0, 0}), false},
		{NewBytes([]byte{12, 48, 255}), NewBytes([]byte{12, 48, 255}), true},
		{NewBytes([]byte{0}), NewBytes([]byte{0}), true},
		{NewBytes([]byte{1}), NewBool(true), false},
		{NewBytes([]byte{0}), NewBool(false), false},
	})
}
