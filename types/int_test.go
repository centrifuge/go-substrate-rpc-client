package types

import (
	"testing"
)

func TestI8_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewI8(0))
	assertRoundtrip(t, NewI8(12))
	assertRoundtrip(t, NewI8(-12))
}

func TestI8_EncodedLength(t *testing.T) {
	assertEncodedLength(t, []encodedLengthAssert{{NewI8(-13), 1}})
}

func TestI8_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{NewI8(-29), mustDecodeHexString("0xe3")},
	})
}

func TestI8_Hash(t *testing.T) {
	assertHash(t, []hashAssert{
		{NewI8(-29), mustDecodeHexString("0xb683f1b6c99388ff3443b35a0051eeaafdc5e364e771bdfc72c7fd5d2be800bc")},
	})
}

func TestI8_Hex(t *testing.T) {
	assertHex(t, []hexAssert{
		{NewI8(-29), "0xe3"},
	})
}

func TestI8_String(t *testing.T) {
	assertString(t, []stringAssert{
		{NewI8(-29), "-29"},
	})
}

func TestI8_Eq(t *testing.T) {
	assertEq(t, []eqAssert{
		{NewI8(23), NewI8(23), true},
		{NewI8(-23), NewI8(23), false},
		{NewI8(23), NewU8(23), false},
		{NewI8(23), NewBool(false), false},
	})
}

func TestI16_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewI16(0))
	assertRoundtrip(t, NewI16(12))
	assertRoundtrip(t, NewI16(-12))
}

func TestI16_EncodedLength(t *testing.T) {
	assertEncodedLength(t, []encodedLengthAssert{{NewI16(-13), 2}})
}

func TestI16_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{NewI16(-29), mustDecodeHexString("0xe3ff")},
	})
}

func TestI16_Hash(t *testing.T) {
	assertHash(t, []hashAssert{
		{NewI16(-29), mustDecodeHexString("0x39fbf34f574b72d1815c602a2fe95b7af4b5dfd7bc92a2fc0824aa55f8b9d7b2")},
	})
}

func TestI16_Hex(t *testing.T) {
	assertHex(t, []hexAssert{
		{NewI16(-29), "0xe3ff"},
	})
}

func TestI16_String(t *testing.T) {
	assertString(t, []stringAssert{
		{NewI16(-29), "-29"},
	})
}

func TestI16_Eq(t *testing.T) {
	assertEq(t, []eqAssert{
		{NewI16(23), NewI16(23), true},
		{NewI16(-23), NewI16(23), false},
		{NewI16(23), NewU16(23), false},
		{NewI16(23), NewBool(false), false},
	})
}

func TestI32_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewI32(0))
	assertRoundtrip(t, NewI32(12))
	assertRoundtrip(t, NewI32(-12))
}

func TestI32_EncodedLength(t *testing.T) {
	assertEncodedLength(t, []encodedLengthAssert{{NewI32(-13), 4}})
}

func TestI32_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{NewI32(-29), mustDecodeHexString("0xe3ffffff")},
	})
}

func TestI32_Hash(t *testing.T) {
	assertHash(t, []hashAssert{
		{NewI32(-29), mustDecodeHexString("0x6ef9d4772b9d657bfa727862d9690d5bf8b9045943279e95d3ae0743684f1b95")},
	})
}

func TestI32_Hex(t *testing.T) {
	assertHex(t, []hexAssert{
		{NewI32(-29), "0xe3ffffff"},
	})
}

func TestI32_String(t *testing.T) {
	assertString(t, []stringAssert{
		{NewI32(-29), "-29"},
	})
}

func TestI32_Eq(t *testing.T) {
	assertEq(t, []eqAssert{
		{NewI32(23), NewI32(23), true},
		{NewI32(-23), NewI32(23), false},
		{NewI32(23), NewU32(23), false},
		{NewI32(23), NewBool(false), false},
	})
}

func TestI64_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewI64(0))
	assertRoundtrip(t, NewI64(12))
	assertRoundtrip(t, NewI64(-12))
}

func TestI64_EncodedLength(t *testing.T) {
	assertEncodedLength(t, []encodedLengthAssert{{NewI64(-13), 8}})
}

func TestI64_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{NewI64(-29), mustDecodeHexString("0xe3ffffffffffffff")},
	})
}

func TestI64_Hash(t *testing.T) {
	assertHash(t, []hashAssert{
		{NewI64(-29), mustDecodeHexString("0x4d42db2aa4a23bde81a3ad3705220affaa457c56a0135080c71db7783fec8f44")},
	})
}

func TestI64_Hex(t *testing.T) {
	assertHex(t, []hexAssert{
		{NewI64(-29), "0xe3ffffffffffffff"},
	})
}

func TestI64_String(t *testing.T) {
	assertString(t, []stringAssert{
		{NewI64(-29), "-29"},
	})
}

func TestI64_Eq(t *testing.T) {
	assertEq(t, []eqAssert{
		{NewI64(23), NewI64(23), true},
		{NewI64(-23), NewI64(23), false},
		{NewI64(23), NewU64(23), false},
		{NewI64(23), NewBool(false), false},
	})
}
