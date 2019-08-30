package codec

import (
	"bytes"
	"fmt"

	"github.com/ChainSafe/gossamer/codec"
)

// Bool represents boolean values
type Bytes []byte

// NewBytes creates a new Bytes type
func NewBytes(b []byte) Bytes {
	return Bytes(b)
}

func (b Bytes) EncodedLength() (int, error) {
	return encodedLength([]byte(b))
}

func (b Bytes) Hash() ([32]byte, error) {
	return hash(b)
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
	return _hex(b)
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
