package codec

import (
	"bytes"
	"fmt"

	"github.com/ChainSafe/gossamer/codec"
	"golang.org/x/crypto/blake2b"
)

// Bool represents boolean values
type Bytes []byte

// NewBytes creates a new Bytes type
func NewBytes(b []byte) Bytes {
	return Bytes(b)
}

func (b Bytes) EncodedLength() (int, error) {
	buffer := bytes.Buffer{}
	se := codec.Encoder{&buffer}
	len, err := se.Encode([]byte(b))
	return len, err
}

func (b Bytes) Hash() (hash [32]byte, err error) {
	e, err := b.Encode()
	if err != nil {
		return hash, err
	}
	return blake2b.Sum256(e), err
}

func (b Bytes) IsEmpty() bool {
	for _, v := range b {
		if v != 0 {
			return false
		}
	}
	return true
}

func (b Bytes) Encode() ([]byte, error) {
	return codec.Encode([]byte(b))
}

func (b Bytes) Hex() (string, error) {
	e, err := b.Encode()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", e), nil
}

func (b Bytes) String() string {
	return fmt.Sprintf("%x", []byte(b))
}

func (b Bytes) Eq(o Codec) bool {
	if ov, ok := o.(Bytes); ok {
		return bytes.Equal(b, ov)
	}
	return false
}
