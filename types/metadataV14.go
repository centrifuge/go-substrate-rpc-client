package types

import (
	"fmt"
	"hash"
	"strings"

	"github.com/centrifuge/go-substrate-rpc-client/v3/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v3/xxhash"
)

// Based on:
// https://github.com/polkadot-js/api/blob/48ef04b8ca21dc4bd06442775d9b7585c75d1253/packages/types/src/interfaces/metadata/v14.ts#L30-L34
type MetadataV14 struct {
	Lookup    PortableRegistry
	Pallets   []PalletMetadataV14
	Extrinsic ExtrinsicV14

	LookUpData map[int64]*Si1Type
}

type ExtrinsicV14 struct {
	Type             Si1LookupTypeID
	Version          U8
	SignedExtensions []SignedExtensionMetadataV14
}

type SignedExtensionMetadataV14 struct {
	Identifier       Text
	Type             Si1LookupTypeID
	AdditionalSigned Si1LookupTypeID
}

func (m *MetadataV14) Decode(decoder scale.Decoder) error {
	var err error
	err = decoder.Decode(&m.Lookup)
	if err != nil {
		return err
	}
	fmt.Println("Decoded types")

	m.LookUpData = make(map[int64]*Si1Type)
	for _, lookUp := range m.Lookup {
		m.LookUpData[lookUp.ID.Int64()] = &lookUp.Type
	}
	fmt.Println("Built LookUpData")

	fmt.Println("Will Decode Pallets")

	err = decoder.Decode(&m.Pallets)
	if err != nil {
		return err
	}
	fmt.Println("Decoded Pallets")

	err = decoder.Decode(&m.Extrinsic)
	if err != nil {
		return err
	}

	fmt.Println("Decoded Extrinsic")

	return nil
}

func (m MetadataV14) Encode(encoder scale.Encoder) error {
	err := encoder.Encode(m.Pallets)
	if err != nil {
		return err
	}
	return encoder.Encode(m.Extrinsic)
}

func (m *MetadataV14) FindCallIndex(call string) (CallIndex, error) {
	s := strings.Split(call, ".")
	for _, mod := range m.Pallets {
		if !mod.HasCalls {
			continue
		}
		if string(mod.Name) != s[0] {
			continue
		}
		callType := mod.Calls.Type

		for _, lookUp := range m.Lookup {
			if lookUp.ID.Int64() == callType.Int64() {
				if len(lookUp.Type.Def.Variant.Variants) > 0 {
					for _, vars := range lookUp.Type.Def.Variant.Variants {
						if string(vars.Name) == s[1] {
							return CallIndex{uint8(mod.Index), uint8(vars.Index)}, nil
						}
					}
				}
			}
		}
	}
	return CallIndex{}, fmt.Errorf("module %v not found in metadata for call %v", s[0], call)
}

func (m *MetadataV14) FindEventNamesForEventID(eventID EventID) (Text, Text, error) {
	for _, mod := range m.Pallets {
		if !mod.HasEvents {
			continue
		}
		if mod.Index != eventID[0] {
			continue
		}
		eventType := mod.Events.Type.Int64()

		for _, lookUp := range m.Lookup {
			if lookUp.ID.Int64() == eventType {
				if len(lookUp.Type.Def.Variant.Variants) > 0 {
					for _, vars := range lookUp.Type.Def.Variant.Variants {
						if uint8(vars.Index) == eventID[1] {
							return mod.Name, vars.Name, nil
						}
					}
				}
			}
		}
	}
	return "", "", fmt.Errorf("module index %v out of range", eventID[0])
}

func (m *MetadataV14) FindStorageEntryMetadata(module string, fn string) (StorageEntryMetadata, error) {
	for _, mod := range m.Pallets {
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

func (m *MetadataV14) FindConstantValue(module Text, constant Text) ([]byte, error) {
	for _, mod := range m.Pallets {
		if mod.Name == module {
			value, err := mod.FindConstantValue(constant)
			if err == nil {
				return value, nil
			}
		}
	}
	return nil, fmt.Errorf("could not find constant %s.%s", module, constant)
}

func (m *MetadataV14) ExistsModuleMetadata(module string) bool {
	for _, mod := range m.Pallets {
		if string(mod.Name) == module {
			return true
		}
	}
	return false
}

type PortableRegistry []PortableTypeV14

type PalletMetadataV14 struct {
	Name       Text
	HasStorage bool
	Storage    StorageMetadataV14
	HasCalls   bool
	Calls      FunctionMetadataV14
	HasEvents  bool
	Events     EventMetadataV14
	Constants  []ConstantMetadataV14
	HasErrors  bool
	Errors     ErrorMetadataV14
	Index      uint8
}

type FunctionMetadataV14 struct {
	Type Si1LookupTypeID
}

type EventMetadataV14 struct {
	Type Si1LookupTypeID
}

type ConstantMetadataV14 struct {
	Name  Text
	Type  Si1LookupTypeID
	Value Bytes
	Docs  []Text
}

type ErrorMetadataV14 struct {
	Type Si1LookupTypeID
}

func (m *PalletMetadataV14) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&m.Name)
	if err != nil {
		return err
	}
	fmt.Println("Decoded Pallet.Name")

	err = decoder.Decode(&m.HasStorage)
	if err != nil {
		return err
	}

	fmt.Println("Decoded Pallet.HasStorage")

	if m.HasStorage {
		err = decoder.Decode(&m.Storage)
		if err != nil {
			return err
		}
	}

	fmt.Println("Decoded Pallet.Storage")

	err = decoder.Decode(&m.HasCalls)
	if err != nil {
		return err
	}

	fmt.Println("Decoded Pallet.HasCalls")

	if m.HasCalls {
		err = decoder.Decode(&m.Calls)
		if err != nil {
			return err
		}
	}

	fmt.Println("Decoded Pallet.Calls")

	err = decoder.Decode(&m.HasEvents)
	if err != nil {
		return err
	}

	fmt.Println("Decoded Pallet.HasEvents")

	if m.HasEvents {
		err = decoder.Decode(&m.Events)
		if err != nil {
			return err
		}
	}

	fmt.Println("Decoded Pallet.Events")

	err = decoder.Decode(&m.Constants)
	if err != nil {
		return err
	}

	fmt.Println("Decoded Pallet.Constants")

	err = decoder.Decode(&m.HasErrors)
	if err != nil {
		return err
	}

	fmt.Println("Decoded Pallet.HasErrors")

	if m.HasErrors {
		err = decoder.Decode(&m.Errors)
		if err != nil {
			return err
		}
	}

	fmt.Println("Decoded Pallet.Errors")

	return decoder.Decode(&m.Index)
}

func (m PalletMetadataV14) Encode(encoder scale.Encoder) error {
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

func (m *PalletMetadataV14) FindConstantValue(constant Text) ([]byte, error) {
	for _, cons := range m.Constants {
		if cons.Name == constant {
			return cons.Value, nil
		}
	}
	return nil, fmt.Errorf("could not find constant %s", constant)
}

type StorageMetadataV14 struct {
	Prefix Text
	Items  []StorageEntryMetadataV14
}

func (storage *StorageMetadataV14) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&storage.Prefix)
	if err != nil {
		return err
	}
	return decoder.Decode(&storage.Items)
}

type StorageEntryMetadataV14 struct {
	Name          Text
	Modifier      StorageFunctionModifierV0
	Type          StorageEntryTypeV14
	Fallback      Bytes
	Documentation []Text
}

type MapTypeV14 struct {
	Hasher  []StorageHasherV10
	KeysId  Si1LookupTypeID
	ValueId Si1LookupTypeID
}

func (s StorageEntryMetadataV14) IsPlain() bool {
	return s.Type.IsPlainType
}

func (s StorageEntryMetadataV14) IsMap() bool {
	return false
}

func (s StorageEntryMetadataV14) IsDoubleMap() bool {
	return false
}
func (s StorageEntryMetadataV14) IsNMap() bool {
	return s.Type.IsMap
}

func (s StorageEntryMetadataV14) Hasher() (hash.Hash, error) {
	if s.IsPlain() {
		return xxhash.New128(nil), nil
	}

	return xxhash.New128(nil), nil
}

func (s StorageEntryMetadataV14) Hasher2() (hash.Hash, error) {
	panic("Not implemented")
}

func (s StorageEntryMetadataV14) Hashers() ([]hash.Hash, error) {
	var hashes []hash.Hash
	if s.Type.IsMap {
		for _, hasher := range s.Type.AsMap.Hasher {
			h, err := hasher.HashFunc()
			if err != nil {
				return nil, err
			}
			hashes = append(hashes, h)
		}
	}
	return hashes, nil
}

type StorageEntryTypeV14 struct {
	IsPlainType bool
	AsPlainType Si1LookupTypeID
	IsMap       bool
	AsMap       MapTypeV14
}

type DoubleMapTypeV14 struct {
	Hasher     StorageHasherV10
	Key1       Si1LookupTypeID
	Key2       Si1LookupTypeID
	Value      Si1LookupTypeID
	Key2Hasher StorageHasherV10
}

type NMapTypeV14 struct {
	Key     Si1LookupTypeID
	Hashers []StorageHasherV10
	Value   Si1LookupTypeID
}

func (d *StorageEntryTypeV14) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}
	switch b {
	case 0:
		d.IsPlainType = true
		err = decoder.Decode(&d.AsPlainType)
		if err != nil {
			return err
		}
	case 1:
		d.IsMap = true
		err = decoder.Decode(&d.AsMap)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("StorageFunctionTypeV14 is not support this type: %d", b)
	}
	return nil
}

func (s StorageEntryTypeV14) Encode(encoder scale.Encoder) error {
	switch {
	case s.IsPlainType:
		err := encoder.PushByte(0)
		if err != nil {
			return err
		}
		err = encoder.Encode(s.AsPlainType)
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
	default:
		return fmt.Errorf("expected to be either type, map, double map or nmap but none was set: %v", s)
	}
	return nil
}
