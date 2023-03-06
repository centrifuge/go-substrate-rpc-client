package registry

import (
	"fmt"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/registry/test"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/stretchr/testify/assert"
)

func TestCreateEventRegistry(t *testing.T) {
	var tests = []struct {
		Chain       string
		MetadataHex string
	}{
		{
			Chain:       "centrifuge",
			MetadataHex: test.CentrifugeMetadataHex,
		},
		{
			Chain:       "polkadot",
			MetadataHex: test.PolkadotMetadataHex,
		},
		{
			Chain:       "acala",
			MetadataHex: test.AcalaMetaHex,
		},
		{
			Chain:       "statemint",
			MetadataHex: test.StatemintMetaHex,
		},
		{
			Chain:       "moonbeam",
			MetadataHex: test.MoonbeamMetaHex,
		},
	}

	for _, test := range tests {
		t.Run(test.Chain, func(t *testing.T) {
			var meta types.Metadata

			err := codec.DecodeFromHex(test.MetadataHex, &meta)
			assert.NoError(t, err)

			t.Log("Metadata was decoded successfully")

			factory := NewFactory()

			reg, err := factory.CreateEventRegistry(&meta)
			assert.NoError(t, err)

			t.Log("Event registry was created successfully")

			testAsserter := newTestAsserter()

			for _, pallet := range meta.AsMetadataV14.Pallets {
				if !pallet.HasEvents {
					continue
				}

				eventsType, ok := meta.AsMetadataV14.EfficientLookup[pallet.Events.Type.Int64()]
				assert.True(t, ok, fmt.Sprintf("Events type %d not found", pallet.Events.Type.Int64()))

				assert.True(t, eventsType.Def.IsVariant, fmt.Sprintf("Events type %d not a variant", pallet.Events.Type.Int64()))

				for _, eventVariant := range eventsType.Def.Variant.Variants {
					eventID := types.EventID{byte(pallet.Index), byte(eventVariant.Index)}

					registryEventType, ok := reg[eventID]
					assert.True(t, ok, fmt.Sprintf("Event with ID %v not found in registry", eventID))

					testAsserter.assertRegistryItemContainsAllTypes(t, meta, registryEventType.Fields, eventVariant.Fields)
				}
			}
		})
	}
}

func TestCreateCallRegistry(t *testing.T) {
	var tests = []struct {
		Chain       string
		MetadataHex string
	}{
		{
			Chain:       "centrifuge",
			MetadataHex: test.CentrifugeMetadataHex,
		},
		{
			Chain:       "polkadot",
			MetadataHex: test.PolkadotMetadataHex,
		},
		{
			Chain:       "acala",
			MetadataHex: test.AcalaMetaHex,
		},
		{
			Chain:       "statemint",
			MetadataHex: test.StatemintMetaHex,
		},
		{
			Chain:       "moonbeam",
			MetadataHex: test.MoonbeamMetaHex,
		},
	}

	for _, test := range tests {
		t.Run(test.Chain, func(t *testing.T) {
			var meta types.Metadata

			err := codec.DecodeFromHex(test.MetadataHex, &meta)
			assert.NoError(t, err)

			t.Log("Metadata was decoded successfully")

			factory := NewFactory()

			reg, err := factory.CreateCallRegistry(&meta)
			assert.NoError(t, err)

			t.Log("Call registry was created successfully")

			testAsserter := newTestAsserter()

			for _, pallet := range meta.AsMetadataV14.Pallets {
				if !pallet.HasCalls {
					continue
				}

				callsType, ok := meta.AsMetadataV14.EfficientLookup[pallet.Calls.Type.Int64()]
				assert.True(t, ok, fmt.Sprintf("Calls type %d not found", pallet.Events.Type.Int64()))

				assert.True(t, callsType.Def.IsVariant, fmt.Sprintf("Calls type %d not a variant", pallet.Events.Type.Int64()))

				for _, callVariant := range callsType.Def.Variant.Variants {
					callName := fmt.Sprintf("%s.%s", pallet.Name, callVariant.Name)

					registryCallType, ok := reg[callName]
					assert.True(t, ok, fmt.Sprintf("Call '%s' not found in registry", callName))

					testAsserter.assertRegistryItemContainsAllTypes(t, meta, registryCallType.Fields, callVariant.Fields)
				}
			}
		})
	}
}

type testAsserter struct {
	recursiveTypeMap map[int64]struct{}
}

func newTestAsserter() *testAsserter {
	return &testAsserter{map[int64]struct{}{}}
}

func (a *testAsserter) assertRegistryItemContainsAllTypes(t *testing.T, meta types.Metadata, registryItemFields []*Field, metaItemFields []types.Si1Field) {
	for i, metaItemField := range metaItemFields {
		registryItemField := registryItemFields[i]
		registryItemFieldType := registryItemField.FieldType
		metaLookupIndex := metaItemField.Type.Int64()

		if _, ok := a.recursiveTypeMap[metaLookupIndex]; ok {
			continue
		}

		if metaLookupIndex != registryItemField.LookupIndex {
			t.Fatalf("Field lookup index mismatch for field with index %d", i)
		}

		fieldType, ok := meta.AsMetadataV14.EfficientLookup[metaLookupIndex]
		assert.True(t, ok, "field type for field with type %d not found", metaItemField.Type.Int64())

		a.assertRegistryItemFieldIsCorrect(t, meta, registryItemFieldType, fieldType)

		if _, ok := registryItemField.FieldType.(*RecursiveFieldType); ok {
			a.recursiveTypeMap[metaLookupIndex] = struct{}{}
		}
	}
}

func (a *testAsserter) assertRegistryItemFieldIsCorrect(t *testing.T, meta types.Metadata, registryItemFieldType FieldType, metaFieldType *types.Si1Type) {
	metaFieldTypeDef := metaFieldType.Def

	switch {
	case metaFieldTypeDef.IsComposite:
		compositeRegistryFieldType, ok := registryItemFieldType.(*CompositeFieldType)

		if !ok {
			_, isRecursive := registryItemFieldType.(*RecursiveFieldType)
			assert.True(t, isRecursive, "expected recursive field")

			return
		}

		a.assertRegistryItemContainsAllTypes(t, meta, compositeRegistryFieldType.Fields, metaFieldTypeDef.Composite.Fields)
	case metaFieldTypeDef.IsVariant:
		variantRegistryFieldType, ok := registryItemFieldType.(*VariantFieldType)

		if !ok {
			_, isRecursive := registryItemFieldType.(*RecursiveFieldType)
			assert.True(t, isRecursive, "expected variant or recursive field")
			return
		}

		for _, variant := range metaFieldTypeDef.Variant.Variants {
			registryVariant, ok := variantRegistryFieldType.FieldTypeMap[byte(variant.Index)]
			assert.True(t, ok, "expected registry variant")

			if len(variant.Fields) == 0 {
				_, ok = registryVariant.(*PrimitiveFieldType[byte])
				assert.True(t, ok, "expected byte field type")
				continue
			}

			compositeRegistryField, ok := registryVariant.(*CompositeFieldType)
			assert.True(t, ok, "expected composite field type")

			a.assertRegistryItemContainsAllTypes(t, meta, compositeRegistryField.Fields, variant.Fields)
		}
	case metaFieldTypeDef.IsSequence:
		sliceRegistryField, ok := registryItemFieldType.(*SliceFieldType)

		if !ok {
			_, isRecursive := registryItemFieldType.(*RecursiveFieldType)
			assert.True(t, isRecursive, "expected recursive field")

			return
		}

		sequenceFieldType, ok := meta.AsMetadataV14.EfficientLookup[metaFieldTypeDef.Sequence.Type.Int64()]
		assert.True(t, ok, "couldn't get sequence field type")

		a.assertRegistryItemFieldIsCorrect(t, meta, sliceRegistryField.ItemType, sequenceFieldType)
	case metaFieldTypeDef.IsArray:
		arrayRegistryField, ok := registryItemFieldType.(*ArrayFieldType)
		assert.True(t, ok, "expected array field type in registry")

		arrayFieldType, ok := meta.AsMetadataV14.EfficientLookup[metaFieldTypeDef.Array.Type.Int64()]
		assert.True(t, ok, "couldn't get array field type")

		a.assertRegistryItemFieldIsCorrect(t, meta, arrayRegistryField.ItemType, arrayFieldType)
	case metaFieldTypeDef.IsTuple:
		if metaFieldTypeDef.Tuple == nil {
			_, ok := registryItemFieldType.(*PrimitiveFieldType[[]any])
			assert.True(t, ok, "expected empty tuple field type")
			return
		}

		compositeRegistryFieldType, ok := registryItemFieldType.(*CompositeFieldType)

		if !ok {
			_, isRecursive := registryItemFieldType.(*RecursiveFieldType)
			assert.True(t, isRecursive, "expected composite or recursive field")
			return
		}

		for i, item := range metaFieldTypeDef.Tuple {
			itemTypeDef, ok := meta.AsMetadataV14.EfficientLookup[item.Int64()]
			assert.True(t, ok, "couldn't get tuple item field type")

			registryTupleItemFieldType := compositeRegistryFieldType.Fields[i].FieldType

			a.assertRegistryItemFieldIsCorrect(t, meta, registryTupleItemFieldType, itemTypeDef)
		}
	case metaFieldTypeDef.IsPrimitive:
		primitiveFieldType, err := getPrimitiveType(metaFieldTypeDef.Primitive.Si0TypeDefPrimitive)
		assert.NoError(t, err, "couldn't get primitive type")

		assert.Equal(t, primitiveFieldType, registryItemFieldType, "primitive field types should match")
	case metaFieldTypeDef.IsCompact:
		compactFieldType, ok := meta.AsMetadataV14.EfficientLookup[metaFieldTypeDef.Compact.Type.Int64()]
		assert.True(t, ok, "couldn't find compact field type")

		switch {
		case compactFieldType.Def.IsPrimitive:
			_, ok = registryItemFieldType.(*PrimitiveFieldType[types.UCompact])
			assert.True(t, ok, "expected compact field type in registry")
		case compactFieldType.Def.IsTuple:
			if metaFieldTypeDef.Tuple == nil {
				_, ok := registryItemFieldType.(*PrimitiveFieldType[any])
				assert.True(t, ok, "expected empty tuple field type")
				return
			}

			compositeRegistryField, ok := registryItemFieldType.(*CompositeFieldType)
			assert.True(t, ok, "expected composite field type in registry")

			for _, field := range compositeRegistryField.Fields {
				_, ok = field.FieldType.(*PrimitiveFieldType[types.UCompact])
				assert.True(t, ok, "expected compact field type in registry")
			}
		case compactFieldType.Def.IsComposite:
			compositeRegistryField, ok := registryItemFieldType.(*CompositeFieldType)
			assert.True(t, ok, "expected composite field type in registry")

			for _, field := range compositeRegistryField.Fields {
				_, ok = field.FieldType.(*PrimitiveFieldType[types.UCompact])
				assert.True(t, ok, "expected compact field type in registry")
			}
		default:
			t.Fatalf("unsupported compact field type")
		}
	case metaFieldTypeDef.IsBitSequence:
		bitSequenceType, ok := registryItemFieldType.(*BitSequenceType)
		assert.True(t, ok, "expected bit sequence field type in registry")

		bitStoreType, ok := meta.AsMetadataV14.EfficientLookup[metaFieldTypeDef.BitSequence.BitStoreType.Int64()]
		assert.True(t, ok, "couldn't get bit store field type")

		a.assertRegistryItemFieldIsCorrect(t, meta, bitSequenceType.BitStoreType, bitStoreType)

		bitOrderType, ok := meta.AsMetadataV14.EfficientLookup[metaFieldTypeDef.BitSequence.BitOrderType.Int64()]
		assert.True(t, ok, "couldn't get bit order field type")

		a.assertRegistryItemFieldIsCorrect(t, meta, bitSequenceType.BitOrderType, bitOrderType)
	case metaFieldTypeDef.IsHistoricMetaCompat:
		t.Fatalf("historic meta compat type not covered")
	}
}
