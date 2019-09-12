package types

import (
	"github.com/ChainSafe/gossamer/codec"
)

type U8 uint8

func NewU8(u uint8) U8 {
	return U8(u)
}

func (u U8) EncodedLength() (int, error) {
	return encodedLength(u)
}

func (u U8) Hash() ([32]byte, error) {
	return hash(u)
}

func (u U8) IsEmpty() bool {
	return u == 0
}

func (u U8) Encode() ([]byte, error) {
	return codec.Encode(u)
}

func (u U8) Hex() (string, error) {
	return _hex(u)
}

func (u U8) String() string {
	return string(u)
}

func (u U8) Eq(o Codec) bool {
	if ov, ok := o.(U8); ok {
		return u == ov
	}
	return false
}
