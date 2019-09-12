package types

import (
	"bytes"
	"fmt"

	"github.com/ChainSafe/gossamer/codec"
)

type H160 [20]byte

func NewH160(b [20]byte) H160 {
	return H160(b)
}

func (b H160) EncodedLength() (int, error) {
	return encodedLength(b[:])
}

func (b H160) Hash() ([32]byte, error) {
	return hash(b)
}

func (b H160) IsEmpty() bool {
	return isEmpty(b[:])
}

func (b H160) Encode() ([]byte, error) {
	return codec.Encode(b[:])
}

func (b H160) Hex() (string, error) {
	return _hex(b)
}

func (b H160) String() string {
	return fmt.Sprintf("%#x", b[:])
}

func (b H160) Eq(o Codec) bool {
	if ov, ok := o.(H160); ok {
		return bytes.Equal(b[:], ov[:])
	}
	return false
}

type H256 [32]byte

func NewH256(b [32]byte) H256 {
	return H256(b)
}

func (b H256) EncodedLength() (int, error) {
	return encodedLength(b[:])
}

func (b H256) Hash() ([32]byte, error) {
	return hash(b)
}

func (b H256) IsEmpty() bool {
	return isEmpty(b[:])
}

func (b H256) Encode() ([]byte, error) {
	return codec.Encode(b[:])
}

func (b H256) Hex() (string, error) {
	return _hex(b)
}

func (b H256) String() string {
	return fmt.Sprintf("%#x", b[:])
}

func (b H256) Eq(o Codec) bool {
	if ov, ok := o.(H256); ok {
		return bytes.Equal(b[:], ov[:])
	}
	return false
}

type H512 [64]byte

func NewH512(b [64]byte) H512 {
	return H512(b)
}

func (b H512) EncodedLength() (int, error) {
	return encodedLength(b[:])
}

func (b H512) Hash() ([32]byte, error) {
	return hash(b)
}

func (b H512) IsEmpty() bool {
	return isEmpty(b[:])
}

func (b H512) Encode() ([]byte, error) {
	return codec.Encode(b[:])
}

func (b H512) Hex() (string, error) {
	return _hex(b)
}

func (b H512) String() string {
	return fmt.Sprintf("%#x", b[:])
}

func (b H512) Eq(o Codec) bool {
	if ov, ok := o.(H512); ok {
		return bytes.Equal(b[:], ov[:])
	}
	return false
}

type Hash H256

func NewHash(b [32]byte) Hash {
	return Hash(b)
}

func (b Hash) EncodedLength() (int, error) {
	return encodedLength(b[:])
}

func (b Hash) Hash() ([32]byte, error) {
	return hash(b)
}

func (b Hash) IsEmpty() bool {
	return isEmpty(b[:])
}

func (b Hash) Encode() ([]byte, error) {
	return codec.Encode(b[:])
}

func (b Hash) Hex() (string, error) {
	return _hex(b)
}

func (b Hash) String() string {
	return fmt.Sprintf("%#x", b[:])
}

func (b Hash) Eq(o Codec) bool {
	if ov, ok := o.(Hash); ok {
		return bytes.Equal(b[:], ov[:])
	}
	return false
}
