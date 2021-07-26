// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
//
// Copyright 2019 Centrifuge GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	"fmt"
	"io"

	"github.com/snowfork/go-substrate-rpc-client/v3/scale"
	"github.com/snowfork/go-substrate-rpc-client/v3/xxhash"
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

// CreateStorageKeyPrefix creates a key prefix for keys of a map.
// Can be used as an input to the state.GetKeys() RPC, in order to list the keys of map.
func CreateStorageKeyPrefix(prefix, method string) []byte {
	return createPrefixedKey(method, prefix)
}

// CreateStorageKey uses the given metadata and to derive the right hashing of method, prefix as well as arguments to
// create a hashed StorageKey
// Using variadic argument, so caller do not need to construct array of arguments
func CreateStorageKey(meta *Metadata, prefix, method string, args ...[]byte) (StorageKey, error) {
	stringKey := []byte(prefix + " " + method)

	validateAndTrimArgs := func(args [][]byte) ([][]byte, error) {
		nonNilCount := -1
		for i, arg := range args {
			if len(arg) == 0 {
				nonNilCount = i
				break
			}
		}

		if nonNilCount == -1 {
			return args, nil
		}

		for i := nonNilCount; i < len(args); i++ {
			if len(args[i]) != 0 {
				return nil, fmt.Errorf("non-nil arguments cannot be preceded by nil arguments")
			}
		}

		trimmedArgs := make([][]byte, nonNilCount)
		for i := 0; i < nonNilCount; i++ {
			trimmedArgs[i] = args[i]
		}

		return trimmedArgs, nil
	}

	validatedArgs, err := validateAndTrimArgs(args)
	if err != nil {
		return nil, err
	}

	entryMeta, err := meta.FindStorageEntryMetadata(prefix, method)
	if err != nil {
		return nil, err
	}

	if entryMeta.IsNMap() {
		hashers, err := entryMeta.Hashers()
		if err != nil {
			return nil, fmt.Errorf("unable to get hashers for %s nmap", method)
		}
		if len(hashers) != len(validatedArgs) {
			return nil, fmt.Errorf("%s:%s is a nmap, therefore requires that number of arguments should "+
				"exactly match number of hashers in metadata. "+
				"Expected: %d, received: %d", prefix, method, len(hashers), len(validatedArgs))
		}
		return createKeyNMap(method, prefix, validatedArgs, entryMeta)
	}

	if entryMeta.IsDoubleMap() {
		if len(validatedArgs) != 2 {
			return nil, fmt.Errorf("%s:%s is a double map, therefore requires precisely two arguments. "+
				"received: %d", prefix, method, len(validatedArgs))
		}
		return createKeyDoubleMap(meta, method, prefix, stringKey, validatedArgs[0], validatedArgs[1], entryMeta)
	}

	if entryMeta.IsMap() {
		if len(validatedArgs) != 1 {
			return nil, fmt.Errorf("%s:%s is a map, therefore requires precisely one argument. "+
				"received: %d", prefix, method, len(validatedArgs))
		}
		return createKey(meta, method, prefix, stringKey, validatedArgs[0], entryMeta)
	}

	if entryMeta.IsPlain() && len(validatedArgs) != 0 {
		return nil, fmt.Errorf("%s:%s is a plain key, therefore requires no argument. "+
			"received: %d", prefix, method, len(validatedArgs))
	}

	return createKey(meta, method, prefix, stringKey, nil, entryMeta)
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

func createKeyNMap(method, prefix string, args [][]byte, entryMeta StorageEntryMetadata) (StorageKey, error) {
	hashers, err := entryMeta.Hashers()
	if err != nil {
		return nil, err
	}

	key := createPrefixedKey(method, prefix)

	for i, arg := range args {
		_, err := hashers[i].Write(arg)
		if err != nil {
			return nil, fmt.Errorf("unable to hash args[%d]: %s Error: %v", i, arg, err)
		}
		key = append(key, hashers[i].Sum(nil)...)
	}

	return key, nil
}

// createKeyDoubleMap creates a key for a DoubleMap type
func createKeyDoubleMap(meta *Metadata, method, prefix string, stringKey, arg, arg2 []byte,
	entryMeta StorageEntryMetadata) (StorageKey, error) {

	hasher, err := entryMeta.Hasher()
	if err != nil {
		return nil, err
	}

	hasher2, err := entryMeta.Hasher2()
	if err != nil {
		return nil, err
	}

	if meta.Version <= 8 {
		_, err := hasher.Write(append(stringKey, arg...))
		if err != nil {
			return nil, err
		}
		_, err = hasher2.Write(arg2)
		if err != nil {
			return nil, err
		}
		return append(hasher.Sum(nil), hasher2.Sum(nil)...), err
	}

	_, err = hasher.Write(arg)
	if err != nil {
		return nil, err
	}
	_, err = hasher2.Write(arg2)
	if err != nil {
		return nil, err
	}

	key := createPrefixedKey(method, prefix)
	key = append(key, hasher.Sum(nil)...)
	key = append(key, hasher2.Sum(nil)...)

	return key, nil
}

// createKey creates a key for either a map or a plain value
func createKey(meta *Metadata, method, prefix string, stringKey, arg []byte, entryMeta StorageEntryMetadata) (
	StorageKey, error) {

	hasher, err := entryMeta.Hasher()
	if err != nil {
		return nil, err
	}

	if meta.Version <= 8 {
		_, err := hasher.Write(append(stringKey, arg...))
		return hasher.Sum(nil), err
	}

	if entryMeta.IsMap() {
		_, err := hasher.Write(arg)
		if err != nil {
			return nil, err
		}
		arg = hasher.Sum(nil)
	}

	return append(createPrefixedKey(method, prefix), arg...), nil
}

func createPrefixedKey(method, prefix string) []byte {
	return append(xxhash.New128([]byte(prefix)).Sum(nil), xxhash.New128([]byte(method)).Sum(nil)...)
}

