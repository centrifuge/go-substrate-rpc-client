package types

import (
	"fmt"
	"hash"
	"strings"

	"github.com/centrifuge/go-substrate-rpc-client/v3/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v3/xxhash"
)

type MetadataV13 struct {
	Modules   []ModuleMetadataV13
	Extrinsic ExtrinsicV11
}

func (m *MetadataV13) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&m.Modules)
	if err != nil {
		return err
	}
	return decoder.Decode(&m.Extrinsic)
}

func (m MetadataV13) Encode(encoder scale.Encoder) error {
	err := encoder.Encode(m.Modules)
	if err != nil {
		return err
	}
	return encoder.Encode(m.Extrinsic)
}

func (m *MetadataV13) FindCallIndex(call string) (CallIndex, error) {
	s := strings.Split(call, ".")
	for _, mod := range m.Modules {
		if !mod.HasCalls {
			continue
		}
		if string(mod.Name) != s[0] {
			continue
		}
		for ci, f := range mod.Calls {
			if string(f.Name) == s[1] {
				return CallIndex{mod.Index, uint8(ci)}, nil
			}
		}
		return CallIndex{}, fmt.Errorf("method %v not found within module %v for call %v", s[1], mod.Name, call)
	}
	return CallIndex{}, fmt.Errorf("module %v not found in metadata for call %v", s[0], call)
}

func (m *MetadataV13) FindEventNamesForEventID(eventID EventID) (Text, Text, error) {
	for _, mod := range m.Modules {
		if !mod.HasEvents {
			continue
		}
		if mod.Index != eventID[0] {
			continue
		}
		if int(eventID[1]) >= len(mod.Events) {
			return "", "", fmt.Errorf("event index %v for module %v out of range", eventID[1], mod.Name)
		}
		return mod.Name, mod.Events[eventID[1]].Name, nil
	}
	return "", "", fmt.Errorf("module index %v out of range", eventID[0])
}

func (m *MetadataV13) FindStorageEntryMetadata(module string, fn string) (StorageEntryMetadata, error) {
	for _, mod := range m.Modules {
		if !mod.HasStorage {
			continue
		}
		if string(mod.Storage.Prefix) != module {
			continue
		}
		for _, s := range mod.Storage.Items {
			if string(s.Name) != fn {
				continue
			}
			return s, nil
		}
		return nil, fmt.Errorf("storage %v not found within module %v", fn, module)
	}
	return nil, fmt.Errorf("module %v not found in metadata", module)
}

func (m *MetadataV13) FindConstantValue(module Text, constant Text) ([]byte, error) {
	for _, mod := range m.Modules {
		if mod.Name == module {
			value, err := mod.FindConstantValue(constant)
			if err == nil {
				return value, nil
			}
		}
	}
	return nil, fmt.Errorf("could not find constant %s.%s", module, constant)
}

func (m *MetadataV13) ExistsModuleMetadata(module string) bool {
	for _, mod := range m.Modules {
		if string(mod.Name) == module {
			return true
		}
	}
	return false
}

type ModuleMetadataV13 struct {
	Name       Text
	HasStorage bool
	Storage    StorageMetadataV13
	HasCalls   bool
	Calls      []FunctionMetadataV4
	HasEvents  bool
	Events     []EventMetadataV4
	Constants  []ModuleConstantMetadataV6
	Errors     []ErrorMetadataV8
	Index      uint8
}

func (m *ModuleMetadataV13) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&m.Name)
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

	err = decoder.Decode(&m.Constants)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.Errors)
	if err != nil {
		return err
	}

	return decoder.Decode(&m.Index)
}

func (m ModuleMetadataV13) Encode(encoder scale.Encoder) error {
	err := encoder.Encode(m.Name)
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

	err = encoder.Encode(m.Constants)
	if err != nil {
		return err
	}

	err = encoder.Encode(m.Errors)
	if err != nil {
		return err
	}

	return encoder.Encode(m.Index)
}

func (m *ModuleMetadataV13) FindConstantValue(constant Text) ([]byte, error) {
	for _, cons := range m.Constants {
		if cons.Name == constant {
			return cons.Value, nil
		}
	}
	return nil, fmt.Errorf("could not find constant %s", constant)
}

type StorageMetadataV13 struct {
	Prefix Text
	Items  []StorageFunctionMetadataV13
}

type StorageFunctionMetadataV13 struct {
	Name          Text
	Modifier      StorageFunctionModifierV0
	Type          StorageFunctionTypeV13
	Fallback      Bytes
	Documentation []Text
}

func (s StorageFunctionMetadataV13) IsPlain() bool {
	return s.Type.IsType
}

func (s StorageFunctionMetadataV13) IsMap() bool {
	return s.Type.IsMap
}

func (s StorageFunctionMetadataV13) IsDoubleMap() bool {
	return s.Type.IsDoubleMap
}

func (s StorageFunctionMetadataV13) IsNMap() bool {
	return s.Type.IsNMap
}

func (s StorageFunctionMetadataV13) Hasher() (hash.Hash, error) {
	if s.Type.IsMap {
		return s.Type.AsMap.Hasher.HashFunc()
	}
	if s.Type.IsDoubleMap {
		return s.Type.AsDoubleMap.Hasher.HashFunc()
	}
	if s.Type.IsNMap {
		return nil, fmt.Errorf("only Map and DoubleMap have a Hasher")
	}
	return xxhash.New128(nil), nil
}

func (s StorageFunctionMetadataV13) Hasher2() (hash.Hash, error) {
	if !s.Type.IsDoubleMap {
		return nil, fmt.Errorf("only DoubleMaps have a Hasher2")
	}
	return s.Type.AsDoubleMap.Key2Hasher.HashFunc()
}

func (s StorageFunctionMetadataV13) Hashers() ([]hash.Hash, error) {
	if !s.Type.IsNMap {
		return nil, fmt.Errorf("only NMaps have Hashers")
	}

	hashers := make([]hash.Hash, len(s.Type.AsNMap.Hashers))
	for i, hasher := range s.Type.AsNMap.Hashers {
		hasherFn, err := hasher.HashFunc()
		if err != nil {
			return nil, err
		}
		hashers[i] = hasherFn
	}
	return hashers, nil
}

type StorageFunctionTypeV13 struct {
	IsType      bool
	AsType      Type // 0
	IsMap       bool
	AsMap       MapTypeV10 // 1
	IsDoubleMap bool
	AsDoubleMap DoubleMapTypeV10 // 2
	IsNMap      bool
	AsNMap      NMapTypeV13 // 3
}

func (s *StorageFunctionTypeV13) Decode(decoder scale.Decoder) error {
	var t uint8
	err := decoder.Decode(&t)
	if err != nil {
		return err
	}

	switch t {
	case 0:
		s.IsType = true
		err = decoder.Decode(&s.AsType)
		if err != nil {
			return err
		}
	case 1:
		s.IsMap = true
		err = decoder.Decode(&s.AsMap)
		if err != nil {
			return err
		}
	case 2:
		s.IsDoubleMap = true
		err = decoder.Decode(&s.AsDoubleMap)
		if err != nil {
			return err
		}
	case 3:
		s.IsNMap = true
		err = decoder.Decode(&s.AsNMap)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("received unexpected type %v", t)
	}
	return nil
}

func (s StorageFunctionTypeV13) Encode(encoder scale.Encoder) error {
	switch {
	case s.IsType:
		err := encoder.PushByte(0)
		if err != nil {
			return err
		}
		err = encoder.Encode(s.AsType)
		if err != nil {
			return err
		}
	case s.IsMap:
		err := encoder.PushByte(1)
		if err != nil {
			return err
		}
		err = encoder.Encode(s.AsMap)
		if err != nil {
			return err
		}
	case s.IsDoubleMap:
		err := encoder.PushByte(2)
		if err != nil {
			return err
		}
		err = encoder.Encode(s.AsDoubleMap)
		if err != nil {
			return err
		}
	case s.IsNMap:
		err := encoder.PushByte(3)
		if err != nil {
			return err
		}
		err = encoder.Encode(s.AsNMap)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("expected to be either type, map, double map or nmap but none was set: %v", s)
	}
	return nil
}

type NMapTypeV13 struct {
	Keys    []Type
	Hashers []StorageHasherV10
	Value   Type
}
