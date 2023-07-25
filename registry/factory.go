package registry

import (
	"errors"
	"fmt"
	"strings"

	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

//go:generate mockery --name Factory --structname FactoryMock --filename factory_mock.go --inpackage

// Factory is the interface responsible for generating the according registries from the metadata.
type Factory interface {
	CreateCallRegistry(meta *types.Metadata) (CallRegistry, error)
	CreateErrorRegistry(meta *types.Metadata) (ErrorRegistry, error)
	CreateEventRegistry(meta *types.Metadata) (EventRegistry, error)
}

// CallRegistry maps a call name to its TypeDecoder.
type CallRegistry map[types.CallIndex]*TypeDecoder

type ErrorID struct {
	ModuleIndex types.U8
	ErrorIndex  [4]types.U8
}

// ErrorRegistry maps an error name to its TypeDecoder.
type ErrorRegistry map[ErrorID]*TypeDecoder

// EventRegistry maps an event ID to its TypeDecoder.
type EventRegistry map[types.EventID]*TypeDecoder

// FieldOverride is used to override the default FieldDecoder for a particular type.
type FieldOverride struct {
	FieldLookupIndex int64
	FieldDecoder     FieldDecoder
}

type factory struct {
	fieldStorage          map[int64]FieldDecoder
	recursiveFieldStorage map[int64]*RecursiveDecoder
	fieldOverrides        []FieldOverride
}

// NewFactory creates a new Factory using the provided overrides, if any.
func NewFactory(fieldOverrides ...FieldOverride) Factory {
	f := &factory{}
	f.fieldOverrides = fieldOverrides

	return f
}

func (f *factory) resetStorages() {
	f.fieldStorage = make(map[int64]FieldDecoder)
	f.recursiveFieldStorage = make(map[int64]*RecursiveDecoder)

	for _, fieldOverride := range f.fieldOverrides {
		f.fieldStorage[fieldOverride.FieldLookupIndex] = fieldOverride.FieldDecoder
	}
}

// CreateErrorRegistry creates the registry that contains the types for errors.
// nolint:dupl
func (f *factory) CreateErrorRegistry(meta *types.Metadata) (ErrorRegistry, error) {
	f.resetStorages()

	errorRegistry := make(map[ErrorID]*TypeDecoder)

	for _, mod := range meta.AsMetadataV14.Pallets {
		if !mod.HasErrors {
			continue
		}

		errorsType, ok := meta.AsMetadataV14.EfficientLookup[mod.Errors.Type.Int64()]

		if !ok {
			return nil, ErrErrorsTypeNotFound.WithMsg("errors type '%d', module '%s'", mod.Errors.Type.Int64(), mod.Name)
		}

		if !errorsType.Def.IsVariant {
			return nil, ErrErrorsTypeNotVariant.WithMsg("errors type '%d', module '%s'", mod.Errors.Type.Int64(), mod.Name)
		}

		for _, errorVariant := range errorsType.Def.Variant.Variants {
			errorName := fmt.Sprintf("%s.%s", mod.Name, errorVariant.Name)

			errorFields, err := f.getTypeFields(meta, errorVariant.Fields)

			if err != nil {
				return nil, ErrErrorFieldsRetrieval.WithMsg(errorName).Wrap(err)
			}

			errorID := ErrorID{
				ModuleIndex: mod.Index,
				ErrorIndex:  [4]types.U8{errorVariant.Index},
			}

			errorRegistry[errorID] = &TypeDecoder{
				Name:   errorName,
				Fields: errorFields,
			}
		}
	}

	if err := f.resolveRecursiveDecoders(); err != nil {
		return nil, ErrRecursiveDecodersResolving.Wrap(err)
	}

	return errorRegistry, nil
}

// CreateCallRegistry creates the registry that contains the types for calls.
// nolint:dupl
func (f *factory) CreateCallRegistry(meta *types.Metadata) (CallRegistry, error) {
	f.resetStorages()

	callRegistry := make(map[types.CallIndex]*TypeDecoder)

	for _, mod := range meta.AsMetadataV14.Pallets {
		if !mod.HasCalls {
			continue
		}

		callsType, ok := meta.AsMetadataV14.EfficientLookup[mod.Calls.Type.Int64()]

		if !ok {
			return nil, ErrCallsTypeNotFound.WithMsg("calls type '%d', module '%s'", mod.Calls.Type.Int64(), mod.Name)
		}

		if !callsType.Def.IsVariant {
			return nil, ErrCallsTypeNotVariant.WithMsg("calls type '%d', module '%s'", mod.Calls.Type.Int64(), mod.Name)
		}

		for _, callVariant := range callsType.Def.Variant.Variants {
			callIndex := types.CallIndex{
				SectionIndex: uint8(mod.Index),
				MethodIndex:  uint8(callVariant.Index),
			}

			callName := fmt.Sprintf("%s.%s", mod.Name, callVariant.Name)

			callFields, err := f.getTypeFields(meta, callVariant.Fields)

			if err != nil {
				return nil, ErrCallFieldsRetrieval.WithMsg(callName).Wrap(err)
			}

			callRegistry[callIndex] = &TypeDecoder{
				Name:   callName,
				Fields: callFields,
			}
		}
	}

	if err := f.resolveRecursiveDecoders(); err != nil {
		return nil, ErrRecursiveDecodersResolving.Wrap(err)
	}

	return callRegistry, nil
}

// CreateEventRegistry creates the registry that contains the types for events.
func (f *factory) CreateEventRegistry(meta *types.Metadata) (EventRegistry, error) {
	f.resetStorages()

	eventRegistry := make(map[types.EventID]*TypeDecoder)

	for _, mod := range meta.AsMetadataV14.Pallets {
		if !mod.HasEvents {
			continue
		}

		eventsType, ok := meta.AsMetadataV14.EfficientLookup[mod.Events.Type.Int64()]

		if !ok {
			return nil, ErrEventsTypeNotFound.WithMsg("events type '%d', module '%s'", mod.Events.Type.Int64(), mod.Name)
		}

		if !eventsType.Def.IsVariant {
			return nil, ErrEventsTypeNotVariant.WithMsg("events type '%d', module '%s'", mod.Events.Type.Int64(), mod.Name)
		}

		for _, eventVariant := range eventsType.Def.Variant.Variants {
			eventID := types.EventID{byte(mod.Index), byte(eventVariant.Index)}

			eventName := fmt.Sprintf("%s.%s", mod.Name, eventVariant.Name)

			eventFields, err := f.getTypeFields(meta, eventVariant.Fields)

			if err != nil {
				return nil, ErrEventFieldsRetrieval.WithMsg(eventName).Wrap(err)
			}

			eventRegistry[eventID] = &TypeDecoder{
				Name:   eventName,
				Fields: eventFields,
			}
		}
	}

	if err := f.resolveRecursiveDecoders(); err != nil {
		return nil, ErrRecursiveDecodersResolving.Wrap(err)
	}

	return eventRegistry, nil
}

// resolveRecursiveDecoders resolves all recursive decoders with their according FieldDecoder.
// nolint:lll
func (f *factory) resolveRecursiveDecoders() error {
	for recursiveFieldLookupIndex, recursiveFieldDecoder := range f.recursiveFieldStorage {
		if recursiveFieldDecoder.FieldDecoder != nil {
			// Skip if the inner FieldDecoder is present, this could be an override.
			continue
		}

		fieldDecoder, ok := f.fieldStorage[recursiveFieldLookupIndex]

		if !ok {
			return ErrFieldDecoderForRecursiveFieldNotFound.
				WithMsg(
					"recursive field lookup index %d",
					recursiveFieldLookupIndex,
				)
		}

		if _, ok := fieldDecoder.(*RecursiveDecoder); ok {
			return ErrRecursiveFieldResolving.
				WithMsg(
					"recursive field lookup index %d",
					recursiveFieldLookupIndex,
				)
		}

		recursiveFieldDecoder.FieldDecoder = fieldDecoder
	}

	return nil
}

// getTypeFields parses and returns all Field(s) for a type.
func (f *factory) getTypeFields(meta *types.Metadata, fields []types.Si1Field) ([]*Field, error) {
	var typeFields []*Field

	for _, field := range fields {
		fieldType, ok := meta.AsMetadataV14.EfficientLookup[field.Type.Int64()]

		if !ok {
			return nil, ErrFieldTypeNotFound.WithMsg(string(field.Name))
		}

		fieldName := getFullFieldName(field, fieldType)

		if storedFieldDecoder, ok := f.getStoredFieldDecoder(field.Type.Int64()); ok {
			typeFields = append(typeFields, &Field{
				Name:         fieldName,
				FieldDecoder: storedFieldDecoder,
				LookupIndex:  field.Type.Int64(),
			})
			continue
		}

		fieldTypeDef := fieldType.Def

		fieldDecoder, err := f.getFieldDecoder(meta, fieldName, fieldTypeDef)

		if err != nil {
			return nil, ErrFieldDecoderRetrieval.WithMsg(fieldName).Wrap(err)
		}

		f.fieldStorage[field.Type.Int64()] = fieldDecoder

		typeFields = append(typeFields, &Field{
			Name:         fieldName,
			FieldDecoder: fieldDecoder,
			LookupIndex:  field.Type.Int64(),
		})
	}

	return typeFields, nil
}

// getFieldDecoder returns the FieldDecoder based on the provided type definition.
// nolint:funlen
func (f *factory) getFieldDecoder(
	meta *types.Metadata,
	fieldName string,
	typeDef types.Si1TypeDef,
) (FieldDecoder, error) {
	switch {
	case typeDef.IsCompact:
		compactFieldType, ok := meta.AsMetadataV14.EfficientLookup[typeDef.Compact.Type.Int64()]

		if !ok {
			return nil, ErrCompactFieldTypeNotFound.WithMsg(fieldName)
		}

		return f.getCompactFieldDecoder(meta, fieldName, compactFieldType.Def)
	case typeDef.IsComposite:
		compositeDecoder := &CompositeDecoder{
			FieldName: fieldName,
		}

		fields, err := f.getTypeFields(meta, typeDef.Composite.Fields)

		if err != nil {
			return nil, ErrCompositeTypeFieldsRetrieval.WithMsg(fieldName).Wrap(err)
		}

		compositeDecoder.Fields = fields

		return compositeDecoder, nil
	case typeDef.IsVariant:
		return f.getVariantFieldDecoder(meta, typeDef)
	case typeDef.IsPrimitive:
		return getPrimitiveDecoder(typeDef.Primitive.Si0TypeDefPrimitive)
	case typeDef.IsArray:
		arrayFieldType, ok := meta.AsMetadataV14.EfficientLookup[typeDef.Array.Type.Int64()]

		if !ok {
			return nil, ErrArrayFieldTypeNotFound.WithMsg(fieldName)
		}

		return f.getArrayFieldDecoder(uint(typeDef.Array.Len), meta, fieldName, arrayFieldType.Def)
	case typeDef.IsSequence:
		vectorFieldType, ok := meta.AsMetadataV14.EfficientLookup[typeDef.Sequence.Type.Int64()]

		if !ok {
			return nil, ErrVectorFieldTypeNotFound.WithMsg(fieldName)
		}

		return f.getSliceFieldDecoder(meta, fieldName, vectorFieldType.Def)
	case typeDef.IsTuple:
		if typeDef.Tuple == nil {
			return &NoopDecoder{}, nil
		}

		return f.getTupleFieldDecoder(meta, fieldName, typeDef.Tuple)
	case typeDef.IsBitSequence:
		return f.getBitSequenceDecoder(meta, fieldName, typeDef.BitSequence)
	default:
		return nil, ErrFieldTypeDefinitionNotSupported.WithMsg(fieldName)
	}
}

const (
	variantItemFieldNameFormat = "variant_item_%d"
)

// getVariantFieldDecoder parses a variant type definition and returns a VariantDecoder.
func (f *factory) getVariantFieldDecoder(meta *types.Metadata, typeDef types.Si1TypeDef) (FieldDecoder, error) {
	variantDecoder := &VariantDecoder{}

	fieldDecoderMap := make(map[byte]FieldDecoder)

	for i, variant := range typeDef.Variant.Variants {
		if len(variant.Fields) == 0 {
			fieldDecoderMap[byte(variant.Index)] = &NoopDecoder{}
			continue
		}

		variantFieldName := fmt.Sprintf(variantItemFieldNameFormat, i)

		compositeDecoder := &CompositeDecoder{
			FieldName: variantFieldName,
		}

		fields, err := f.getTypeFields(meta, variant.Fields)

		if err != nil {
			return nil, ErrVariantTypeFieldsRetrieval.WithMsg("variant '%d'", variant.Index).Wrap(err)
		}

		compositeDecoder.Fields = fields

		fieldDecoderMap[byte(variant.Index)] = compositeDecoder
	}

	variantDecoder.FieldDecoderMap = fieldDecoderMap

	return variantDecoder, nil
}

const (
	tupleItemFieldNameFormat = "tuple_item_%d"
)

// getCompactFieldDecoder parses a compact type definition and returns the according field decoder.
// nolint:funlen,lll
func (f *factory) getCompactFieldDecoder(meta *types.Metadata, fieldName string, typeDef types.Si1TypeDef) (FieldDecoder, error) {
	switch {
	case typeDef.IsPrimitive:
		return &ValueDecoder[types.UCompact]{}, nil
	case typeDef.IsTuple:
		if typeDef.Tuple == nil {
			return &ValueDecoder[any]{}, nil
		}

		compositeDecoder := &CompositeDecoder{
			FieldName: fieldName,
		}

		for i, item := range typeDef.Tuple {
			itemTypeDef, ok := meta.AsMetadataV14.EfficientLookup[item.Int64()]

			if !ok {
				return nil, ErrCompactTupleItemTypeNotFound.WithMsg("tuple item '%d'", item.Int64())
			}

			fieldName := fmt.Sprintf(tupleItemFieldNameFormat, i)

			itemFieldDecoder, err := f.getCompactFieldDecoder(meta, fieldName, itemTypeDef.Def)

			if err != nil {
				return nil, ErrCompactTupleItemFieldDecoderRetrieval.
					WithMsg("tuple item '%d'", item.Int64()).
					Wrap(err)
			}

			compositeDecoder.Fields = append(compositeDecoder.Fields, &Field{
				Name:         fieldName,
				FieldDecoder: itemFieldDecoder,
				LookupIndex:  item.Int64(),
			})
		}

		return compositeDecoder, nil
	case typeDef.IsComposite:
		compactCompositeFields := typeDef.Composite.Fields

		compositeDecoder := &CompositeDecoder{
			FieldName: fieldName,
		}

		for _, compactCompositeField := range compactCompositeFields {
			compactCompositeFieldType, ok := meta.AsMetadataV14.EfficientLookup[compactCompositeField.Type.Int64()]

			if !ok {
				return nil, ErrCompactCompositeFieldTypeNotFound
			}

			compactFieldName := getFullFieldName(compactCompositeField, compactCompositeFieldType)

			compactCompositeDecoder, err := f.getCompactFieldDecoder(meta, compactFieldName, compactCompositeFieldType.Def)

			if err != nil {
				return nil, ErrCompactCompositeFieldDecoderRetrieval.Wrap(err)
			}

			compositeDecoder.Fields = append(compositeDecoder.Fields, &Field{
				Name:         compactFieldName,
				FieldDecoder: compactCompositeDecoder,
				LookupIndex:  compactCompositeField.Type.Int64(),
			})
		}

		return compositeDecoder, nil
	default:
		return nil, errors.New("unsupported compact field type")
	}
}

// getArrayFieldDecoder parses an array type definition and returns an ArrayDecoder.
// nolint:lll
func (f *factory) getArrayFieldDecoder(arrayLen uint, meta *types.Metadata, fieldName string, typeDef types.Si1TypeDef) (FieldDecoder, error) {
	itemFieldDecoder, err := f.getFieldDecoder(meta, fieldName, typeDef)

	if err != nil {
		return nil, ErrArrayItemFieldDecoderRetrieval.Wrap(err)
	}

	return &ArrayDecoder{Length: arrayLen, ItemDecoder: itemFieldDecoder}, nil
}

// getSliceFieldDecoder parses a slice type definition and returns an SliceDecoder.
func (f *factory) getSliceFieldDecoder(
	meta *types.Metadata,
	fieldName string,
	typeDef types.Si1TypeDef,
) (FieldDecoder, error) {
	itemFieldDecoder, err := f.getFieldDecoder(meta, fieldName, typeDef)

	if err != nil {
		return nil, ErrSliceItemFieldDecoderRetrieval.Wrap(err)
	}

	return &SliceDecoder{itemFieldDecoder}, nil
}

// getTupleFieldDecoder parses a tuple type definition and returns a CompositeDecoder.
func (f *factory) getTupleFieldDecoder(
	meta *types.Metadata,
	fieldName string,
	tuple types.Si1TypeDefTuple,
) (FieldDecoder, error) {
	compositeDecoder := &CompositeDecoder{
		FieldName: fieldName,
	}

	for i, item := range tuple {
		itemTypeDef, ok := meta.AsMetadataV14.EfficientLookup[item.Int64()]

		if !ok {
			return nil, ErrTupleItemTypeNotFound.WithMsg("tuple item '%d'", i)
		}

		tupleFieldName := fmt.Sprintf(tupleItemFieldNameFormat, i)

		itemFieldDecoder, err := f.getFieldDecoder(meta, tupleFieldName, itemTypeDef.Def)

		if err != nil {
			return nil, ErrTupleItemFieldDecoderRetrieval.Wrap(err)
		}

		compositeDecoder.Fields = append(compositeDecoder.Fields, &Field{
			Name:         tupleFieldName,
			FieldDecoder: itemFieldDecoder,
			LookupIndex:  item.Int64(),
		})
	}

	return compositeDecoder, nil
}

func (f *factory) getBitSequenceDecoder(
	meta *types.Metadata,
	fieldName string,
	bitSequenceTypeDef types.Si1TypeDefBitSequence,
) (FieldDecoder, error) {
	bitStoreType, ok := meta.AsMetadataV14.EfficientLookup[bitSequenceTypeDef.BitStoreType.Int64()]

	if !ok {
		return nil, ErrBitStoreTypeNotFound.WithMsg(fieldName)
	}

	if bitStoreType.Def.Primitive.Si0TypeDefPrimitive != types.IsU8 {
		return nil, ErrBitStoreTypeNotSupported.WithMsg(fieldName)
	}

	bitOrderType, ok := meta.AsMetadataV14.EfficientLookup[bitSequenceTypeDef.BitOrderType.Int64()]

	if !ok {
		return nil, ErrBitOrderTypeNotFound.WithMsg(fieldName)
	}

	bitOrder, err := types.NewBitOrderFromString(getBitOrderString(bitOrderType.Path))

	if err != nil {
		return nil, ErrBitOrderCreation.Wrap(err)
	}

	bitSequenceDecoder := &BitSequenceDecoder{
		FieldName: fieldName,
		BitOrder:  bitOrder,
	}

	return bitSequenceDecoder, nil
}

func getBitOrderString(path types.Si1Path) string {
	pathLen := len(path)

	if pathLen == 0 {
		return ""
	}

	return string(path[pathLen-1])
}

// getPrimitiveDecoder parses a primitive type definition and returns a ValueDecoder.
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
		return nil, ErrPrimitiveTypeNotSupported.WithMsg("primitive type %v", primitiveTypeDef)
	}
}

// getStoredFieldDecoder will attempt to return a FieldDecoder from storage,
// and perform an extra check for recursive decoders.
func (f *factory) getStoredFieldDecoder(fieldLookupIndex int64) (FieldDecoder, bool) {
	if ft, ok := f.fieldStorage[fieldLookupIndex]; ok {
		if rt, ok := ft.(*RecursiveDecoder); ok {
			f.recursiveFieldStorage[fieldLookupIndex] = rt
		}

		return ft, ok
	}

	// Ensure that a recursive type such as Xcm::TransferReserveAsset does not cause an infinite loop
	// by adding the RecursiveDecoder the first time the field is encountered.
	f.fieldStorage[fieldLookupIndex] = &RecursiveDecoder{}

	return nil, false
}

const (
	fieldSeparator    = "."
	lookupIndexFormat = "lookup_index_%d"
)

func getFieldPath(fieldType *types.Si1Type) string {
	var nameParts []string

	for _, pathEntry := range fieldType.Path {
		nameParts = append(nameParts, string(pathEntry))
	}

	return strings.Join(nameParts, fieldSeparator)
}

func getFullFieldName(field types.Si1Field, fieldType *types.Si1Type) string {
	fieldName := getFieldName(field)

	if fieldPath := getFieldPath(fieldType); fieldPath != "" {
		return fmt.Sprintf("%s%s%s", fieldPath, fieldSeparator, fieldName)
	}

	return getFieldName(field)
}

func getFieldName(field types.Si1Field) string {
	switch {
	case field.HasName:
		return string(field.Name)
	case field.HasTypeName:
		return string(field.TypeName)
	default:
		return fmt.Sprintf(lookupIndexFormat, field.Type.Int64())
	}
}

// FieldDecoder is the interface implemented by all the different types that are available.
type FieldDecoder interface {
	Decode(decoder *scale.Decoder) (any, error)
}

// NoopDecoder is a FieldDecoder that does not decode anything. It comes in handy for nil tuples or variants
// with no inner types.
type NoopDecoder struct{}

func (n *NoopDecoder) Decode(_ *scale.Decoder) (any, error) {
	return nil, nil
}

// VariantDecoder holds a FieldDecoder for each variant/enum.
type VariantDecoder struct {
	FieldDecoderMap map[byte]FieldDecoder
}

func (v *VariantDecoder) Decode(decoder *scale.Decoder) (any, error) {
	variantByte, err := decoder.ReadOneByte()

	if err != nil {
		return nil, ErrVariantByteDecoding.Wrap(err)
	}

	variantDecoder, ok := v.FieldDecoderMap[variantByte]

	if !ok {
		return nil, ErrVariantFieldDecoderNotFound.WithMsg("variant '%d'", variantByte)
	}

	if _, ok := variantDecoder.(*NoopDecoder); ok {
		return variantByte, nil
	}

	return variantDecoder.Decode(decoder)
}

// ArrayDecoder holds information about the length of the array and the FieldDecoder used for its items.
type ArrayDecoder struct {
	Length      uint
	ItemDecoder FieldDecoder
}

func (a *ArrayDecoder) Decode(decoder *scale.Decoder) (any, error) {
	if a.ItemDecoder == nil {
		return nil, ErrArrayItemDecoderNotFound
	}

	slice := make([]any, 0, a.Length)

	for i := uint(0); i < a.Length; i++ {
		item, err := a.ItemDecoder.Decode(decoder)

		if err != nil {
			return nil, ErrArrayItemDecoding.Wrap(err)
		}

		slice = append(slice, item)
	}

	return slice, nil
}

// SliceDecoder holds a FieldDecoder for the items of a vector/slice.
type SliceDecoder struct {
	ItemDecoder FieldDecoder
}

func (s *SliceDecoder) Decode(decoder *scale.Decoder) (any, error) {
	if s.ItemDecoder == nil {
		return nil, ErrSliceItemDecoderNotFound
	}

	sliceLen, err := decoder.DecodeUintCompact()

	if err != nil {
		return nil, ErrSliceLengthDecoding.Wrap(err)
	}

	slice := make([]any, 0, sliceLen.Uint64())

	for i := uint64(0); i < sliceLen.Uint64(); i++ {
		item, err := s.ItemDecoder.Decode(decoder)

		if err != nil {
			return nil, ErrSliceItemDecoding.Wrap(err)
		}

		slice = append(slice, item)
	}

	return slice, nil
}

// CompositeDecoder holds all the information required to decoder a struct/composite.
type CompositeDecoder struct {
	FieldName string
	Fields    []*Field
}

func (e *CompositeDecoder) Decode(decoder *scale.Decoder) (any, error) {
	var decodedFields DecodedFields

	for _, field := range e.Fields {
		value, err := field.FieldDecoder.Decode(decoder)

		if err != nil {
			return nil, ErrCompositeFieldDecoding.Wrap(err)
		}

		decodedFields = append(decodedFields, &DecodedField{
			Name:        field.Name,
			Value:       value,
			LookupIndex: field.LookupIndex,
		})
	}

	return decodedFields, nil
}

// ValueDecoder decodes a primitive type.
type ValueDecoder[T any] struct{}

func (v *ValueDecoder[T]) Decode(decoder *scale.Decoder) (any, error) {
	var t T

	if err := decoder.Decode(&t); err != nil {
		return nil, ErrValueDecoding.Wrap(err)
	}

	return t, nil
}

// RecursiveDecoder is a wrapper for a FieldDecoder that is recursive.
type RecursiveDecoder struct {
	FieldDecoder FieldDecoder
}

func (r *RecursiveDecoder) Decode(decoder *scale.Decoder) (any, error) {
	if r.FieldDecoder == nil {
		return nil, ErrRecursiveFieldDecoderNotFound
	}

	return r.FieldDecoder.Decode(decoder)
}

// BitSequenceDecoder holds decoding information for a bit sequence.
type BitSequenceDecoder struct {
	FieldName string
	BitOrder  types.BitOrder
}

func (b *BitSequenceDecoder) Decode(decoder *scale.Decoder) (any, error) {
	bitVec := types.NewBitVec(b.BitOrder)

	if err := bitVec.Decode(*decoder); err != nil {
		return nil, ErrBitVecDecoding.Wrap(err)
	}

	return map[string]string{
		b.FieldName: bitVec.String(),
	}, nil
}

// TypeDecoder holds all information required to decode a particular type.
type TypeDecoder struct {
	Name   string
	Fields []*Field
}

func (t *TypeDecoder) Decode(decoder *scale.Decoder) (DecodedFields, error) {
	if t == nil {
		return nil, ErrNilTypeDecoder
	}

	var decodedFields DecodedFields

	for _, field := range t.Fields {
		decodedField, err := field.Decode(decoder)

		if err != nil {
			return nil, ErrTypeFieldDecoding.Wrap(err)
		}

		decodedFields = append(decodedFields, decodedField)
	}

	return decodedFields, nil
}

// Field represents one field of a TypeDecoder.
type Field struct {
	Name         string
	FieldDecoder FieldDecoder
	LookupIndex  int64
}

func (f *Field) Decode(decoder *scale.Decoder) (*DecodedField, error) {
	if f == nil {
		return nil, ErrNilField
	}

	if f.FieldDecoder == nil {
		return nil, ErrNilFieldDecoder
	}

	value, err := f.FieldDecoder.Decode(decoder)

	if err != nil {
		return nil, err
	}

	return &DecodedField{
		Name:        f.Name,
		Value:       value,
		LookupIndex: f.LookupIndex,
	}, nil
}

// DecodedField holds the name, value and lookup index of a field that was decoded.
type DecodedField struct {
	Name        string
	Value       any
	LookupIndex int64
}

func (d DecodedField) Encode(encoder scale.Encoder) error {
	if d.Value == nil {
		return nil
	}

	return encoder.Encode(d.Value)
}

type DecodedFields []*DecodedField

type DecodedFieldPredicateFn func(fieldIndex int, field *DecodedField) bool
type DecodedValueProcessingFn[T any] func(value any) (T, error)

// ProcessDecodedFieldValue applies the processing func to the value of the field
// that matches the provided predicate func.
func ProcessDecodedFieldValue[T any](
	decodedFields DecodedFields,
	fieldPredicateFn DecodedFieldPredicateFn,
	valueProcessingFn DecodedValueProcessingFn[T],
) (T, error) {
	var t T

	for decodedFieldIndex, decodedField := range decodedFields {
		if !fieldPredicateFn(decodedFieldIndex, decodedField) {
			continue
		}

		res, err := valueProcessingFn(decodedField.Value)

		if err != nil {
			return t, ErrDecodedFieldValueProcessingError.Wrap(err)
		}

		return res, nil
	}

	return t, ErrDecodedFieldNotFound
}

// GetDecodedFieldAsType returns the value of the field that matches the provided predicate func
// as the provided generic argument.
func GetDecodedFieldAsType[T any](
	decodedFields DecodedFields,
	fieldPredicateFn DecodedFieldPredicateFn,
) (T, error) {
	return ProcessDecodedFieldValue(
		decodedFields,
		fieldPredicateFn,
		func(value any) (T, error) {
			if res, ok := value.(T); ok {
				return res, nil
			}

			var t T

			err := fmt.Errorf("expected %T, got %T", t, value)

			return t, ErrDecodedFieldValueTypeMismatch.Wrap(err)
		},
	)
}

// GetDecodedFieldAsSliceOfType returns the value of the field that matches the provided predicate func
// as a slice of the provided generic argument.
func GetDecodedFieldAsSliceOfType[T any](
	decodedFields DecodedFields,
	fieldPredicateFn DecodedFieldPredicateFn,
) ([]T, error) {
	return ProcessDecodedFieldValue(
		decodedFields,
		fieldPredicateFn,
		func(value any) ([]T, error) {
			v, ok := value.([]any)

			if !ok {
				return nil, ErrDecodedFieldValueNotAGenericSlice
			}

			res, err := convertSliceToType[T](v)

			if err != nil {
				return nil, ErrDecodedFieldValueTypeMismatch.Wrap(err)
			}

			return res, nil
		},
	)
}

func convertSliceToType[T any](slice []any) ([]T, error) {
	res := make([]T, 0)

	for _, item := range slice {
		if v, ok := item.(T); ok {
			res = append(res, v)
			continue
		}

		var t T

		return nil, fmt.Errorf("expected %T, got %T", t, item)
	}

	return res, nil
}
