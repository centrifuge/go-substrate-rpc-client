package registry

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

// Factory is the interface responsible for generating the according registries from the metadata.
type Factory interface {
	CreateCallRegistry(meta *types.Metadata) (CallRegistry, error)
	CreateErrorRegistry(meta *types.Metadata) (ErrorRegistry, error)
	CreateEventRegistry(meta *types.Metadata) (EventRegistry, error)
}

// CallRegistry maps a call name to its Type.
type CallRegistry map[string]*Type

// ErrorRegistry maps an error name to its Type.
type ErrorRegistry map[string]*Type

// EventRegistry maps an event ID to its Type.
type EventRegistry map[types.EventID]*Type

type factory struct {
	fieldStorage          map[int64]FieldType
	recursiveFieldStorage map[int64]*RecursiveFieldType
}

// NewFactory creates a new Factory.
func NewFactory() Factory {
	return &factory{
		fieldStorage:          make(map[int64]FieldType),
		recursiveFieldStorage: make(map[int64]*RecursiveFieldType),
	}
}

// CreateErrorRegistry creates the registry that contains the types for errors.
// nolint:dupl
func (r *factory) CreateErrorRegistry(meta *types.Metadata) (ErrorRegistry, error) {
	errorRegistry := make(map[string]*Type)

	for _, mod := range meta.AsMetadataV14.Pallets {
		if !mod.HasErrors {
			continue
		}

		errorsType, ok := meta.AsMetadataV14.EfficientLookup[mod.Errors.Type.Int64()]

		if !ok {
			return nil, fmt.Errorf("errors type %d not found for module '%s'", mod.Errors.Type.Int64(), mod.Name)
		}

		if !errorsType.Def.IsVariant {
			return nil, fmt.Errorf("errors type %d for module '%s' is not a variant", mod.Errors.Type.Int64(), mod.Name)
		}

		for _, errorVariant := range errorsType.Def.Variant.Variants {
			errorName := fmt.Sprintf("%s.%s", mod.Name, errorVariant.Name)

			errorFields, err := r.getTypeFields(meta, errorVariant.Fields)

			if err != nil {
				return nil, fmt.Errorf("couldn't get fields for error '%s': %w", errorName, err)
			}

			errorRegistry[errorName] = &Type{
				Name:   errorName,
				Fields: errorFields,
			}
		}
	}

	if err := r.resolveRecursiveTypes(); err != nil {
		return nil, err
	}

	return errorRegistry, nil
}

// CreateCallRegistry creates the registry that contains the types for calls.
// nolint:dupl
func (r *factory) CreateCallRegistry(meta *types.Metadata) (CallRegistry, error) {
	callRegistry := make(map[string]*Type)

	for _, mod := range meta.AsMetadataV14.Pallets {
		if !mod.HasCalls {
			continue
		}

		callsType, ok := meta.AsMetadataV14.EfficientLookup[mod.Calls.Type.Int64()]

		if !ok {
			return nil, fmt.Errorf("calls type %d not found for module '%s'", mod.Calls.Type.Int64(), mod.Name)
		}

		if !callsType.Def.IsVariant {
			return nil, fmt.Errorf("calls type %d for module '%s' is not a variant", mod.Calls.Type.Int64(), mod.Name)
		}

		for _, callVariant := range callsType.Def.Variant.Variants {
			callName := fmt.Sprintf("%s.%s", mod.Name, callVariant.Name)

			callFields, err := r.getTypeFields(meta, callVariant.Fields)

			if err != nil {
				return nil, fmt.Errorf("couldn't get fields for call '%s': %w", callName, err)
			}

			callRegistry[callName] = &Type{
				Name:   callName,
				Fields: callFields,
			}
		}
	}

	if err := r.resolveRecursiveTypes(); err != nil {
		return nil, err
	}

	return callRegistry, nil
}

// CreateEventRegistry creates the registry that contains the types for events.
func (r *factory) CreateEventRegistry(meta *types.Metadata) (EventRegistry, error) {
	eventRegistry := make(map[types.EventID]*Type)

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

			eventFields, err := r.getTypeFields(meta, eventVariant.Fields)

			if err != nil {
				return nil, fmt.Errorf("couldn't get fields for event '%s': %w", eventName, err)
			}

			eventRegistry[eventID] = &Type{
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

// resolveRecursiveTypes resolves all recursive types with their according field type.
// nolint:lll
func (r *factory) resolveRecursiveTypes() error {
	for recursiveFieldLookupIndex, recursiveFieldType := range r.recursiveFieldStorage {
		fieldType, ok := r.fieldStorage[recursiveFieldLookupIndex]

		if !ok {
			return fmt.Errorf("couldn't get field type for recursive type %d", recursiveFieldLookupIndex)
		}

		if _, ok := fieldType.(*RecursiveFieldType); ok {
			return fmt.Errorf("recursive field type %d cannot be resolved with a non-recursive field type", recursiveFieldLookupIndex)
		}

		recursiveFieldType.ResolvedItemType = fieldType
	}

	return nil
}

// getTypeFields parses and returns all fields for a type.
func (r *factory) getTypeFields(meta *types.Metadata, fields []types.Si1Field) ([]*Field, error) {
	var typeFields []*Field

	for _, field := range fields {
		fieldType, ok := meta.AsMetadataV14.EfficientLookup[field.Type.Int64()]

		if !ok {
			return nil, fmt.Errorf("type not found for field '%s'", field.Name)
		}

		fieldName := getFieldName(field, fieldType)

		if storedFieldType, ok := r.getStoredFieldType(fieldName, field.Type.Int64()); ok {
			typeFields = append(typeFields, &Field{
				Name:        fieldName,
				FieldType:   storedFieldType,
				LookupIndex: field.Type.Int64(),
			})
			continue
		}

		fieldTypeDef := fieldType.Def

		resolvedFieldType, err := r.getFieldType(meta, fieldName, fieldTypeDef)

		if err != nil {
			return nil, fmt.Errorf("couldn't get field type for '%s': %w", fieldName, err)
		}

		r.storeFieldType(field.Type.Int64(), resolvedFieldType)

		typeFields = append(typeFields, &Field{
			Name:        fieldName,
			FieldType:   resolvedFieldType,
			LookupIndex: field.Type.Int64(),
		})
	}

	return typeFields, nil
}

// getFieldType returns the FieldType based on the provided type definition.
// nolint:funlen
func (r *factory) getFieldType(meta *types.Metadata, fieldName string, typeDef types.Si1TypeDef) (FieldType, error) {
	switch {
	case typeDef.IsCompact:
		compactFieldType, ok := meta.AsMetadataV14.EfficientLookup[typeDef.Compact.Type.Int64()]

		if !ok {
			return nil, errors.New("type not found for compact field")
		}

		return r.getCompactFieldType(meta, fieldName, compactFieldType.Def)
	case typeDef.IsComposite:
		compositeFieldType := &CompositeFieldType{
			FieldName: fieldName,
		}

		fields, err := r.getTypeFields(meta, typeDef.Composite.Fields)

		if err != nil {
			return nil, fmt.Errorf("couldn't get composite fields: %w", err)
		}

		compositeFieldType.Fields = fields

		return compositeFieldType, nil
	case typeDef.IsVariant:
		return r.getVariantFieldType(meta, typeDef)
	case typeDef.IsPrimitive:
		return getPrimitiveType(typeDef.Primitive.Si0TypeDefPrimitive)
	case typeDef.IsArray:
		arrayFieldType, ok := meta.AsMetadataV14.EfficientLookup[typeDef.Array.Type.Int64()]

		if !ok {
			return nil, fmt.Errorf("type not found for array field")
		}

		return r.getArrayFieldType(uint(typeDef.Array.Len), meta, fieldName, arrayFieldType.Def)
	case typeDef.IsSequence:
		vectorFieldType, ok := meta.AsMetadataV14.EfficientLookup[typeDef.Sequence.Type.Int64()]

		if !ok {
			return nil, errors.New("type not found for vector field")
		}

		return r.getSliceFieldType(meta, fieldName, vectorFieldType.Def)
	case typeDef.IsTuple:
		if typeDef.Tuple == nil {
			return &PrimitiveFieldType[[]any]{}, nil
		}

		return r.getTupleType(meta, fieldName, typeDef.Tuple)
	case typeDef.IsBitSequence:
		bitStoreType, ok := meta.AsMetadataV14.EfficientLookup[typeDef.BitSequence.BitStoreType.Int64()]

		if !ok {
			return nil, errors.New("bit store type not found")
		}

		bitStoreFieldType, err := r.getFieldType(meta, "bitStoreType", bitStoreType.Def)

		if err != nil {
			return nil, fmt.Errorf("couldn't get bit store field type: %w", err)
		}

		bitOrderType, ok := meta.AsMetadataV14.EfficientLookup[typeDef.BitSequence.BitOrderType.Int64()]

		if !ok {
			return nil, errors.New("bit order type not found")
		}

		bitOrderFieldType, err := r.getFieldType(meta, "bitOrderType", bitOrderType.Def)

		if err != nil {
			return nil, fmt.Errorf("couldn't get bit order field type: %w", err)
		}

		return &BitSequenceType{
			BitStoreType: bitStoreFieldType,
			BitOrderType: bitOrderFieldType,
		}, nil
	default:
		return nil, errors.New("unsupported field type definition")
	}
}

const (
	variantItemFieldNameFormat = "variant_item_%d"
)

// getVariantFieldType parses a variant type definition and returns a VariantFieldType.
func (r *factory) getVariantFieldType(meta *types.Metadata, typeDef types.Si1TypeDef) (FieldType, error) {
	variantFieldType := &VariantFieldType{}

	fieldTypeMap := make(map[byte]FieldType)

	for i, variant := range typeDef.Variant.Variants {
		if len(variant.Fields) == 0 {
			fieldTypeMap[byte(variant.Index)] = &PrimitiveFieldType[byte]{}
			continue
		}

		variantFieldName := fmt.Sprintf(variantItemFieldNameFormat, i)

		compositeFieldType := &CompositeFieldType{
			FieldName: variantFieldName,
		}

		fields, err := r.getTypeFields(meta, variant.Fields)

		if err != nil {
			return nil, fmt.Errorf("couldn't get field types for variant '%d': %w", variant.Index, err)
		}

		compositeFieldType.Fields = fields

		fieldTypeMap[byte(variant.Index)] = compositeFieldType
	}

	variantFieldType.FieldTypeMap = fieldTypeMap

	return variantFieldType, nil
}

const (
	tupleItemFieldNameFormat = "tuple_item_%d"
)

// getCompactFieldType parses a compact type definition and returns the according field type.
// nolint:funlen,lll
func (r *factory) getCompactFieldType(meta *types.Metadata, fieldName string, typeDef types.Si1TypeDef) (FieldType, error) {
	switch {
	case typeDef.IsPrimitive:
		return &PrimitiveFieldType[types.UCompact]{}, nil
	case typeDef.IsTuple:
		if typeDef.Tuple == nil {
			return &PrimitiveFieldType[any]{}, nil
		}

		compositeFieldType := &CompositeFieldType{
			FieldName: fieldName,
		}

		for i, item := range typeDef.Tuple {
			itemTypeDef, ok := meta.AsMetadataV14.EfficientLookup[item.Int64()]

			if !ok {
				return nil, fmt.Errorf("type definition for tuple item %d not found", item.Int64())
			}

			fieldName := fmt.Sprintf(tupleItemFieldNameFormat, i)

			itemFieldType, err := r.getCompactFieldType(meta, fieldName, itemTypeDef.Def)

			if err != nil {
				return nil, fmt.Errorf("couldn't get tuple field type: %w", err)
			}

			compositeFieldType.Fields = append(compositeFieldType.Fields, &Field{
				Name:        fieldName,
				FieldType:   itemFieldType,
				LookupIndex: item.Int64(),
			})
		}

		return compositeFieldType, nil
	case typeDef.IsComposite:
		compactCompositeFields := typeDef.Composite.Fields

		compositeFieldType := &CompositeFieldType{
			FieldName: fieldName,
		}

		for _, compactCompositeField := range compactCompositeFields {
			compactCompositeFieldType, ok := meta.AsMetadataV14.EfficientLookup[compactCompositeField.Type.Int64()]

			if !ok {
				return nil, errors.New("compact composite field type not found")
			}

			compactFieldName := getFieldName(compactCompositeField, compactCompositeFieldType)

			compactCompositeType, err := r.getCompactFieldType(meta, compactFieldName, compactCompositeFieldType.Def)

			if err != nil {
				return nil, fmt.Errorf("couldn't decode compact composite type: %w", err)
			}

			compositeFieldType.Fields = append(compositeFieldType.Fields, &Field{
				Name:        compactFieldName,
				FieldType:   compactCompositeType,
				LookupIndex: compactCompositeField.Type.Int64(),
			})
		}

		return compositeFieldType, nil
	default:
		return nil, errors.New("unsupported compact field type")
	}
}

// getArrayFieldType parses an array type definition and returns an ArrayFieldType.
// nolint:lll
func (r *factory) getArrayFieldType(arrayLen uint, meta *types.Metadata, fieldName string, typeDef types.Si1TypeDef) (FieldType, error) {
	itemFieldType, err := r.getFieldType(meta, fieldName, typeDef)

	if err != nil {
		return nil, fmt.Errorf("couldn't get array item field type: %w", err)
	}

	arrayFieldType := &ArrayFieldType{Length: arrayLen, ItemType: itemFieldType}

	return arrayFieldType, nil
}

// getSliceFieldType parses a slice type definition and returns an SliceFieldType.
// nolint:lll
func (r *factory) getSliceFieldType(meta *types.Metadata, fieldName string, typeDef types.Si1TypeDef) (FieldType, error) {
	itemFieldType, err := r.getFieldType(meta, fieldName, typeDef)

	if err != nil {
		return nil, fmt.Errorf("couldn't get slice item field type: %w", err)
	}

	sliceFieldType := &SliceFieldType{itemFieldType}

	return sliceFieldType, nil
}

// getTupleType parses a tuple type definition and returns a CompositeFieldType.
func (r *factory) getTupleType(meta *types.Metadata, fieldName string, tuple types.Si1TypeDefTuple) (FieldType, error) {
	compositeFieldType := &CompositeFieldType{
		FieldName: fieldName,
	}

	for i, item := range tuple {
		itemTypeDef, ok := meta.AsMetadataV14.EfficientLookup[item.Int64()]

		if !ok {
			return nil, fmt.Errorf("type definition for tuple item %d not found", i)
		}

		tupleFieldName := fmt.Sprintf(tupleItemFieldNameFormat, i)

		itemFieldType, err := r.getFieldType(meta, tupleFieldName, itemTypeDef.Def)

		if err != nil {
			return nil, fmt.Errorf("couldn't get tuple field type: %w", err)
		}

		compositeFieldType.Fields = append(compositeFieldType.Fields, &Field{
			Name:        tupleFieldName,
			FieldType:   itemFieldType,
			LookupIndex: item.Int64(),
		})
	}

	return compositeFieldType, nil
}

// getPrimitiveType parses a primitive type definition and returns a PrimitiveFieldType.
func getPrimitiveType(primitiveTypeDef types.Si0TypeDefPrimitive) (FieldType, error) {
	switch primitiveTypeDef {
	case types.IsBool:
		return &PrimitiveFieldType[bool]{}, nil
	case types.IsChar:
		return &PrimitiveFieldType[byte]{}, nil
	case types.IsStr:
		return &PrimitiveFieldType[string]{}, nil
	case types.IsU8:
		return &PrimitiveFieldType[types.U8]{}, nil
	case types.IsU16:
		return &PrimitiveFieldType[types.U16]{}, nil
	case types.IsU32:
		return &PrimitiveFieldType[types.U32]{}, nil
	case types.IsU64:
		return &PrimitiveFieldType[types.U64]{}, nil
	case types.IsU128:
		return &PrimitiveFieldType[types.U128]{}, nil
	case types.IsU256:
		return &PrimitiveFieldType[types.U256]{}, nil
	case types.IsI8:
		return &PrimitiveFieldType[types.I8]{}, nil
	case types.IsI16:
		return &PrimitiveFieldType[types.I16]{}, nil
	case types.IsI32:
		return &PrimitiveFieldType[types.I32]{}, nil
	case types.IsI64:
		return &PrimitiveFieldType[types.I64]{}, nil
	case types.IsI128:
		return &PrimitiveFieldType[types.I128]{}, nil
	case types.IsI256:
		return &PrimitiveFieldType[types.I256]{}, nil
	default:
		return nil, fmt.Errorf("unsupported primitive type %v", primitiveTypeDef)
	}
}

// getStoredFieldType will attempt to return a field type from storage, and perform an extra check for recursive types.
func (r *factory) getStoredFieldType(fieldName string, fieldType int64) (FieldType, bool) {
	if ft, ok := r.fieldStorage[fieldType]; ok {
		if rt, ok := ft.(*RecursiveFieldType); ok {
			r.recursiveFieldStorage[fieldType] = rt
		}

		return ft, ok
	}

	// Ensure that a recursive type such as Xcm::TransferReserveAsset does not cause an infinite loop
	// by adding the RecursiveFieldType the first time the field is encountered.
	r.fieldStorage[fieldType] = &RecursiveFieldType{
		FieldName: fieldName,
	}

	return nil, false
}

func (r *factory) storeFieldType(fieldType int64, registryFieldType FieldType) {
	r.fieldStorage[fieldType] = registryFieldType
}

const (
	unknownFieldName = "unknown_field_name"

	fieldPathSeparator = "::"
)

func getFieldPath(fieldType *types.Si1Type) string {
	var nameParts []string

	for _, pathEntry := range fieldType.Path {
		nameParts = append(nameParts, string(pathEntry))
	}

	return strings.Join(nameParts, fieldPathSeparator)
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
		return unknownFieldName
	}
}

// Type represents a parsed metadata type.
type Type struct {
	Name   string
	Fields []*Field
}

// Field represents one field of a Type.
type Field struct {
	Name        string
	FieldType   FieldType
	LookupIndex int64
}

func (f *Field) GetFieldMap() (map[string]any, error) {
	fieldType, err := f.FieldType.GetFieldTypeString()

	if err != nil {
		return nil, fmt.Errorf("couldn't get field type: %w", err)
	}

	fieldMap := map[string]any{
		"field_name":   f.Name,
		"field_type":   fieldType,
		"lookup_index": f.LookupIndex,
	}

	return fieldMap, nil
}

func (f *Field) MarshalJSON() ([]byte, error) {
	fieldMap, err := f.GetFieldMap()

	if err != nil {
		return nil, err
	}

	return json.Marshal(fieldMap)
}

// FieldType is the interface implemented by all the different types that are available.
type FieldType interface {
	GetFieldTypeString() (string, error)
}

// VariantFieldType represents an enum.
type VariantFieldType struct {
	FieldTypeMap map[byte]FieldType
}

func (v *VariantFieldType) GetFieldTypeString() (string, error) {
	return "enum", nil
}

// ArrayFieldType holds information about the length of the array and the type of its items.
type ArrayFieldType struct {
	Length   uint
	ItemType FieldType
}

func (a *ArrayFieldType) GetFieldTypeString() (string, error) {
	arrayItemTypeString, err := a.ItemType.GetFieldTypeString()

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("[%d]%s", a.Length, arrayItemTypeString), nil
}

// SliceFieldType represents a vector.
type SliceFieldType struct {
	ItemType FieldType
}

func (s *SliceFieldType) GetFieldTypeString() (string, error) {
	sliceItemTypeString, err := s.ItemType.GetFieldTypeString()

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("[]%s", sliceItemTypeString), nil
}

// CompositeFieldType represents a struct.
type CompositeFieldType struct {
	FieldName string
	Fields    []*Field
}

func (c *CompositeFieldType) GetFieldTypeString() (string, error) {
	return "struct", nil
}

// PrimitiveFieldType holds a primitive type.
type PrimitiveFieldType[T any] struct{}

func (f *PrimitiveFieldType[T]) GetFieldTypeString() (string, error) {
	var t T
	return fmt.Sprintf("%T", t), nil
}

// RecursiveFieldType is a wrapper for a FieldType that is recursive.
type RecursiveFieldType struct {
	depth int

	FieldName        string
	ResolvedItemType FieldType
}

func (r *RecursiveFieldType) GetFieldTypeString() (string, error) {
	if r.ResolvedItemType == nil {
		return "", fmt.Errorf("recursive type not resolved")
	}

	if r.depth > 0 {
		return "recursive", nil
	}

	r.depth++

	return r.ResolvedItemType.GetFieldTypeString()
}

// BitSequenceType represents a bit sequence.
type BitSequenceType struct {
	BitStoreType FieldType
	BitOrderType FieldType
}

func (b *BitSequenceType) GetFieldTypeString() (string, error) {
	bitStoreFieldType, err := b.BitStoreType.GetFieldTypeString()

	if err != nil {
		return "", err
	}

	bitOrderFieldType, err := b.BitOrderType.GetFieldTypeString()

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s.%s", bitStoreFieldType, bitOrderFieldType), nil
}
