// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
// Copyright (C) 2019  Centrifuge GmbH
//
// This file is part of Go Substrate RPC Client (GSRPC).
//
// GSRPC is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// GSRPC is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package types

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"reflect"

	"github.com/centrifuge/go-substrate-rpc-client/scale"
	"golang.org/x/crypto/blake2b"
)

type Hexer interface {
	Hex() string
}

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

// Eq compares the value of the input to see if there is a match
func Eq(one, other interface{}) bool {
	return reflect.DeepEqual(one, other)
}

// Hex returns a hex string representation of the value
func Hex(value interface{}) (string, error) {
	switch v := value.(type) {
	case Hexer:
		return v.Hex(), nil
	case []byte:
		return fmt.Sprintf("%#x", v), nil
	default:
		return "", fmt.Errorf("does not support %T", v)
	}
}
