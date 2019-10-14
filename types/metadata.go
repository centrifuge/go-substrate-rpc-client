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
	"errors"
	"hash"

	"github.com/centrifuge/go-substrate-rpc-client/scale"
	"golang.org/x/crypto/blake2b"
)

// Modelled after https://github.com/paritytech/substrate/blob/v1.0.0rc2/srml/metadata/src/lib.rs

type Metadata struct {
	MagicNumber uint32
	Version     uint8
	Metadata    RuntimeMetadataV4
}

func NewMetadata() *Metadata {
	return &Metadata{Version: 4, Metadata: RuntimeMetadataV4{make([]ModuleMetadata, 0)}}
}

func (m *Metadata) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&m.MagicNumber)
	if err != nil {
		return err
	}
	// TODO: we need to decide which struct to use based on the following number(enum), for now its hardcoded
	err = decoder.Decode(&m.Version)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.Metadata)
	if err != nil {
		return err
	}

	return nil
}

func (m Metadata) Encode(encoder scale.Encoder) error {
	err := encoder.Encode(m.MagicNumber)
	if err != nil {
		return err
	}
	// TODO: we need to decide which struct to use based on the following number(enum), for now its hardcoded
	err = encoder.Encode(m.Version)
	if err != nil {
		return err
	}

	err = encoder.Encode(m.Metadata)
	if err != nil {
		return err
	}

	return nil
}

type RuntimeMetadataV4 struct {
	Modules []ModuleMetadata
}

func (m *RuntimeMetadataV4) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&m.Modules)
	if err != nil {
		return err
	}
	return nil
}

func (m RuntimeMetadataV4) Encode(encoder scale.Encoder) error {
	err := encoder.Encode(m.Modules)
	if err != nil {
		return err
	}
	return nil
}

type ModuleMetadata struct {
	Name       string
	Prefix     string
	HasStorage bool
	Storage    []StorageFunctionMetadata
	HasCalls   bool
	Calls      []FunctionMetadata
	HasEvents  bool
	Events     []EventMetadata
}

func (m *ModuleMetadata) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&m.Name)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.Prefix)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.HasStorage)
	if err != nil {
		return err
	}

	if m.HasStorage {
		err = decoder.Decode(&m.Storage)
		if err != nil {
			return err
		}
	}

	err = decoder.Decode(&m.HasCalls)
	if err != nil {
		return err
	}

	if m.HasCalls {
		err = decoder.Decode(&m.Calls)
		if err != nil {
			return err
		}
	}

	err = decoder.Decode(&m.HasEvents)
	if err != nil {
		return err
	}

	if m.HasEvents {
		err = decoder.Decode(&m.Events)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m ModuleMetadata) Encode(encoder scale.Encoder) error {
	err := encoder.Encode(m.Name)
	if err != nil {
		return err
	}

	err = encoder.Encode(m.Prefix)
	if err != nil {
		return err
	}

	err = encoder.Encode(m.HasStorage)
	if err != nil {
		return err
	}

	if m.HasStorage {
		err = encoder.Encode(m.Storage)
		if err != nil {
			return err
		}
	}

	err = encoder.Encode(m.HasCalls)
	if err != nil {
		return err
	}

	if m.HasCalls {
		err = encoder.Encode(m.Calls)
		if err != nil {
			return err
		}
	}

	err = encoder.Encode(m.HasEvents)
	if err != nil {
		return err
	}

	if m.HasEvents {
		err = encoder.Encode(m.Events)
		if err != nil {
			return err
		}
	}
	return nil
}

type StorageFunctionMetadata struct {
	Name          string
	Modifier      uint8
	Type          uint8
	Plane         string       // TODO: rename to: Plain         string
	Map           TypMap       // TODO: rename to: Map           MapType
	DMap          TypDoubleMap // TODO: rename to: DoubleMap     DoubleMapType
	Fallback      []byte
	Documentation []string
}

func (s StorageFunctionMetadata) isMap() bool {
	return s.Type == 1
}

func (s StorageFunctionMetadata) isDMap() bool {
	return s.Type == 2
}

func (s *StorageFunctionMetadata) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&s.Name)
	if err != nil {
		return err
	}

	err = decoder.Decode(&s.Modifier)
	if err != nil {
		return err
	}

	err = decoder.Decode(&s.Type)
	if err != nil {
		return err
	}

	switch s.Type {
	case 0:
		err = decoder.Decode(&s.Plane)
		if err != nil {
			return err
		}
	case 1:
		err = decoder.Decode(&s.Map)
		if err != nil {
			return err
		}
	default:
		err = decoder.Decode(&s.DMap)
		if err != nil {
			return err
		}
	}
	err = decoder.Decode(&s.Fallback)
	if err != nil {
		return err
	}

	err = decoder.Decode(&s.Documentation)
	if err != nil {
		return err
	}
	// log.Println(metadataVersioned.Documentation)
	return nil
}

func (s StorageFunctionMetadata) Encode(encoder scale.Encoder) error {
	err := encoder.Encode(s.Name)
	if err != nil {
		return err
	}

	err = encoder.Encode(s.Modifier)
	if err != nil {
		return err
	}

	err = encoder.Encode(s.Type)
	if err != nil {
		return err
	}

	switch s.Type {
	case 0:
		err = encoder.Encode(s.Plane)
		if err != nil {
			return err
		}
	case 1:
		err = encoder.Encode(s.Map)
		if err != nil {
			return err
		}
	default:
		err = encoder.Encode(s.DMap)
		if err != nil {
			return err
		}
	}
	err = encoder.Encode(s.Fallback)
	if err != nil {
		return err
	}

	err = encoder.Encode(s.Documentation)
	if err != nil {
		return err
	}

	return nil
}

type FunctionMetadata struct {
	Name          string
	Args          []FunctionArgumentMetadata
	Documentation []string
}

func (m *FunctionMetadata) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&m.Name)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.Args)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.Documentation)
	if err != nil {
		return err
	}

	return nil
}

func (m FunctionMetadata) Encode(encoder scale.Encoder) error {
	err := encoder.Encode(m.Name)
	if err != nil {
		return err
	}

	err = encoder.Encode(m.Args)
	if err != nil {
		return err
	}

	err = encoder.Encode(m.Documentation)
	if err != nil {
		return err
	}

	return nil
}

type EventMetadata struct {
	Name          string
	Args          []string
	Documentation []string
}

func (m *EventMetadata) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&m.Name)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.Args)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.Documentation)
	if err != nil {
		return err
	}

	return nil
}

func (m EventMetadata) Encode(encoder scale.Encoder) error {
	err := encoder.Encode(m.Name)
	if err != nil {
		return err
	}

	err = encoder.Encode(m.Args)
	if err != nil {
		return err
	}

	err = encoder.Encode(m.Documentation)
	if err != nil {
		return err
	}

	return nil
}

type TypMap struct {
	Hasher   uint8
	Key      string
	Value    string
	IsLinked bool
}

func (t TypMap) HashFunc() (hash.Hash, error) {
	// Blake2_128
	// if t.Hasher == 0 {
	// 	// TODO implement Blake2_128
	// }

	// Blake2_256
	if t.Hasher == 1 {
		return blake2b.New256(nil)
	}

	// Twox128
	// if t.Hasher == 2 {
	// 	// TODO implement Twox128
	// }

	// Twox256
	// if t.Hasher == 3 {
	// 	// TODO implement Twox256
	// }

	// Twox64Concat
	// if t.Hasher == 4 {
	// 	// TODO implement Twox64Concat
	// }

	return nil, errors.New("hash function type not yet supported")
}

func (t *TypMap) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&t.Hasher)
	if err != nil {
		return err
	}

	err = decoder.Decode(&t.Key)
	if err != nil {
		return err
	}

	err = decoder.Decode(&t.Value)
	if err != nil {
		return err
	}

	err = decoder.Decode(&t.IsLinked)
	if err != nil {
		return err
	}

	return nil
}

func (t TypMap) Encode(encoder scale.Encoder) error {
	err := encoder.Encode(t.Hasher)
	if err != nil {
		return err
	}

	err = encoder.Encode(t.Key)
	if err != nil {
		return err
	}

	err = encoder.Encode(t.Value)
	if err != nil {
		return err
	}

	err = encoder.Encode(t.IsLinked)
	if err != nil {
		return err
	}

	return nil
}

type TypDoubleMap struct {
	Hasher     uint8
	Key        string
	Key2       string
	Value      string
	Key2Hasher string
}

func (m *TypDoubleMap) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&m.Hasher)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.Key)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.Key2)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.Value)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.Key2Hasher)
	if err != nil {
		return err
	}

	return nil
}

func (m TypDoubleMap) Encode(encoder scale.Encoder) error {
	err := encoder.Encode(m.Hasher)
	if err != nil {
		return err
	}

	err = encoder.Encode(m.Key)
	if err != nil {
		return err
	}

	err = encoder.Encode(m.Key2)
	if err != nil {
		return err
	}

	err = encoder.Encode(m.Value)
	if err != nil {
		return err
	}

	err = encoder.Encode(m.Key2Hasher)
	if err != nil {
		return err
	}

	return nil
}

type FunctionArgumentMetadata struct {
	Name string
	Type string
}

func (m *FunctionArgumentMetadata) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&m.Name)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.Type)
	if err != nil {
		return err
	}

	return nil
}

func (m FunctionArgumentMetadata) Encode(encoder scale.Encoder) error {
	err := encoder.Encode(m.Name)
	if err != nil {
		return err
	}

	err = encoder.Encode(m.Type)
	if err != nil {
		return err
	}

	return nil
}
