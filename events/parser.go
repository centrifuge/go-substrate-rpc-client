package events

import (
	"bytes"
	"encoding/json"
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

type OverrideDecodeFn func(decoder *scale.Decoder) (any, error)

var (
	overrideDecodeFnMap = map[string]OverrideDecodeFn{
		//"ref_time": func(decoder *scale.Decoder) (any, error) {
		//	return decode[types.UCompact](decoder)
		//},
	}
)

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
			return nil, fmt.Errorf("unable to decode Phase for event #%v: %v", i, err)
		}

		// decode EventID
		var eventID types.EventID

		if err := decoder.Decode(&eventID); err != nil {
			return nil, fmt.Errorf("unable to decode EventID for event #%v: %v", i, err)
		}

		// ask metadata for method & event name for event
		event, err := parseEvent(meta, decoder, eventID)
		if err != nil {
			return nil, fmt.Errorf("unable to find event with EventID %v in metadata for event #%v: %s", eventID, i, err)
		}

		var topics []types.Hash

		if err := decoder.Decode(&topics); err != nil {
			return nil, fmt.Errorf("unable to decode topics for event #%v: %v", i, err)
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

			fmt.Println("Parsing event", eventName)

			if len(variant.Fields) == 0 {
				return &Event{eventName, nil, nil, nil}, nil
			}

			eventFields, err := parseFields(meta, decoder, variant.Fields)

			if err != nil {
				return nil, fmt.Errorf("couldn't parse event fields: %w", err)
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

		fieldName := string(field.Name)

		if overrideDecodeFn, ok := overrideDecodeFnMap[fieldName]; ok {
			fieldValue, err := overrideDecodeFn(decoder)

			if err != nil {
				return nil, fmt.Errorf("couldn't decode '%s' using override func: %w", fieldName, err)
			}

			eventFields[fieldName] = fieldValue

			jsonPrint(fieldName, fieldValue)

			continue
		}

		fieldTypeDef := fieldType.Def

		switch {
		case fieldTypeDef.IsComposite:
			compositeFields, err := parseFields(meta, decoder, fieldTypeDef.Composite.Fields)

			if err != nil {
				return nil, fmt.Errorf("couldn't parse composite fields for %s: %w", fieldName, err)
			}

			eventFields[fieldName] = compositeFields

			jsonPrint(fieldName, compositeFields)
		case fieldTypeDef.IsVariant:
			variantByte, err := decoder.ReadOneByte()

			if err != nil {
				return nil, fmt.Errorf("couldn't read variant byte for %s: %w", fieldName, err)
			}

			variantFound := false

			for _, variant := range fieldTypeDef.Variant.Variants {
				if byte(variant.Index) != variantByte {
					continue
				}

				variantFound = true

				if len(variant.Fields) == 0 {
					eventFields[fieldName] = string(variant.Name)

					jsonPrint(fieldName, string(variant.Name))

					break
				}

				variantFields, err := parseFields(meta, decoder, variant.Fields)

				if err != nil {
					return nil, fmt.Errorf("couldn't parse variant fields for %s: %w", fieldName, err)
				}

				eventFields[fieldName] = variantFields

				jsonPrint(fieldName, variantFields)
			}

			if !variantFound {
				return nil, fmt.Errorf("variant %d not found for %s", variantByte, fieldName)
			}
		case fieldTypeDef.IsCompact:
			compactFieldType, ok := meta.AsMetadataV14.EfficientLookup[fieldTypeDef.Compact.Type.Int64()]

			if !ok {
				return nil, fmt.Errorf("type not found for compact field %s", field.Name)
			}

			switch {
			case compactFieldType.Def.IsPrimitive:
				fieldValue, err := decode[types.UCompact](decoder)

				if err != nil {
					return nil, fmt.Errorf("couldn't decode primitive type for %s: %w", fieldName, err)
				}

				eventFields[fieldName] = fieldValue

				jsonPrint(fieldName, fieldValue)
			default:
				return nil, fmt.Errorf("unsupported compact field type for %s", fieldName)
			}
		case fieldTypeDef.IsPrimitive:
			primitiveValue, err := decodePrimitive(decoder, fieldTypeDef.Primitive.Si0TypeDefPrimitive)

			if err != nil {
				return nil, fmt.Errorf("couldn't decode primitive type for %s: %w", fieldName, err)
			}

			eventFields[fieldName] = primitiveValue

			jsonPrint(fieldName, primitiveValue)
		case fieldTypeDef.IsArray:
			arrayFieldType, ok := meta.AsMetadataV14.EfficientLookup[fieldTypeDef.Array.Type.Int64()]

			if !ok {
				return nil, fmt.Errorf("type not found for array field '%s'", fieldName)
			}

			switch {
			case arrayFieldType.Def.IsPrimitive:
				arrayValue, err := decodePrimitiveArrayOfLength(uint32(fieldTypeDef.Array.Len), decoder, arrayFieldType.Def.Primitive.Si0TypeDefPrimitive)

				if err != nil {
					return nil, fmt.Errorf("couldn't get primitive slice for array field '%s': %w", fieldName, err)
				}

				eventFields[fieldName] = arrayValue

				jsonPrint(fieldName, arrayValue)
			default:
				return nil, errors.New("unspupported array field type definition")
			}
		default:
			return nil, errors.New("unsupported field type definition")
		}
	}

	return eventFields, nil
}

func jsonPrint(fieldName string, obj any) {
	b, _ := json.Marshal(obj)

	fmt.Printf("Field name - '%s': %s\n", fieldName, string(b))
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
