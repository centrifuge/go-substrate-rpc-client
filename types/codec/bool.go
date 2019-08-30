package codec

import (
	"bytes"
	"fmt"

	"github.com/ChainSafe/gossamer/codec"
	"golang.org/x/crypto/blake2b"
)

// Bool represents boolean values
type Bool bool

// NewBool creates a new Bool
func NewBool(b bool) Bool {
	return Bool(b)
}

func (b Bool) EncodedLength() (int, error) {
	buffer := bytes.Buffer{}
	se := codec.Encoder{&buffer}
	len, err := se.Encode(bool(b))
	return len, err
}

func (b Bool) Hash() (hash [32]byte, err error) {
	e, err := b.Encode()
	if err != nil {
		return hash, err
	}
	return blake2b.Sum256(e), err
}

func (b Bool) IsEmpty() bool {
	return bool(b)
}

func (b Bool) Encode() ([]byte, error) {
	return codec.Encode(bool(b))
}

func (b Bool) Hex() (string, error) {
	e, err := b.Encode()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", e), nil
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
