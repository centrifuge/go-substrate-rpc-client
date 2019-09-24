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
	"hash"

	"github.com/centrifuge/go-substrate-rpc-client/scale"
	"github.com/pierrec/xxHash/xxHash64"
)

type StorageKey []byte

func NewStorageKey(b []byte) StorageKey {
	return b
}

func CreateStorageKey(meta Metadata, module string, fn string, key []byte) (StorageKey, error) {
	var fnMeta *StorageFunctionMetadata
	for _, m := range meta.Metadata.Modules {
		if m.Prefix == module {
			for _, s := range m.Storage {
				if s.Name == fn {
					fnMeta = &s //nolint:scopelint
					break
				}
			}
		}
	}
	if fnMeta == nil {
		return nil, fmt.Errorf("no metadata found for module %s function %s", module, fn)
	}

	var hasher hash.Hash
	var err error
	if fnMeta.isMap() {
		hasher, err = fnMeta.Map.HashFunc()
		if err != nil {
			return nil, err
		}
	} else if fnMeta.isDMap() {
		// TODO define hashing for 2 keys
		return nil, fmt.Errorf("double map storage keys are not yet implemented")
	}

	afn := []byte(module + " " + fn)

	// TODO why is add length prefix step in JS client doesn't add anything to the hashed key?
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
	return createMultiXxhash(append(afn), 2)
}

// Encode implements encoding for StorageKey, which just unwraps the bytes of StorageKey
func (s StorageKey) Encode(encoder scale.Encoder) error {
	return encoder.Write(s)
}

// Decode implements decoding for StorageKey, which just wraps the bytes in StorageKey
func (s *StorageKey) Decode(decoder scale.Decoder) error {
	return decoder.Read(*s)
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
