// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
// Copyright (C) 2019  Philip Stanislaus, Philip Stehlik, Vimukthi Wickramasinghe
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
	"strings"

	"github.com/centrifuge/go-substrate-rpc-client/scale"
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

func (m *RuntimeMetadataV4) MethodIndex(method string) MethodIDX {
	s := strings.Split(method, ".")
	var sIDX, mIDX uint8 = 0, 0
	// section index
	var sCounter = 0

	for _, n := range m.Modules {
		if n.HasCalls {
			if n.Name == s[0] {
				sIDX = uint8(sCounter)
				for j, f := range n.Calls {
					if f.Name == s[1] {
						mIDX = uint8(j)
					}
				}
			}
			sCounter++
		}
	}

	return MethodIDX{sIDX, mIDX}
}

// MethodIDX [sectionIndex, methodIndex] 16bits
type MethodIDX struct {
	SectionIndex uint8
	MethodIndex  uint8
}

func (m *MethodIDX) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&m.SectionIndex)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.MethodIndex)
	if err != nil {
		return err
	}

	return nil
}

func (m MethodIDX) Encode(encoder scale.Encoder) error {
	err := encoder.Encode(m.SectionIndex)
	if err != nil {
		return err
	}

	err = encoder.Encode(m.MethodIndex)
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
	Plane         string
	Map           TypMap
	DMap          TypDoubleMap
	Fallback      []byte
	Documentation []string
}

// TODO add again, write test
//func (s StorageFunctionMetadata) isMap() bool {
//	return s.Type == 1
//}

// TODO add again, write test
//func (s StorageFunctionMetadata) isDMap() bool {
//	return s.Type == 2
//}

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

//func (t TypMap) HashFunc() (hash.Hash, error) {
//	if t.Hasher == 1 {
//		return blake2b.New(&blake2b.Config{Size: 32})
//	}
//	return hash.Hash{}, errors.New("hash function type not supported")
//}

func (m *TypMap) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&m.Hasher)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.Key)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.Value)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.IsLinked)
	if err != nil {
		return err
	}

	return nil
}

func (m TypMap) Encode(encoder scale.Encoder) error {
	err := encoder.Encode(m.Hasher)
	if err != nil {
		return err
	}

	err = encoder.Encode(m.Key)
	if err != nil {
		return err
	}

	err = encoder.Encode(m.Value)
	if err != nil {
		return err
	}

	err = encoder.Encode(m.IsLinked)
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

//type StorageKey []byte

//func NewStorageKey(meta Metadata, module string, fn string, key []byte) (StorageKey, error) {
//	var fnMeta *StorageFunctionMetadata
//	for _, m := range meta.Metadata.Modules {
//		if m.Prefix == module {
//			for _, s := range m.Storage {
//				if s.Name == fn {
//					fnMeta = &s
//					break
//				}
//			}
//		}
//	}
//	if fnMeta == nil {
//		return nil, fmt.Errorf("no meta data found for module %s function %s", module, fn)
//	}
//
//	var hasher hash.Hash
//	var err error
//	if fnMeta.isMap() {
//		hasher, err = fnMeta.Map.HashFunc()
//		if err != nil {
//			return nil, err
//		}
//	} else if fnMeta.isDMap() {
//		// TODO define hashing for 2 keys
//	}
//
//	afn := []byte(module + " " + fn)
//	// TODO why is add length prefix step in JS client doesn't add anything to the hashed key?
//	if hasher != nil {
//		hasher.Write(append(afn, key...))
//		return hasher.Sum(nil), nil
//	} else {
//		if key != nil {
//			return createMultiXxhash(append(afn, key...), 2), nil
//		}
//		return createMultiXxhash(append(afn), 2), nil
//	}
//}

//func (s StorageKey) Encode(encoder scale.Encoder) error {
//	return encoder.Encode(s)
//}
//
//type StorageData []byte
//
//func (s StorageData) Decoder() *scale.Decoder {
//	buf := bytes.NewBuffer(s[:])
//	return scale.NewDecoder(buf)
//}
//
//func (s *State) Storage(key StorageKey, block []byte) (StorageData, error) {
//	var res string
//	var err error
//	if block != nil {
//		err = s.client.Call(&res, "state_getStorage", hexutil.Encode(key), hexutil.Encode(block))
//	} else {
//		err = s.client.Call(&res, "state_getStorage", hexutil.Encode(key))
//	}
//
//	if err != nil {
//		return nil, err
//	}
//
//	if res == "" {
//		return nil, errors.New("empty result")
//	}
//
//	return hexutil.Decode(res)
//}
//
//func createMultiXxhash(data []byte, rounds int) []byte {
//	res := make([]byte, 0)
//	for i := 0; i < rounds; i++ {
//		h := xxHash64.New(uint64(i))
//		h.Write(data)
//		res = append(res, h.Sum(nil)...)
//	}
//	return res
//}
