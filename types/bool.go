package types

import (
	"fmt"

	"github.com/ChainSafe/gossamer/codec"
)

// Bool represents boolean values
type Bool bool

// NewBool creates a new Bool
func NewBool(b bool) Bool {
	return Bool(b)
}

func (b Bool) EncodedLength() (int, error) {
	return encodedLength(bool(b))
}

func (b Bool) Hash() ([32]byte, error) {
	return hash(b)
}

func (b Bool) IsEmpty() bool {
	return bool(b)
}

func (b Bool) Encode() ([]byte, error) {
	return codec.Encode(bool(b))
}

func (b Bool) Hex() (string, error) {
	return _hex(b)
}

func (b Bool) String() string {
	return fmt.Sprintf("%t", b)
}

func (b Bool) Eq(o Codec) bool {
	if ov, ok := o.(Bool); ok {
		return ov == b
	}
	return false
}
