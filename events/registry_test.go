package events

import (
	"fmt"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/events/test"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/stretchr/testify/assert"
)

func TestCreateEventRegistry(t *testing.T) {
	var meta types.Metadata

	err := codec.DecodeFromHex(test.MetadataHex, &meta)
	assert.NoError(t, err)

	t.Log("Metadata was decoded successfully")

	reg, err := CreateEventRegistry(&meta)
	assert.NoError(t, err)

	t.Log("Event registry was created successfully")

	for _, pallet := range meta.AsMetadataV14.Pallets {
		if !pallet.HasEvents {
			continue
		}

		eventsType, ok := meta.AsMetadataV14.EfficientLookup[pallet.Events.Type.Int64()]
		assert.True(t, ok, fmt.Sprintf("Events type %d not found", pallet.Events.Type.Int64()))

		assert.True(t, eventsType.Def.IsVariant, fmt.Sprintf("Events type %d not a variant", pallet.Events.Type.Int64()))

		for _, eventVariant := range eventsType.Def.Variant.Variants {
			eventID := types.EventID{byte(pallet.Index), byte(eventVariant.Index)}

			eventType, ok := reg[eventID]
			assert.True(t, ok, fmt.Sprintf("Event with ID %v not found in registry", eventID))

			assertRegistryEventContainsAllTypes(t, meta, eventType.Fields, eventVariant.Fields)

			eventStr, err := eventType.String()
			assert.NoError(t, err)

			t.Logf("String representation for event with ID %v:\n%s", eventID, eventStr)
		}
	}
}

func assertRegistryEventContainsAllTypes(t *testing.T, meta types.Metadata, registryEventFields []*Field, metaEventFields []types.Si1Field) {
	for i, metaEventField := range metaEventFields {
		registryEventField := registryEventFields[i]

		if metaEventField.Type.Int64() != registryEventField.LookupIndex {
			t.Fatalf("Event field lookup index mismatch for field with index %d", i)
		}

		eventFieldType, ok := meta.AsMetadataV14.EfficientLookup[metaEventField.Type.Int64()]
		assert.True(t, ok, "event field type for event field with type %d not found", metaEventField.Type.Int64())

		assertRegistryEventFieldIsCorrect(t, meta, registryEventField.FieldType, eventFieldType)
	}
}

func assertRegistryEventFieldIsCorrect(t *testing.T, meta types.Metadata, registryEventFieldType EventFieldType, metaEventType *types.Si1Type) {
	metaEventFieldTypeDef := metaEventType.Def

	switch {
	case metaEventFieldTypeDef.IsComposite:
		compositeRegistryFieldType, ok := registryEventFieldType.(*CompositeFieldType)

		if !ok {
			_, isRecursive := registryEventFieldType.(*RecursiveFieldType)
			assert.True(t, isRecursive, "expected composite or recursive event field")
			return
		}

		assertRegistryEventContainsAllTypes(t, meta, compositeRegistryFieldType.Fields, metaEventFieldTypeDef.Composite.Fields)
	case metaEventFieldTypeDef.IsVariant:
		variantRegistryFieldType, ok := registryEventFieldType.(*VariantFieldType)
		assert.True(t, ok, "expected variant field type in registry")

		for _, variant := range metaEventFieldTypeDef.Variant.Variants {
			registryVariant, ok := variantRegistryFieldType.FieldTypeMap[byte(variant.Index)]
			assert.True(t, ok, "expected registry variant")

			if len(variant.Fields) == 0 {
				_, ok = registryVariant.(*FieldType[byte])
				assert.True(t, ok, "expected byte field type")
				continue
			}

			compositeRegistryEventField, ok := registryVariant.(*CompositeFieldType)
			assert.True(t, ok, "expected composite field type")

			assertRegistryEventContainsAllTypes(t, meta, compositeRegistryEventField.Fields, variant.Fields)
		}
	case metaEventFieldTypeDef.IsSequence:
		sliceRegistryEventField, ok := registryEventFieldType.(*SliceFieldType)
		assert.True(t, ok, "expected slice field type in registry")

		sequenceFieldType, ok := meta.AsMetadataV14.EfficientLookup[metaEventFieldTypeDef.Sequence.Type.Int64()]
		assert.True(t, ok, "couldn't get sequence field type")

		assertRegistryEventFieldIsCorrect(t, meta, sliceRegistryEventField.ItemType, sequenceFieldType)
	case metaEventFieldTypeDef.IsArray:
		arrayRegsitryEventField, ok := registryEventFieldType.(*ArrayFieldType)
		assert.True(t, ok, "expected array field type in registry")

		arrayFieldType, ok := meta.AsMetadataV14.EfficientLookup[metaEventFieldTypeDef.Array.Type.Int64()]
		assert.True(t, ok, "couldn't get array field type")

		assertRegistryEventFieldIsCorrect(t, meta, arrayRegsitryEventField.ItemType, arrayFieldType)
	case metaEventFieldTypeDef.IsTuple:
		if metaEventFieldTypeDef.Tuple == nil {
			_, ok := registryEventFieldType.(*FieldType[[]any])
			assert.True(t, ok, "expected empty tuple field type")
			return
		}

		compositeRegistryFieldType, ok := registryEventFieldType.(*CompositeFieldType)

		if !ok {
			_, isRecursive := registryEventFieldType.(*RecursiveFieldType)
			assert.True(t, isRecursive, "expected composite or recursive event field")
			return
		}

		for i, item := range metaEventFieldTypeDef.Tuple {
			itemTypeDef, ok := meta.AsMetadataV14.EfficientLookup[item.Int64()]
			assert.True(t, ok, "couldn't get tuple item field type")

			registryTupleItemFieldType := compositeRegistryFieldType.Fields[i].FieldType

			assertRegistryEventFieldIsCorrect(t, meta, registryTupleItemFieldType, itemTypeDef)
		}
	case metaEventFieldTypeDef.IsPrimitive:
		primitiveEventFieldType, err := getPrimitiveType(metaEventFieldTypeDef.Primitive.Si0TypeDefPrimitive)
		assert.NoError(t, err, "couldn't get primitive type")

		assert.Equal(t, primitiveEventFieldType, registryEventFieldType, "primitive field types should match")
	case metaEventFieldTypeDef.IsCompact:
		compactFieldType, ok := meta.AsMetadataV14.EfficientLookup[metaEventFieldTypeDef.Compact.Type.Int64()]
		assert.True(t, ok, "couldn't find compact field type")

		switch {
		case compactFieldType.Def.IsPrimitive:
			_, ok = registryEventFieldType.(*FieldType[types.UCompact])
			assert.True(t, ok, "expected compact field type in registry")
		case compactFieldType.Def.IsComposite:
			compositeRegistryEventField, ok := registryEventFieldType.(*CompositeFieldType)
			assert.True(t, ok, "expected composite field type in registry")

			for _, field := range compositeRegistryEventField.Fields {
				_, ok = field.FieldType.(*FieldType[types.UCompact])
				assert.True(t, ok, "expected compact field type in registry")
			}
		default:
			t.Fatalf("unsupported compact event field type")
		}
	case metaEventFieldTypeDef.IsBitSequence:
		t.Fatalf("bit sequence type not covered")
	case metaEventFieldTypeDef.IsHistoricMetaCompat:
		t.Fatalf("historic meta compat type not covered")
	}
}
