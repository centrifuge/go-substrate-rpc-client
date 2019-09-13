package types

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/centrifuge/go-substrate-rpc-client/scale"
)

// Codec is the base interface that all types implement. The Codec Base is required for operating as an
// encoding/decoding layer.
type Codec interface {
	// The length of the value when encoded as a byte array
	EncodedLength() (int, error)
	// Returns a hash of the value
	Hash() ([32]byte, error)
	// Checks if the value is an empty value
	IsEmpty() bool
	// Compares the value of the input to see if there is a match
	Eq(o Codec) bool
	// Returns a hex string representation of the value. isLe returns a LE (number-only) representation
	Hex() (string, error)
	// Returns the string representation of the value
	String() string
	// Encodes the value as a byte array as per the SCALE specifications
	Encode() ([]byte, error)
}

func EncodeToBytes(value interface{}) ([]byte, error) {
	var buffer = bytes.Buffer{}
	err := scale.NewEncoder(&buffer).Encode(value)
	if err != nil {
		return buffer.Bytes(), err
	}
	return buffer.Bytes(), nil
}

func EncodeToHexString(value interface{}) (string, error) {
	bz, err := EncodeToBytes(value)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%#x", bz), nil
}

func DecodeFromBytes(bz []byte, target interface{}) error {
	return scale.NewDecoder(bytes.NewReader(bz)).Decode(target)
}

func DecodeFromHexString(str string, target interface{}) error {
	bz, err := hex.DecodeString(str[2:])
	if err != nil {
		return err
	}
	return DecodeFromBytes(bz, target)
}
