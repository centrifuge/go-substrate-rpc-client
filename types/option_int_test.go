package types

import (
	"testing"
)

func TestOptionI8_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionI8(NewI8(7)))
	assertRoundtrip(t, NewOptionI8(NewI8(0)))
	assertRoundtrip(t, NewOptionI8Empty())
}

func TestOptionI16_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionI16(NewI16(14)))
	assertRoundtrip(t, NewOptionI16(NewI16(0)))
	assertRoundtrip(t, NewOptionI16Empty())
}

func TestOptionI32_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionI32(NewI32(21)))
	assertRoundtrip(t, NewOptionI32(NewI32(0)))
	assertRoundtrip(t, NewOptionI32Empty())
}

func TestOptionI64_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionI64(NewI64(28)))
	assertRoundtrip(t, NewOptionI64(NewI64(0)))
	assertRoundtrip(t, NewOptionI64Empty())
}
