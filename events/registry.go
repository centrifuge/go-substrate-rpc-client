package events

import (
	"errors"
	"fmt"
	"strings"

	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

type EventRegistry map[types.EventID]*EventDecoder

type RegistryFactory interface {
	CreateEventRegistry(meta *types.Metadata) (EventRegistry, error)
}

type registryFactory struct {
	fieldStorage          map[int64]FieldDecoder
	recursiveFieldStorage map[int64]*RecursiveDecoder
}

func NewRegistryFactory() RegistryFactory {
	return &registryFactory{
		fieldStorage:          make(map[int64]FieldDecoder),
		recursiveFieldStorage: make(map[int64]*RecursiveDecoder),
	}
}

func (r *registryFactory) CreateEventRegistry(meta *types.Metadata) (EventRegistry, error) {
	defer r.resetStorages()

	eventRegistry := make(map[types.EventID]*EventDecoder)

	for _, mod := range meta.AsMetadataV14.Pallets {
		if !mod.HasEvents {
			continue
		}

		eventsType, ok := meta.AsMetadataV14.EfficientLookup[mod.Events.Type.Int64()]

		if !ok {
			return nil, fmt.Errorf("events type %d not found for module '%s'", mod.Events.Type.Int64(), mod.Name)
		}

		if !eventsType.Def.IsVariant {
			return nil, fmt.Errorf("events type %d for module '%s' is not a variant", mod.Events.Type.Int64(), mod.Name)
		}

		for _, eventVariant := range eventsType.Def.Variant.Variants {
			eventID := types.EventID{byte(mod.Index), byte(eventVariant.Index)}

			eventName := fmt.Sprintf("%s.%s", mod.Name, eventVariant.Name)

			eventFields, err := r.getEventFields(meta, eventVariant.Fields)

			if err != nil {
				return nil, fmt.Errorf("couldn't get fields for event '%s': %w", eventName, err)
			}

			eventRegistry[eventID] = &EventDecoder{
				Name:   eventName,
				Fields: eventFields,
			}
		}
	}

	if err := r.resolveRecursiveTypes(); err != nil {
		return nil, err
	}

	return eventRegistry, nil
}

func (r *registryFactory) resetStorages() {
	r.fieldStorage = make(map[int64]FieldDecoder)
	r.recursiveFieldStorage = make(map[int64]*RecursiveDecoder)
}

func (r *registryFactory) resolveRecursiveTypes() error {
	for recursiveFieldDecoderLookupIndex, recursiveDecoder := range r.recursiveFieldStorage {
		fieldDecoder, ok := r.fieldStorage[recursiveFieldDecoderLookupIndex]

		if !ok {
			return fmt.Errorf("couldn't get field decoder for recursive type %d", recursiveFieldDecoderLookupIndex)
		}

		if _, ok := fieldDecoder.(*RecursiveDecoder); ok {
			return fmt.Errorf("recursive field type %d cannot be resolved with a non-recursive field decoder", recursiveFieldDecoderLookupIndex)
		}

		recursiveDecoder.FieldDecoder = fieldDecoder
	}

	return nil
}

func (r *registryFactory) getEventFields(meta *types.Metadata, fields []types.Si1Field) ([]*Field, error) {
	var eventFields []*Field

	for _, field := range fields {
		fieldType, ok := meta.AsMetadataV14.EfficientLookup[field.Type.Int64()]

		if !ok {
			return nil, fmt.Errorf("type not found for field '%s'", field.Name)
		}

		fieldName := getFieldName(field, fieldType)

		if fieldDecoder, ok := r.getStoredFieldDecoder(fieldName, field.Type.Int64()); ok {
			eventFields = append(eventFields, &Field{
				Name:         fieldName,
				FieldDecoder: fieldDecoder,
			})
			continue
		}

		fieldTypeDef := fieldType.Def

		fieldDecoder, err := r.getFieldDecoder(meta, fieldName, fieldTypeDef)

		if err != nil {
			return nil, fmt.Errorf("couldn't get field decoder for field '%s': %w", fieldName, err)
		}

		r.storeFieldDecoder(field.Type.Int64(), fieldDecoder)

		eventFields = append(eventFields, &Field{
			Name:         fieldName,
			FieldDecoder: fieldDecoder,
		})
	}

	return eventFields, nil
}

func (r *registryFactory) getFieldDecoder(meta *types.Metadata, fieldName string, typeDef types.Si1TypeDef) (FieldDecoder, error) {
	switch {
	case typeDef.IsCompact:
		compactFieldType, ok := meta.AsMetadataV14.EfficientLookup[typeDef.Compact.Type.Int64()]

		if !ok {
			return nil, errors.New("type not found for compact field")
		}

		return r.getCompactFieldDecoder(meta, fieldName, compactFieldType.Def)
	case typeDef.IsComposite:
		compositeDecoder := &CompositeDecoder{}

		compositeFieldDecoders, err := r.getEventFields(meta, typeDef.Composite.Fields)

		if err != nil {
			return nil, fmt.Errorf("couldn't get composite fields: %w", err)
		}

		compositeDecoder.Fields = compositeFieldDecoders

		return compositeDecoder, nil
	case typeDef.IsVariant:
		return r.getVariantFieldDecoder(meta, typeDef)
	case typeDef.IsPrimitive:
		return getPrimitiveDecoder(typeDef.Primitive.Si0TypeDefPrimitive)
	case typeDef.IsArray:
		arrayFieldType, ok := meta.AsMetadataV14.EfficientLookup[typeDef.Array.Type.Int64()]

		if !ok {
			return nil, fmt.Errorf("type not found for array field")
		}

		return r.getArrayFieldDecoder(uint(typeDef.Array.Len), arrayFieldType.Def)
	case typeDef.IsSequence:
		vectorFieldType, ok := meta.AsMetadataV14.EfficientLookup[typeDef.Sequence.Type.Int64()]

		if !ok {
			return nil, errors.New("type not found for slice field")
		}

		return r.getSliceFieldDecoder(meta, fieldName, vectorFieldType.Def)
	case typeDef.IsTuple:
		if typeDef.Tuple == nil {
			return &NoopDecoder{}, nil
		}

		return r.getTupleDecoder(meta, fieldName, typeDef.Tuple)
	default:
		return nil, errors.New("unsupported field type definition")
	}
}

func (r *registryFactory) getVariantFieldDecoder(meta *types.Metadata, typeDef types.Si1TypeDef) (FieldDecoder, error) {
	variantDecoder := &VariantDecoder{}

	fieldDecoderMap := make(map[byte]FieldDecoder)

	for _, variant := range typeDef.Variant.Variants {
		if len(variant.Fields) == 0 {
			fieldDecoderMap[byte(variant.Index)] = &NoopDecoder{}
			continue
		}

		compositeDecoder := &CompositeDecoder{}

		fieldDecoders, err := r.getEventFields(meta, variant.Fields)

		if err != nil {
			return nil, fmt.Errorf("couldn't get field decoders for variant '%d': %w", variant.Index, err)
		}

		compositeDecoder.Fields = fieldDecoders

		fieldDecoderMap[byte(variant.Index)] = compositeDecoder
	}

	variantDecoder.FieldDecoderMap = fieldDecoderMap

	return variantDecoder, nil
}

func (r *registryFactory) getCompactFieldDecoder(meta *types.Metadata, fieldName string, typeDef types.Si1TypeDef) (FieldDecoder, error) {
	// TODO(cdamian): Cover all types.
	switch {
	case typeDef.IsPrimitive:
		// TODO(cdamian): Confirm that this covers all primitive types.
		return &ValueDecoder[types.UCompact]{}, nil
	case typeDef.IsComposite:
		compactCompositeFields := typeDef.Composite.Fields

		compositeDecoder := &CompositeDecoder{}

		for _, compactCompositeField := range compactCompositeFields {
			compactCompositeFieldType, ok := meta.AsMetadataV14.EfficientLookup[compactCompositeField.Type.Int64()]

			if !ok {
				return nil, errors.New("compact composite field type not found")
			}

			compactFieldName := getFieldName(compactCompositeField, compactCompositeFieldType)

			compactCompositeFieldDecoder, err := r.getCompactFieldDecoder(meta, compactFieldName, compactCompositeFieldType.Def)

			if err != nil {
				return nil, fmt.Errorf("couldn't get compact field decoder: %w", err)
			}

			compositeDecoder.Fields = append(compositeDecoder.Fields, &Field{Name: fieldName, FieldDecoder: compactCompositeFieldDecoder})
		}

		return compositeDecoder, nil
	default:
		return nil, errors.New("unsupported compact field type")
	}
}

func (r *registryFactory) getArrayFieldDecoder(arrayLen uint, typeDef types.Si1TypeDef) (FieldDecoder, error) {
	arrayDecoder := &ArrayDecoder{
		Length: arrayLen,
	}

	// TODO(cdamian): Cover all types.
	switch {
	case typeDef.IsPrimitive:
		primitiveDecoder, err := getPrimitiveDecoder(typeDef.Primitive.Si0TypeDefPrimitive)

		if err != nil {
			return nil, fmt.Errorf("couldn't get primitive decoder for array: %w", err)
		}

		arrayDecoder.ItemDecoder = primitiveDecoder

		return arrayDecoder, nil
	default:
		return nil, errors.New("unsupported array field type definition")
	}
}

func (r *registryFactory) getSliceFieldDecoder(meta *types.Metadata, fieldName string, typeDef types.Si1TypeDef) (FieldDecoder, error) {
	sliceDecoder := &SliceDecoder{}

	itemDecoder, err := r.getFieldDecoder(meta, fieldName, typeDef)

	if err != nil {
		return nil, fmt.Errorf("couldn't get item decoder for slice: %w", err)
	}

	sliceDecoder.ItemDecoder = itemDecoder

	return sliceDecoder, nil
}

func (r *registryFactory) getTupleDecoder(meta *types.Metadata, fieldName string, tuple types.Si1TypeDefTuple) (FieldDecoder, error) {
	compositeDecoder := &CompositeDecoder{}

	for _, item := range tuple {
		itemTypeDef, ok := meta.AsMetadataV14.EfficientLookup[item.Int64()]

		if !ok {
			return nil, fmt.Errorf("type definition for tuple item %d not found", item.Int64())
		}

		fieldDecoder, err := r.getFieldDecoder(meta, fieldName, itemTypeDef.Def)

		if err != nil {
			return nil, fmt.Errorf("couldn't get tuple field decoder: %w", err)
		}

		compositeDecoder.Fields = append(compositeDecoder.Fields, &Field{Name: fieldName, FieldDecoder: fieldDecoder})
	}

	return compositeDecoder, nil
}

func getPrimitiveDecoder(primitiveTypeDef types.Si0TypeDefPrimitive) (FieldDecoder, error) {
	switch primitiveTypeDef {
	case types.IsBool:
		return &ValueDecoder[bool]{}, nil
	case types.IsChar:
		return &ValueDecoder[byte]{}, nil
	case types.IsStr:
		return &ValueDecoder[string]{}, nil
	case types.IsU8:
		return &ValueDecoder[types.U8]{}, nil
	case types.IsU16:
		return &ValueDecoder[types.U16]{}, nil
	case types.IsU32:
		return &ValueDecoder[types.U32]{}, nil
	case types.IsU64:
		return &ValueDecoder[types.U64]{}, nil
	case types.IsU128:
		return &ValueDecoder[types.U128]{}, nil
	case types.IsU256:
		return &ValueDecoder[types.U256]{}, nil
	case types.IsI8:
		return &ValueDecoder[types.I8]{}, nil
	case types.IsI16:
		return &ValueDecoder[types.I16]{}, nil
	case types.IsI32:
		return &ValueDecoder[types.I32]{}, nil
	case types.IsI64:
		return &ValueDecoder[types.I64]{}, nil
	case types.IsI128:
		return &ValueDecoder[types.I128]{}, nil
	case types.IsI256:
		return &ValueDecoder[types.I256]{}, nil
	default:
		return nil, fmt.Errorf("unsupported primitive type %v", primitiveTypeDef)
	}
}

func (r *registryFactory) getStoredFieldDecoder(fieldName string, fieldType int64) (FieldDecoder, bool) {
	if ft, ok := r.fieldStorage[fieldType]; ok {
		if rt, ok := ft.(*RecursiveDecoder); ok {
			r.recursiveFieldStorage[fieldType] = rt
		}

		return ft, ok
	}

	// Ensure that a recursive type such as Xcm::TransferReserveAsset does not cause an infinite loop
	// by adding the RecursiveDecoder the first time a field is encountered.
	r.fieldStorage[fieldType] = &RecursiveDecoder{}

	return nil, false
}

func (r *registryFactory) storeFieldDecoder(fieldType int64, fieldDecoder FieldDecoder) {
	r.fieldStorage[fieldType] = fieldDecoder
}

func getFieldPath(fieldType *types.Si1Type) string {
	var nameParts []string

	for _, pathEntry := range fieldType.Path {
		nameParts = append(nameParts, string(pathEntry))
	}

	return strings.Join(nameParts, "::")
}

func getFieldName(field types.Si1Field, fieldType *types.Si1Type) string {
	if fieldPath := getFieldPath(fieldType); fieldPath != "" {
		return fieldPath
	}

	switch {
	case field.HasName:
		return string(field.Name)
	case field.HasTypeName:
		return string(field.TypeName)
	default:
		return fmt.Sprintf("lookup_index_%d", field.Type.Int64())
	}
}

type FieldDecoder interface {
	Decode(decoder *scale.Decoder) (any, error)
}

type Field struct {
	Name         string
	FieldDecoder FieldDecoder
}

type EventDecoder struct {
	Name   string
	Fields []*Field
}

func (e *EventDecoder) Decode(decoder *scale.Decoder) (map[string]any, error) {
	fieldMap := make(map[string]any)

	for _, field := range e.Fields {
		value, err := field.FieldDecoder.Decode(decoder)

		if err != nil {
			return nil, err
		}

		fieldMap[field.Name] = value
	}

	return fieldMap, nil
}

type ValueDecoder[T any] struct{}

func (v *ValueDecoder[T]) Decode(decoder *scale.Decoder) (any, error) {
	var t T

	if err := decoder.Decode(&t); err != nil {
		return nil, err
	}

	return t, nil
}

type CompositeDecoder struct {
	Fields []*Field
}

func (e *CompositeDecoder) Decode(decoder *scale.Decoder) (any, error) {
	fieldMap := make(map[string]any)

	for _, field := range e.Fields {
		value, err := field.FieldDecoder.Decode(decoder)

		if err != nil {
			return nil, err
		}

		fieldMap[field.Name] = value
	}

	return fieldMap, nil
}

type VariantDecoder struct {
	FieldDecoderMap map[byte]FieldDecoder
}

func (v *VariantDecoder) Decode(decoder *scale.Decoder) (any, error) {
	variantByte, err := decoder.ReadOneByte()

	if err != nil {
		return nil, fmt.Errorf("couldn't read variant byte: %w", err)
	}

	variantDecoder, ok := v.FieldDecoderMap[variantByte]

	if !ok {
		return nil, fmt.Errorf("variant decoder for variant %d not found", variantByte)
	}

	if _, ok := variantDecoder.(*NoopDecoder); ok {
		return variantByte, nil
	}

	return variantDecoder.Decode(decoder)
}

type NoopDecoder struct{}

func (n *NoopDecoder) Decode(_ *scale.Decoder) (any, error) {
	return nil, nil
}

type ArrayDecoder struct {
	Length      uint
	ItemDecoder FieldDecoder
}

func (a *ArrayDecoder) Decode(decoder *scale.Decoder) (any, error) {
	if a.ItemDecoder == nil {
		return nil, errors.New("array item decoder not found")
	}

	slice := make([]any, 0, a.Length)

	for i := uint(0); i < a.Length; i++ {
		item, err := a.ItemDecoder.Decode(decoder)

		if err != nil {
			return nil, err
		}

		slice = append(slice, item)
	}

	return slice, nil
}

type SliceDecoder struct {
	ItemDecoder FieldDecoder
}

func (s *SliceDecoder) Decode(decoder *scale.Decoder) (any, error) {
	if s.ItemDecoder == nil {
		return nil, errors.New("slice item decoder not found")
	}

	sliceLen, err := decoder.DecodeUintCompact()

	if err != nil {
		return nil, fmt.Errorf("couldn't decode slice length: %w", err)
	}

	slice := make([]any, 0, sliceLen.Uint64())

	for i := uint64(0); i < sliceLen.Uint64(); i++ {
		item, err := s.ItemDecoder.Decode(decoder)

		if err != nil {
			return nil, err
		}

		slice = append(slice, item)
	}

	return slice, nil
}

type RecursiveDecoder struct {
	FieldDecoder FieldDecoder
}

func (r *RecursiveDecoder) Decode(decoder *scale.Decoder) (any, error) {
	if r.FieldDecoder == nil {
		return nil, errors.New("recursive field decoder not found")
	}

	return r.FieldDecoder.Decode(decoder)
}
