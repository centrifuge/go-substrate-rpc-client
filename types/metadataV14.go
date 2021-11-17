package types

import (
	"fmt"
	"hash"
	"strings"

	"github.com/centrifuge/go-substrate-rpc-client/v3/scale"
)

// nolint:lll
// Based on https://github.com/polkadot-js/api/blob/80b581f0df87108c59f71e67d7c5fc5f8c89ec33/packages/types/src/interfaces/metadata/v14.ts
type MetadataV14 struct {
	Lookup    PortableRegistryV14
	Pallets   []PalletMetadataV14
	Extrinsic ExtrinsicV14
	Type      Si1LookupTypeID

	// Custom field to help us lookup a type from the registry
	// more efficiently. This field is built while decoding and
	// it is not to be encoded.
	EfficientLookup map[int64]*Si1Type `scale:"-"`
}

// Decode implementation for MetadataV14
// Note: We opt for a custom impl build `EfficientLookup`
// on the fly.
func (m *MetadataV14) Decode(decoder scale.Decoder) error {
	var err error
	err = decoder.Decode(&m.Lookup)
	if err != nil {
		return err
	}

	m.EfficientLookup = m.Lookup.toMap()

	err = decoder.Decode(&m.Pallets)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.Extrinsic)
	if err != nil {
		return err
	}

	return decoder.Decode(&m.Type)
}

// Build a map of type id to pointer to the PortableTypeV14 itself.
func (lookup *PortableRegistryV14) toMap() map[int64]*Si1Type {
	var efficientLookup = make(map[int64]*Si1Type)
	var t PortableTypeV14
	for _, t = range lookup.Types {
		// We need to copy t so that the pointer doesn't get
		// overwritten by the next assignment.
		typ := t
		efficientLookup[typ.ID.Int64()] = &typ.Type
	}
	return efficientLookup
}

/* Metadata interface functions implementation */

func (m *MetadataV14) FindCallIndex(call string) (CallIndex, error) {
	s := strings.Split(call, ".")
	for _, mod := range m.Pallets {
		if !mod.HasCalls {
			continue
		}
		if string(mod.Name) != s[0] {
			continue
		}
		callType := mod.Calls.Type.Int64()

		if typ, ok := m.EfficientLookup[callType]; ok {
			if len(typ.Def.Variant.Variants) > 0 {
				for _, vars := range typ.Def.Variant.Variants {
					if string(vars.Name) == s[1] {
						return CallIndex{uint8(mod.Index), uint8(vars.Index)}, nil
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
		if mod.Index != NewU8(eventID[0]) {
			continue
		}
		eventType := mod.Events.Type.Int64()

		if typ, ok := m.EfficientLookup[eventType]; ok {
			if len(typ.Def.Variant.Variants) > 0 {
				for _, vars := range typ.Def.Variant.Variants {
					if uint8(vars.Index) == eventID[1] {
						return mod.Name, vars.Name, nil
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
			if string(s.Name) == fn {
				return s, nil
			}
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

/* Supporting types */

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

type PortableRegistryV14 struct {
	Types []PortableTypeV14
}

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
	Index      U8
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

	err = decoder.DecodeOption(&m.HasStorage, &m.Storage)
	if err != nil {
		return err
	}

	err = decoder.DecodeOption(&m.HasCalls, &m.Calls)
	if err != nil {
		return err
	}

	err = decoder.DecodeOption(&m.HasEvents, &m.Events)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.Constants)
	if err != nil {
		return err
	}

	err = decoder.DecodeOption(&m.HasErrors, &m.Errors)
	if err != nil {
		return err
	}

	return decoder.Decode(&m.Index)
}

func (m PalletMetadataV14) Encode(encoder scale.Encoder) error {
	err := encoder.Encode(m.Name)
	if err != nil {
		return err
	}

	err = encoder.EncodeOption(m.HasStorage, m.Storage)
	if err != nil {
		return err
	}

	err = encoder.EncodeOption(m.HasCalls, m.Calls)
	if err != nil {
		return err
	}

	err = encoder.EncodeOption(m.HasEvents, m.Events)
	if err != nil {
		return err
	}

	err = encoder.Encode(m.Constants)
	if err != nil {
		return err
	}

	err = encoder.EncodeOption(m.HasErrors, m.Errors)
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

type StorageEntryMetadataV14 struct {
	Name          Text
	Modifier      StorageFunctionModifierV0
	Type          StorageEntryTypeV14
	Fallback      Bytes
	Documentation []Text
}

type MapTypeV14 struct {
	Hashers []StorageHasherV10
	Key     Si1LookupTypeID
	Value   Si1LookupTypeID
}

func (s StorageEntryMetadataV14) IsPlain() bool {
	return s.Type.IsPlainType
}

func (s StorageEntryMetadataV14) IsMap() bool {
	return s.Type.IsMap
}

func (s StorageEntryMetadataV14) IsDoubleMap() bool {
	panic(unsupportedMapVariantCheck("IsDoubleMap"))
}
func (s StorageEntryMetadataV14) IsNMap() bool {
	panic(unsupportedMapVariantCheck("IsNMap"))
}

func (s StorageEntryMetadataV14) Hasher() (hash.Hash, error) {
	return nil, fmt.Errorf("StorageEntryMetadataV14 does not implement Hasher()")
}

func (s StorageEntryMetadataV14) Hasher2() (hash.Hash, error) {
	return nil, fmt.Errorf("StorageEntryMetadataV14 does not implement Hasher2()")
}

func (s StorageEntryMetadataV14) Hashers() ([]hash.Hash, error) {
	if !s.IsMap() {
		return nil, fmt.Errorf("StorageEntryMetadataV14.Hashers() should be called on a Map entry")
	}

	var hashes []hash.Hash
	for _, hasher := range s.Type.AsMap.Hashers {
		h, err := hasher.HashFunc()
		if err != nil {
			return nil, err
		}
		hashes = append(hashes, h)
	}

	return hashes, nil
}

type StorageEntryTypeV14 struct {
	IsPlainType bool
	AsPlainType Si1LookupTypeID
	IsMap       bool
	AsMap       MapTypeV14
}

func (s *StorageEntryTypeV14) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}
	switch b {
	case 0:
		s.IsPlainType = true
		err = decoder.Decode(&s.AsPlainType)
		if err != nil {
			return err
		}
	case 1:
		s.IsMap = true
		err = decoder.Decode(&s.AsMap)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("StorageFunctionTypeV14 does not support this type: %d", b)
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
		return fmt.Errorf("expected to be either plain type or map, but none was set: %v", s)
	}
	return nil
}

func unsupportedMapVariantCheck(variant string) error {
	return fmt.Errorf("StorageEntryMetadataV14 does not implement %s "+
		"as now there is only one Map variant with n keys",
		variant,
	)
}
