package types

import (
	"testing"
)

var hash20 = [20]byte{
	1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
}

func TestH160EncodedLength(t *testing.T) {
	testEncodedLength(t, []encodedLengthTest{{NewH160(hash20), 21}})
}

func TestH160Encode(t *testing.T) {
	testEncode(t, []encodingTest{
		{NewH160(hash20), mustDecodeHexString("0x500102030405060708090001020304050607080900")},
	})
}

func TestH160Hash(t *testing.T) {
	testHash(t, []hashTest{
		{NewH160(hash20), mustDecodeHexString("0x8cfcee28b5f749ec8bad9c058abb739942fccc5498bcb8b7cfa660ea2d3994b0")},
	})
}

func TestH160Hex(t *testing.T) {
	testHex(t, []hexTest{
		{NewH160(hash20), "0x500102030405060708090001020304050607080900"},
	})
}

func TestH160String(t *testing.T) {
	testString(t, []stringTest{
		{NewH160(hash20), "0x0102030405060708090001020304050607080900"},
	})
}

func TestH160Eq(t *testing.T) {
	testEq(t, []eqTest{
		{NewH160(hash20), NewH160(hash20), true},
		{NewH160(hash20), NewBytes(hash20[:]), false},
		{NewH160(hash20), NewBool(false), false},
	})
}

var hash32 = [32]byte{
	1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
	1, 2,
}

func TestH256EncodedLength(t *testing.T) {
	testEncodedLength(t, []encodedLengthTest{{NewH256(hash32), 33}})
}

func TestH256Encode(t *testing.T) {
	testEncode(t, []encodingTest{
		{NewH256(hash32), mustDecodeHexString("0x800102030405060708090001020304050607080900010203040506070809000102")},
	})
}

func TestH256Hash(t *testing.T) {
	testHash(t, []hashTest{
		{NewH256(hash32), mustDecodeHexString("0xde5b09770bf1e1f93bf1a11c3fb060affc6cb8658f33154ce53629a3752954d6")},
	})
}

func TestH256Hex(t *testing.T) {
	testHex(t, []hexTest{
		{NewH256(hash32), "0x800102030405060708090001020304050607080900010203040506070809000102"},
	})
}

func TestH256String(t *testing.T) {
	testString(t, []stringTest{
		{NewH256(hash32), "0x0102030405060708090001020304050607080900010203040506070809000102"},
	})
}

func TestH256Eq(t *testing.T) {
	testEq(t, []eqTest{
		{NewH256(hash32), NewH256(hash32), true},
		{NewH256(hash32), NewBytes(hash32[:]), false},
		{NewH256(hash32), NewBool(false), false},
	})
}

func TestHashEncodedLength(t *testing.T) {
	testEncodedLength(t, []encodedLengthTest{{NewHash(hash32), 33}})
}

func TestHashEncode(t *testing.T) {
	testEncode(t, []encodingTest{
		{NewHash(hash32), mustDecodeHexString("0x800102030405060708090001020304050607080900010203040506070809000102")},
	})
}

func TestHashHash(t *testing.T) {
	testHash(t, []hashTest{
		{NewHash(hash32), mustDecodeHexString("0xde5b09770bf1e1f93bf1a11c3fb060affc6cb8658f33154ce53629a3752954d6")},
	})
}

func TestHashHex(t *testing.T) {
	testHex(t, []hexTest{
		{NewHash(hash32), "0x800102030405060708090001020304050607080900010203040506070809000102"},
	})
}

func TestHashString(t *testing.T) {
	testString(t, []stringTest{
		{NewHash(hash32), "0x0102030405060708090001020304050607080900010203040506070809000102"},
	})
}

func TestHashEq(t *testing.T) {
	testEq(t, []eqTest{
		{NewHash(hash32), NewHash(hash32), true},
		{NewHash(hash32), NewBytes(hash32[:]), false},
		{NewHash(hash32), NewBool(false), false},
	})
}

var hash64 = [64]byte{
	1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
	1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
	1, 2, 3, 4,
}

func TestH512EncodedLength(t *testing.T) {
	testEncodedLength(t, []encodedLengthTest{{NewH512(hash64), 66}})
}

func TestH512Encode(t *testing.T) {
	testEncode(t, []encodingTest{
		{NewH512(hash64), mustDecodeHexString("0x010101020304050607080900010203040506070809000102030405060708090001020304050607080900010203040506070809000102030405060708090001020304")}, //nolint:lll
	})
}

func TestH512Hash(t *testing.T) {
	testHash(t, []hashTest{
		{NewH512(hash64), mustDecodeHexString("0x0926d23398a248b1c7723651a5ad05a5626cc8f9450512d6c3b5b2156615bcd5")},
	})
}

func TestH512Hex(t *testing.T) {
	testHex(t, []hexTest{
		{NewH512(hash64), "0x010101020304050607080900010203040506070809000102030405060708090001020304050607080900010203040506070809000102030405060708090001020304"}, //nolint:lll
	})
}

func TestH512String(t *testing.T) {
	testString(t, []stringTest{
		{NewH512(hash64), "0x01020304050607080900010203040506070809000102030405060708090001020304050607080900010203040506070809000102030405060708090001020304"}, //nolint:lll
	})
}

func TestH512Eq(t *testing.T) {
	testEq(t, []eqTest{
		{NewH512(hash64), NewH512(hash64), true},
		{NewH512(hash64), NewBytes(hash64[:]), false},
		{NewH512(hash64), NewBool(false), false},
	})
}
