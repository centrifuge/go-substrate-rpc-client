package events

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

type Event struct {
	Name   string
	Fields map[string]any
	Phase  *types.Phase
	Topics []types.Hash
}

func ParseEvents(meta *types.Metadata, sd *types.StorageDataRaw) ([]*Event, error) {
	decoder := scale.NewDecoder(bytes.NewReader(*sd))

	// determine number of events
	n, err := decoder.DecodeUintCompact()
	if err != nil {
		return nil, fmt.Errorf("couldn't get number of events: %w", err)
	}

	var events []*Event

	for i := uint64(0); i < n.Uint64(); i++ {
		// decode Phase
		var phase types.Phase

		if err := decoder.Decode(&phase); err != nil {
			return nil, fmt.Errorf("couldn't decode Phase for event #%v: %w", i, err)
		}

		// decode EventID
		var eventID types.EventID

		if err := decoder.Decode(&eventID); err != nil {
			return nil, fmt.Errorf("couldn't decode EventID for event #%v: %w", i, err)
		}

		// ask metadata for method & event name for event
		event, err := parseEvent(meta, decoder, eventID)
		if err != nil {
			return nil, fmt.Errorf("couldn't parse event #%v with EventID %v: %w", i, eventID, err)
		}

		var topics []types.Hash

		if err := decoder.Decode(&topics); err != nil {
			return nil, fmt.Errorf("unable to decode topics for event #%v: %w", i, err)
		}

		event.Phase = &phase
		event.Topics = topics

		events = append(events, event)
	}

	return events, nil
}

func parseEvent(meta *types.Metadata, decoder *scale.Decoder, eventID types.EventID) (*Event, error) {
	for _, mod := range meta.AsMetadataV14.Pallets {
		if !mod.HasEvents {
			continue
		}

		if mod.Index != types.NewU8(eventID[0]) {
			continue
		}

		eventType := mod.Events.Type.Int64()

		typ, ok := meta.AsMetadataV14.EfficientLookup[eventType]

		if !ok {
			return nil, fmt.Errorf("event with type %d not found", eventType)
		}

		if !typ.Def.IsVariant {
			return nil, fmt.Errorf("event with type %d is not a variant", eventType)
		}

		variants := typ.Def.Variant.Variants

		for _, variant := range variants {
			if uint8(variant.Index) != eventID[1] {
				continue
			}

			eventName := fmt.Sprintf("%s.%s", mod.Name, variant.Name)

			if len(variant.Fields) == 0 {
				return &Event{eventName, nil, nil, nil}, nil
			}

			eventFields, err := parseFields(meta, decoder, variant.Fields)

			if err != nil {
				return nil, fmt.Errorf("couldn't parse fields for event '%s': %w", eventName, err)
			}

			return &Event{eventName, eventFields, nil, nil}, nil
		}

		return nil, fmt.Errorf("event with index %d not found", eventID[1])
	}

	return nil, fmt.Errorf("module index %v out of range", eventID[0])
}

func parseFields(meta *types.Metadata, decoder *scale.Decoder, fields []types.Si1Field) (map[string]any, error) {
	eventFields := make(map[string]any)

	for _, field := range fields {
		fieldType, ok := meta.AsMetadataV14.EfficientLookup[field.Type.Int64()]

		if !ok {
			return nil, fmt.Errorf("type not found for field %s", field.Name)
		}

		fieldName := getFieldName(field)

		fieldTypeDef := fieldType.Def

		fieldValue, err := decodeTypeDef(meta, decoder, fieldTypeDef)

		if err != nil {
			return nil, fmt.Errorf("couldn't decode type definition for field '%s': %w", fieldName, err)
		}

		eventFields[fieldName] = fieldValue
	}

	return eventFields, nil
}

func decodeTypeDef(meta *types.Metadata, decoder *scale.Decoder, typeDef types.Si1TypeDef) (any, error) {
	switch {
	case typeDef.IsCompact:
		compactFieldType, ok := meta.AsMetadataV14.EfficientLookup[typeDef.Compact.Type.Int64()]

		if !ok {
			return nil, errors.New("type not found for compact field")
		}

		return decodeCompact(meta, decoder, compactFieldType.Def)
	case typeDef.IsComposite:
		compositeFields, err := parseFields(meta, decoder, typeDef.Composite.Fields)

		if err != nil {
			return nil, fmt.Errorf("couldn't parse composite fields: %w", err)
		}

		return compositeFields, nil
	case typeDef.IsVariant:
		return decodeVariant(meta, decoder, typeDef)
	case typeDef.IsPrimitive:
		primitiveValue, err := decodePrimitive(decoder, typeDef.Primitive.Si0TypeDefPrimitive)

		if err != nil {
			return nil, fmt.Errorf("couldn't decode primitive type: %w", err)
		}

		return primitiveValue, nil
	case typeDef.IsArray:
		arrayFieldType, ok := meta.AsMetadataV14.EfficientLookup[typeDef.Array.Type.Int64()]

		if !ok {
			return nil, fmt.Errorf("type not found for array field")
		}

		// TODO(cdamian): Cover all types.
		switch {
		case arrayFieldType.Def.IsPrimitive:
			arrayValue, err := decodePrimitiveArrayOfLength(uint32(typeDef.Array.Len), decoder, arrayFieldType.Def.Primitive.Si0TypeDefPrimitive)

			if err != nil {
				return nil, fmt.Errorf("couldn't get primitive slice for array field: %w", err)
			}

			return arrayValue, nil
		default:
			return nil, errors.New("unsupported array field type definition")
		}
	case typeDef.IsSequence:
		vectorFieldType, ok := meta.AsMetadataV14.EfficientLookup[typeDef.Sequence.Type.Int64()]

		if !ok {
			return nil, errors.New("type not found for vector field")
		}

		// TODO(cdamian): Cover all types.
		switch {
		case vectorFieldType.Def.IsPrimitive:
			vectorValue, err := decodePrimitiveVector(decoder, vectorFieldType.Def.Primitive.Si0TypeDefPrimitive)

			if err != nil {
				return nil, fmt.Errorf("couldn't get primitive slice for vector field: %w", err)
			}

			return vectorValue, nil
		case vectorFieldType.Def.IsTuple:
			if vectorFieldType.Def.Tuple == nil {
				return nil, nil
			}

			vectorLen, err := decoder.DecodeUintCompact()

			if err != nil {
				return nil, fmt.Errorf("couldn't decode vector length: %w", err)
			}

			var tupleVector []any

			for i := uint64(0); i < vectorLen.Uint64(); i++ {
				tuple, err := decodeTuple(meta, decoder, vectorFieldType.Def.Tuple)

				if err != nil {
					return nil, fmt.Errorf("couldn't decode vector tuple: %w", err)
				}

				tupleVector = append(tupleVector, tuple)
			}

			return tupleVector, nil
		case vectorFieldType.Def.IsComposite:
			vectorLen, err := decoder.DecodeUintCompact()

			if err != nil {
				return nil, fmt.Errorf("couldn't decode vector length: %w", err)
			}

			var compositeVector []any

			for i := uint64(0); i < vectorLen.Uint64(); i++ {
				compositeFields, err := parseFields(meta, decoder, vectorFieldType.Def.Composite.Fields)

				if err != nil {
					return nil, fmt.Errorf("couldn't parse composite vector fields: %w", err)
				}

				compositeVector = append(compositeVector, compositeFields)
			}

			return compositeVector, nil
		case vectorFieldType.Def.IsVariant:
			vectorLen, err := decoder.DecodeUintCompact()

			if err != nil {
				return nil, fmt.Errorf("couldn't decode vector length: %w", err)
			}

			var variantVector []any

			for i := uint64(0); i < vectorLen.Uint64(); i++ {
				variant, err := decodeVariant(meta, decoder, vectorFieldType.Def)

				if err != nil {
					return nil, fmt.Errorf("couldn't parse composite vector fields: %w", err)
				}

				variantVector = append(variantVector, variant)
			}

			return variantVector, nil
		default:
			return nil, errors.New("unsupported vector field type definition")
		}
	case typeDef.IsTuple:
		if typeDef.Tuple == nil {
			return nil, nil
		}

		tuple, err := decodeTuple(meta, decoder, typeDef.Tuple)

		if err != nil {
			return nil, fmt.Errorf("couldn't decode touple: %w", err)
		}

		return tuple, nil
	default:
		return nil, errors.New("unsupported field type definition")
	}
}

func decodeVariant(meta *types.Metadata, decoder *scale.Decoder, typeDef types.Si1TypeDef) (any, error) {
	variantByte, err := decoder.ReadOneByte()

	if err != nil {
		return nil, fmt.Errorf("couldn't read variant byte: %w", err)
	}

	for _, variant := range typeDef.Variant.Variants {
		if byte(variant.Index) != variantByte {
			continue
		}

		if len(variant.Fields) == 0 {
			return string(variant.Name), nil
		}

		fieldMap := make(map[string]any)

		variantFields, err := parseFields(meta, decoder, variant.Fields)

		if err != nil {
			return nil, fmt.Errorf("couldn't parse variant fields: %w", err)
		}

		fieldMap[string(variant.Name)] = variantFields

		return fieldMap, nil
	}

	return nil, fmt.Errorf("variant %d not found", variantByte)
}

func decodeCompact(meta *types.Metadata, decoder *scale.Decoder, typeDef types.Si1TypeDef) (any, error) {
	// TODO(cdamian): Cover all types.
	switch {
	case typeDef.IsPrimitive:
		// TODO(cdamian): Confirm that this covers all primitive types.
		fieldValue, err := decode[types.UCompact](decoder)

		if err != nil {
			return nil, fmt.Errorf("couldn't decode primitive type: %w", err)
		}

		return fieldValue, nil
	case typeDef.IsComposite:
		compactCompositeFields := typeDef.Composite.Fields

		compactCompositeMap := make(map[string]any)

		for _, compactCompositeField := range compactCompositeFields {
			compactCompositeFieldType, ok := meta.AsMetadataV14.EfficientLookup[compactCompositeField.Type.Int64()]

			if !ok {
				return nil, errors.New("compact composite field type not found")
			}

			compactComposite, err := decodeCompact(meta, decoder, compactCompositeFieldType.Def)

			if err != nil {
				return nil, fmt.Errorf("couldn't decode compact composite type: %w", err)
			}

			compactCompositeMap[string(compactCompositeField.Name)] = compactComposite
		}

		return compactCompositeMap, nil
	default:
		return nil, errors.New("unsupported compact field type")
	}
}

func decodeTuple(meta *types.Metadata, decoder *scale.Decoder, tuple types.Si1TypeDefTuple) (any, error) {
	var res []any

	for _, item := range tuple {
		// TODO(cdamian): Add a tuple item struct that has 2 fields for name (itemTypeDef.Path) and value?
		itemTypeDef, ok := meta.AsMetadataV14.EfficientLookup[item.Int64()]

		if !ok {
			return nil, fmt.Errorf("type definition for tuple item %d not found", item.Int64())
		}

		tupleItem, err := decodeTypeDef(meta, decoder, itemTypeDef.Def)

		if err != nil {
			return nil, fmt.Errorf("couldn't decode tuple item: %w", err)
		}

		res = append(res, tupleItem)
	}

	return res, nil
}

func decodePrimitiveVector(decoder *scale.Decoder, primitiveTypeDef types.Si0TypeDefPrimitive) (any, error) {
	switch primitiveTypeDef {
	case types.IsBool:
		return decodeVector[bool](decoder)
	case types.IsChar:
		return decodeVector[byte](decoder)
	case types.IsStr:
		return decodeVector[string](decoder)
	case types.IsU8:
		return decodeVector[types.U8](decoder)
	case types.IsU16:
		return decodeVector[types.U16](decoder)
	case types.IsU32:
		return decodeVector[types.U32](decoder)
	case types.IsU64:
		return decodeVector[types.U64](decoder)
	case types.IsU128:
		return decodeVector[types.U128](decoder)
	case types.IsU256:
		return decodeVector[types.U256](decoder)
	case types.IsI8:
		return decodeVector[types.I8](decoder)
	case types.IsI16:
		return decodeVector[types.I16](decoder)
	case types.IsI32:
		return decodeVector[types.I32](decoder)
	case types.IsI64:
		return decodeVector[types.I64](decoder)
	case types.IsI128:
		return decodeVector[types.I128](decoder)
	case types.IsI256:
		return decodeVector[types.I256](decoder)
	}

	return nil, fmt.Errorf("unsupported primitive type %v", primitiveTypeDef)
}

func decodeVector[T any](decoder *scale.Decoder) ([]T, error) {
	var slice []T

	if err := decoder.Decode(&slice); err != nil {
		return nil, err
	}

	return slice, nil
}

func decodePrimitiveArrayOfLength(length uint32, decoder *scale.Decoder, primitiveTypeDef types.Si0TypeDefPrimitive) (any, error) {
	switch primitiveTypeDef {
	case types.IsBool:
		return decodeArrayOfLength[bool](decoder, length)
	case types.IsChar:
		return decodeArrayOfLength[byte](decoder, length)
	case types.IsStr:
		return decodeArrayOfLength[string](decoder, length)
	case types.IsU8:
		return decodeArrayOfLength[types.U8](decoder, length)
	case types.IsU16:
		return decodeArrayOfLength[types.U16](decoder, length)
	case types.IsU32:
		return decodeArrayOfLength[types.U32](decoder, length)
	case types.IsU64:
		return decodeArrayOfLength[types.U64](decoder, length)
	case types.IsU128:
		return decodeArrayOfLength[types.U128](decoder, length)
	case types.IsU256:
		return decodeArrayOfLength[types.U256](decoder, length)
	case types.IsI8:
		return decodeArrayOfLength[types.I8](decoder, length)
	case types.IsI16:
		return decodeArrayOfLength[types.I16](decoder, length)
	case types.IsI32:
		return decodeArrayOfLength[types.I32](decoder, length)
	case types.IsI64:
		return decodeArrayOfLength[types.I64](decoder, length)
	case types.IsI128:
		return decodeArrayOfLength[types.I128](decoder, length)
	case types.IsI256:
		return decodeArrayOfLength[types.I256](decoder, length)
	}

	return nil, fmt.Errorf("unsupported primitive type %v", primitiveTypeDef)
}

func decodeArrayOfLength[T any](decoder *scale.Decoder, length uint32) ([]T, error) {
	slice := make([]T, 0, length)

	for i := uint32(0); i < length; i++ {
		var t T

		if err := decoder.Decode(&t); err != nil {
			return nil, err
		}

		slice = append(slice, t)
	}

	return slice, nil
}

func decodePrimitive(decoder *scale.Decoder, primitiveTypeDef types.Si0TypeDefPrimitive) (any, error) {
	switch primitiveTypeDef {
	case types.IsBool:
		return decode[bool](decoder)
	case types.IsChar:
		return decode[byte](decoder)
	case types.IsStr:
		return decode[string](decoder)
	case types.IsU8:
		return decode[types.U8](decoder)
	case types.IsU16:
		return decode[types.U16](decoder)
	case types.IsU32:
		return decode[types.U32](decoder)
	case types.IsU64:
		return decode[types.U64](decoder)
	case types.IsU128:
		return decode[types.U128](decoder)
	case types.IsU256:
		return decode[types.U256](decoder)
	case types.IsI8:
		return decode[types.I8](decoder)
	case types.IsI16:
		return decode[types.I16](decoder)
	case types.IsI32:
		return decode[types.I32](decoder)
	case types.IsI64:
		return decode[types.I64](decoder)
	case types.IsI128:
		return decode[types.I128](decoder)
	case types.IsI256:
		return decode[types.I256](decoder)
	}

	return nil, fmt.Errorf("unsupported primitive type %v", primitiveTypeDef)
}

func decode[T any](decoder *scale.Decoder) (T, error) {
	var t T

	if err := decoder.Decode(&t); err != nil {
		return t, err
	}

	return t, nil
}

func getFieldName(field types.Si1Field) string {
	switch {
	case field.HasName:
		return string(field.Name)
	case field.HasTypeName:
		return string(field.TypeName)
	default:
		return "unknown_field_name"
	}
}
