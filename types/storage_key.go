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
	"fmt"
	"io"

	"github.com/centrifuge/go-substrate-rpc-client/scale"
	"github.com/pierrec/xxHash/xxHash64"
)

// StorageKey represents typically hashed storage keys of the system.
// Be careful using this in your own structs – it only works as the last value in a struct since it will consume the
// remainder of the encoded data. The reason for this is that it does not contain any length encoding, so it would
// not know where to stop.
type StorageKey []byte

// NewStorageKey creates a new StorageKey type
func NewStorageKey(b []byte) StorageKey {
	return b
}

// CreateStorageKey uses the given metadata and to derive the right hashing of module, fn names and keys to create a
// hashed StorageKey
func CreateStorageKey(meta *Metadata, module string, fn string, key []byte) (StorageKey, error) {
	hasher, err := meta.FindStorageKeyHasher(module, fn)
	if err != nil {
		return nil, err
	}

	afn := []byte(module + " " + fn)

	if hasher != nil {
		_, err = hasher.Write(append(afn, key...))
		if err != nil {
			return nil, err
		}
		return hasher.Sum(nil), nil
	}

	if key != nil {
		return createMultiXxhash(append(afn, key...), 2)
	}
	return createMultiXxhash(afn, 2)
}

// Encode implements encoding for StorageKey, which just unwraps the bytes of StorageKey
func (s StorageKey) Encode(encoder scale.Encoder) error {
	return encoder.Write(s)
}

// Decode implements decoding for StorageKey, which just reads all the remaining bytes into StorageKey
func (s *StorageKey) Decode(decoder scale.Decoder) error {
	for i := 0; true; i++ {
		b, err := decoder.ReadOneByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		*s = append((*s)[:i], b)
	}
	return nil
}

// Hex returns a hex string representation of the value (not of the encoded value)
func (s StorageKey) Hex() string {
	return fmt.Sprintf("%#x", s)
}

func createMultiXxhash(data []byte, rounds int) ([]byte, error) {
	res := make([]byte, 0)
	for i := 0; i < rounds; i++ {
		h := xxHash64.New(uint64(i))
		_, err := h.Write(data)
		if err != nil {
			return nil, err
		}
		res = append(res, h.Sum(nil)...)
	}
	return res, nil
}
