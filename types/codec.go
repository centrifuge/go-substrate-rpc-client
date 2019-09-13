package types

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"reflect"

	"github.com/centrifuge/go-substrate-rpc-client/scale"
	"golang.org/x/crypto/blake2b"
)

func EncodeToBytes(value interface{}) ([]byte, error) { // TODO rename to Encode
	var buffer = bytes.Buffer{}
	err := scale.NewEncoder(&buffer).Encode(value)
	if err != nil {
		return buffer.Bytes(), err
	}
	return buffer.Bytes(), nil
}

func EncodeToHexString(value interface{}) (string, error) { // TODO rename to EncodeToHex
	bz, err := EncodeToBytes(value)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%#x", bz), nil
}

func DecodeFromBytes(bz []byte, target interface{}) error { // TODO rename to Decode
	return scale.NewDecoder(bytes.NewReader(bz)).Decode(target)
}

func DecodeFromHexString(str string, target interface{}) error { // TODO rename to DecodeFromHex
	bz, err := hex.DecodeString(str[2:])
	if err != nil {
		return err
	}
	return DecodeFromBytes(bz, target)
}

// EncodedLength returns the length of the value when encoded as a byte array
func EncodedLength(value interface{}) (int, error) {
	var buffer = bytes.Buffer{}
	err := scale.NewEncoder(&buffer).Encode(value)
	if err != nil {
		return 0, err
	}
	return buffer.Len(), nil
}

// GetHash returns a hash of the value
func GetHash(value interface{}) (Hash, error) {
	enc, err := EncodeToBytes(value)
	if err != nil {
		return Hash{}, err
	}
	return blake2b.Sum256(enc), err
}

// IsEmpty checks if the value is an empty value
func IsEmpty(value interface{}) bool {
	// TODO
	//for _, v := range value {
	//	if v != 0 {
	//		return false
	//	}
	//}
	//return true
	return false
}

// Eq compares the value of the input to see if there is a match
func Eq(one, other interface{}) bool {
	return reflect.DeepEqual(one, other)
}

// Hex returns a hex string representation of the value
func Hex(value interface{}) (string, error) {
	return EncodeToHexString(value)
}
