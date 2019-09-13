package types

import (
	"testing"
)

func TestOptionU8_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionU8(NewU8(7)))
	assertRoundtrip(t, NewOptionU8(NewU8(0)))
	assertRoundtrip(t, NewOptionU8Empty())
}

func TestOptionU16_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionU16(NewU16(14)))
	assertRoundtrip(t, NewOptionU16(NewU16(0)))
	assertRoundtrip(t, NewOptionU16Empty())
}

func TestOptionU32_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionU32(NewU32(21)))
	assertRoundtrip(t, NewOptionU32(NewU32(0)))
	assertRoundtrip(t, NewOptionU32Empty())
}

func TestOptionU64_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionU64(NewU64(28)))
	assertRoundtrip(t, NewOptionU64(NewU64(0)))
	assertRoundtrip(t, NewOptionU64Empty())
}
