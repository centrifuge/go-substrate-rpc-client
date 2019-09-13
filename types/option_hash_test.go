package types

import (
	"testing"
)

func TestOptionH160_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionH160(NewH160(hash20)))
	assertRoundtrip(t, NewOptionH160Empty())
}

func TestOptionH256_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionH256(NewH256(hash32)))
	assertRoundtrip(t, NewOptionH256Empty())
}

func TestOptionH512_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionH512(NewH512(hash64)))
	assertRoundtrip(t, NewOptionH512Empty())
}

func TestOptionHash_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionHash(NewHash(hash32)))
	assertRoundtrip(t, NewOptionHashEmpty())
}
