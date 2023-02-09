package events

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

type RegistryFactory interface {
	CreateEventRegistry(meta *types.Metadata) (EventRegistry, error)
}

type EventRegistry map[types.EventID]*EventType

func (e EventRegistry) MarshalJSON() ([]byte, error) {
	eventRegistryMap := make(map[string]*EventType)

	for _, eventType := range e {
		eventRegistryMap[eventType.Name] = eventType
	}

	return json.Marshal(eventRegistryMap)
}

type registryFactory struct {
	fieldStorage          map[int64]EventFieldType
	recursiveFieldStorage map[int64]*RecursiveFieldType
}

func NewRegistryFactory() RegistryFactory {
	return &registryFactory{
		fieldStorage:          make(map[int64]EventFieldType),
		recursiveFieldStorage: make(map[int64]*RecursiveFieldType),
	}
}

func (r *registryFactory) CreateEventRegistry(meta *types.Metadata) (EventRegistry, error) {
	eventRegistry := make(map[types.EventID]*EventType)

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

			eventRegistry[eventID] = &EventType{
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

func (r *registryFactory) resolveRecursiveTypes() error {
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

func (r *registryFactory) getEventFields(meta *types.Metadata, fields []types.Si1Field) ([]*Field, error) {
	var eventFields []*Field

	for _, field := range fields {
		fieldType, ok := meta.AsMetadataV14.EfficientLookup[field.Type.Int64()]

		if !ok {
			return nil, fmt.Errorf("type not found for field '%s'", field.Name)
		}

		fieldName := getFieldName(field, fieldType)

		if eventFieldType, ok := r.getStoredEventFieldType(fieldName, field.Type.Int64()); ok {
			eventFields = append(eventFields, &Field{
				Name:        fieldName,
				FieldType:   eventFieldType,
				LookupIndex: field.Type.Int64(),
			})
			continue
		}

		fieldTypeDef := fieldType.Def

		eventFieldType, err := r.getEventFieldType(meta, fieldName, fieldTypeDef)

		if err != nil {
			return nil, fmt.Errorf("couldn't get event field type for '%s': %w", fieldName, err)
		}

		r.storeEventFieldType(field.Type.Int64(), eventFieldType)

		eventFields = append(eventFields, &Field{
			Name:        fieldName,
			FieldType:   eventFieldType,
			LookupIndex: field.Type.Int64(),
		})
	}

	return eventFields, nil
}

func (r *registryFactory) getEventFieldType(meta *types.Metadata, fieldName string, typeDef types.Si1TypeDef) (EventFieldType, error) {
	switch {
	case typeDef.IsCompact:
		compactFieldType, ok := meta.AsMetadataV14.EfficientLookup[typeDef.Compact.Type.Int64()]

		if !ok {
			return nil, errors.New("type not found for compact field")
		}

		return r.getCompactFieldType(meta, fieldName, compactFieldType.Def)
	case typeDef.IsComposite:
		compositeFieldType := &CompositeFieldType{}

		fields, err := r.getEventFields(meta, typeDef.Composite.Fields)

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
			return &FieldType[[]any]{}, nil
		}

		return r.getTupleType(meta, fieldName, typeDef.Tuple)
	default:
		return nil, errors.New("unsupported field type definition")
	}
}

func (r *registryFactory) getVariantFieldType(meta *types.Metadata, typeDef types.Si1TypeDef) (EventFieldType, error) {
	variantFieldType := &VariantFieldType{}

	fieldTypeMap := make(map[byte]EventFieldType)

	for _, variant := range typeDef.Variant.Variants {
		if len(variant.Fields) == 0 {
			fieldTypeMap[byte(variant.Index)] = &FieldType[byte]{}
			continue
		}

		compositeFieldType := &CompositeFieldType{}

		fields, err := r.getEventFields(meta, variant.Fields)

		if err != nil {
			return nil, fmt.Errorf("couldn't get field types for variant '%d': %w", variant.Index, err)
		}

		compositeFieldType.Fields = fields

		fieldTypeMap[byte(variant.Index)] = compositeFieldType
	}

	variantFieldType.FieldTypeMap = fieldTypeMap

	return variantFieldType, nil
}

func (r *registryFactory) getCompactFieldType(meta *types.Metadata, fieldName string, typeDef types.Si1TypeDef) (EventFieldType, error) {
	switch {
	case typeDef.IsPrimitive:
		return &FieldType[types.UCompact]{}, nil
	case typeDef.IsComposite:
		compactCompositeFields := typeDef.Composite.Fields

		compositeFieldType := &CompositeFieldType{}

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
				Name:        fieldName,
				FieldType:   compactCompositeType,
				LookupIndex: compactCompositeField.Type.Int64(),
			})
		}

		return compositeFieldType, nil
	default:
		return nil, errors.New("unsupported compact field type")
	}
}

func (r *registryFactory) getArrayFieldType(arrayLen uint, meta *types.Metadata, fieldName string, typeDef types.Si1TypeDef) (EventFieldType, error) {
	itemFieldType, err := r.getEventFieldType(meta, fieldName, typeDef)

	if err != nil {
		return nil, fmt.Errorf("couldn't get array item field type: %w", err)
	}

	arrayFieldType := &ArrayFieldType{Length: arrayLen, ItemType: itemFieldType}

	return arrayFieldType, nil
}

func (r *registryFactory) getSliceFieldType(meta *types.Metadata, fieldName string, typeDef types.Si1TypeDef) (EventFieldType, error) {
	itemFieldType, err := r.getEventFieldType(meta, fieldName, typeDef)

	if err != nil {
		return nil, fmt.Errorf("couldn't get slice item field type: %w", err)
	}

	sliceFieldType := &SliceFieldType{itemFieldType}

	return sliceFieldType, nil
}

func (r *registryFactory) getTupleType(meta *types.Metadata, fieldName string, tuple types.Si1TypeDefTuple) (EventFieldType, error) {
	compositeFieldType := &CompositeFieldType{}

	for _, item := range tuple {
		itemTypeDef, ok := meta.AsMetadataV14.EfficientLookup[item.Int64()]

		if !ok {
			return nil, fmt.Errorf("type definition for tuple item %d not found", item.Int64())
		}

		itemFieldType, err := r.getEventFieldType(meta, fieldName, itemTypeDef.Def)

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
}

func getPrimitiveType(primitiveTypeDef types.Si0TypeDefPrimitive) (EventFieldType, error) {
	switch primitiveTypeDef {
	case types.IsBool:
		return &FieldType[bool]{}, nil
	case types.IsChar:
		return &FieldType[byte]{}, nil
	case types.IsStr:
		return &FieldType[string]{}, nil
	case types.IsU8:
		return &FieldType[types.U8]{}, nil
	case types.IsU16:
		return &FieldType[types.U16]{}, nil
	case types.IsU32:
		return &FieldType[types.U32]{}, nil
	case types.IsU64:
		return &FieldType[types.U64]{}, nil
	case types.IsU128:
		return &FieldType[types.U128]{}, nil
	case types.IsU256:
		return &FieldType[types.U256]{}, nil
	case types.IsI8:
		return &FieldType[types.I8]{}, nil
	case types.IsI16:
		return &FieldType[types.I16]{}, nil
	case types.IsI32:
		return &FieldType[types.I32]{}, nil
	case types.IsI64:
		return &FieldType[types.I64]{}, nil
	case types.IsI128:
		return &FieldType[types.I128]{}, nil
	case types.IsI256:
		return &FieldType[types.I256]{}, nil
	default:
		return nil, fmt.Errorf("unsupported primitive type %v", primitiveTypeDef)
	}
}

func (r *registryFactory) getStoredEventFieldType(fieldName string, fieldType int64) (EventFieldType, bool) {
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

func (r *registryFactory) storeEventFieldType(fieldType int64, eventFieldType EventFieldType) {
	r.fieldStorage[fieldType] = eventFieldType
}

const (
	UnknownFieldName = "unknown_field_name"
)

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
		return UnknownFieldName
	}
}

type EventType struct {
	Name   string
	Fields []*Field
}

func (e *EventType) String() (string, error) {
	eventTypeMap := map[string][]*Field{
		e.Name: e.Fields,
	}

	res, err := json.Marshal(eventTypeMap)

	if err != nil {
		return "", err
	}

	return string(res), nil
}

type Field struct {
	Name        string
	FieldType   EventFieldType
	LookupIndex int64
}

func (f *Field) GetFieldMap() (map[string]any, error) {
	fieldType, err := f.FieldType.GetFieldType()

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

type EventFieldType interface {
	GetFieldType() (string, error)
}

type VariantFieldType struct {
	FieldTypeMap map[byte]EventFieldType
}

func (v *VariantFieldType) GetFieldType() (string, error) {
	return "enum", nil
}

type ArrayFieldType struct {
	Length   uint
	ItemType EventFieldType
}

func (a *ArrayFieldType) GetFieldType() (string, error) {
	arrayItemTypeString, err := a.ItemType.GetFieldType()

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("[%d]%s", a.Length, arrayItemTypeString), nil
}

type SliceFieldType struct {
	ItemType EventFieldType
}

func (s *SliceFieldType) GetFieldType() (string, error) {
	sliceItemTypeString, err := s.ItemType.GetFieldType()

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("[]%s", sliceItemTypeString), nil
}

type CompositeFieldType struct {
	Fields []*Field
}

func (c *CompositeFieldType) GetFieldType() (string, error) {
	return "struct", nil
}

type FieldType[T any] struct{}

func (f *FieldType[T]) GetFieldType() (string, error) {
	var t T
	return fmt.Sprintf("%T", t), nil
}

type RecursiveFieldType struct {
	depth int

	FieldName        string
	ResolvedItemType EventFieldType
}

func (r *RecursiveFieldType) GetFieldType() (string, error) {
	if r.ResolvedItemType == nil {
		return "", fmt.Errorf("recursive type not resolved")
	}

	if r.depth > 0 {
		return "recursive", nil
	}

	r.depth++

	return r.ResolvedItemType.GetFieldType()
}
