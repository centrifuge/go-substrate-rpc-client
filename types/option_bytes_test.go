package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptionBytes8_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionBytes8(NewBytes8([8]byte{12})))
	assertRoundtrip(t, NewOptionBytes8(NewBytes8([8]byte{})))
	assertRoundtrip(t, NewOptionBytes8Empty())
}

func TestOptionBytes8_EncodedLength(t *testing.T) {
	assertEncodedLength(t, []encodedLengthAssert{
		{NewOptionBytes8(NewBytes8([8]byte{})), 10},
		{NewOptionBytes8(NewBytes8([8]byte{7, 6, 5, 4, 3, 2, 1, 0})), 10},
		{NewOptionBytes8Empty(), 1},
	})
}

func TestOptionBytes8_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{NewOptionBytes8(NewBytes8([8]byte{0, 0, 0})), mustDecodeHexString("0x01200000000000000000")},
		{NewOptionBytes8(NewBytes8([8]byte{171, 18, 52})), mustDecodeHexString("0x0120ab12340000000000")},
		{NewOptionBytes8Empty(), mustDecodeHexString("0x00")},
	})
}

func TestOptionBytes8_Hash(t *testing.T) {
	assertHash(t, []hashAssert{
		{NewOptionBytes8(NewBytes8([8]byte{0, 42, 254})), mustDecodeHexString(
			"0x654f5633f50b9db847bbffe819bf119952ebb5fcfee02561e8c4981bf769b351")},
		{NewOptionBytes8(NewBytes8([8]byte{0, 0})), mustDecodeHexString(
			"0x774d708a6d4a5a0122f6e8def11bfc4e37e1d29496b3673cd3fbdef1762f2406")},
		{NewOptionBytes8Empty(), mustDecodeHexString(
			"0x03170a2e7597b7b7e3d84c05391d139a62b157e78786d8c082f29dcf4c111314")},
	})
}

func TestOptionBytes8_Eq(t *testing.T) {
	assertEq(t, []eqAssert{
		{NewOptionBytes8(NewBytes8([8]byte{1, 0, 0})), NewBytes8([8]byte{1, 0}), false},
		{NewOptionBytes8(NewBytes8([8]byte{0, 0, 1})), NewOptionBytes8(NewBytes8([8]byte{0, 0, 1})), true},
		{NewOptionBytes8Empty(), NewOptionBytes8Empty(), true},
	})
}

func TestOptionBytes8(t *testing.T) {
	bz := NewOptionBytes8(NewBytes8([8]byte{1, 0, 0}))
	assert.False(t, bz.IsNone())
	assert.True(t, bz.IsSome())
	ok, val := bz.Unwrap()
	assert.True(t, ok)
	assert.Equal(t, val, NewBytes8([8]byte{1, 0, 0}))
	bz.SetNone()
	assert.True(t, bz.IsNone())
	assert.False(t, bz.IsSome())
	ok2, val2 := bz.Unwrap()
	assert.False(t, ok2)
	assert.Equal(t, val2, NewBytes8([8]byte{}))
	bz.SetSome(NewBytes8([8]byte{3}))
	assert.False(t, bz.IsNone())
	assert.True(t, bz.IsSome())
	ok3, val3 := bz.Unwrap()
	assert.True(t, ok3)
	assert.Equal(t, val3, NewBytes8([8]byte{3}))
}

func TestOptionBytes_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionBytes(NewBytes([]byte{12})))
	assertRoundtrip(t, NewOptionBytes(NewBytes([]byte{2})))
	assertRoundtrip(t, NewOptionBytesEmpty())
}

func TestOptionBytes_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{NewOptionBytes(NewBytes([]byte{0, 0, 0})), mustDecodeHexString("0x010c000000")},
		{NewOptionBytes(NewBytes([]byte{171, 1, 52})), mustDecodeHexString("0x010cab0134")},
		{NewOptionBytesEmpty(), mustDecodeHexString("0x00")},
	})
}

func TestOptionBytes_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{mustDecodeHexString("0x010c000000"), NewOptionBytes(NewBytes([]byte{0, 0, 0}))},
		{mustDecodeHexString("0x010cab0134"), NewOptionBytes(NewBytes([]byte{171, 1, 52}))},
		{mustDecodeHexString("0x00"), NewOptionBytesEmpty()},
	})
}

func TestOptionBytes16_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionBytes16(NewBytes16([16]byte{12})))
	assertRoundtrip(t, NewOptionBytes16(NewBytes16([16]byte{})))
	assertRoundtrip(t, NewOptionBytes16Empty())
}

func TestOptionBytes32_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionBytes32(NewBytes32([32]byte{12})))
	assertRoundtrip(t, NewOptionBytes32(NewBytes32([32]byte{})))
	assertRoundtrip(t, NewOptionBytes32Empty())
}

func TestOptionBytes64_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionBytes64(NewBytes64([64]byte{12})))
	assertRoundtrip(t, NewOptionBytes64(NewBytes64([64]byte{})))
	assertRoundtrip(t, NewOptionBytes64Empty())
}

func TestOptionBytes128_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionBytes128(NewBytes128([128]byte{12})))
	assertRoundtrip(t, NewOptionBytes128(NewBytes128([128]byte{})))
	assertRoundtrip(t, NewOptionBytes128Empty())
}

func TestOptionBytes256_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionBytes256(NewBytes256([256]byte{12})))
	assertRoundtrip(t, NewOptionBytes256(NewBytes256([256]byte{})))
	assertRoundtrip(t, NewOptionBytes256Empty())
}

func TestOptionBytes512_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionBytes512(NewBytes512([512]byte{12})))
	assertRoundtrip(t, NewOptionBytes512(NewBytes512([512]byte{})))
	assertRoundtrip(t, NewOptionBytes512Empty())
}

func TestOptionBytes1024_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionBytes1024(NewBytes1024([1024]byte{12})))
	assertRoundtrip(t, NewOptionBytes1024(NewBytes1024([1024]byte{})))
	assertRoundtrip(t, NewOptionBytes1024Empty())
}

func TestOptionBytes2048_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionBytes2048(NewBytes2048([2048]byte{12})))
	assertRoundtrip(t, NewOptionBytes2048(NewBytes2048([2048]byte{})))
	assertRoundtrip(t, NewOptionBytes2048Empty())
}
