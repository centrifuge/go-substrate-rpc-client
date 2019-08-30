package codec

import (
	"bytes"
	"fmt"

	"github.com/ChainSafe/gossamer/codec"
	"golang.org/x/crypto/blake2b"
)

func encodedLength(in interface{}) (int, error) {
	buffer := bytes.Buffer{}
	se := codec.Encoder{&buffer}
	len, err := se.Encode(in)
	return len, err
}

func hash(in Codec) (hash [32]byte, err error) {
	enc, err := in.Encode()
	if err != nil {
		return hash, err
	}
	return blake2b.Sum256(enc), err
}

func _hex(in Codec) (string, error) {
	e, err := in.Encode()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", e), nil
}

func isEmpty(b []byte) bool {
	for _, v := range b {
		if v != 0 {
			return false
		}
	}
	return true
}
