package registry

import (
	"fmt"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/registry/test"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/stretchr/testify/assert"
)

func TestFactory_CreateErrorRegistryWithLiveMetadata(t *testing.T) {
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

			reg, err := factory.CreateErrorRegistry(&meta)
			assert.NoError(t, err)

			t.Log("Error registry was created successfully")

			testAsserter := newTestAsserter()

			for _, pallet := range meta.AsMetadataV14.Pallets {
				if !pallet.HasErrors {
					continue
				}

				errorsType, ok := meta.AsMetadataV14.EfficientLookup[pallet.Errors.Type.Int64()]
				assert.True(t, ok, fmt.Sprintf("Error type %d not found", pallet.Errors.Type.Int64()))

				assert.True(t, errorsType.Def.IsVariant, fmt.Sprintf("Error type %d not a variant", pallet.Events.Type.Int64()))

				for _, errorVariant := range errorsType.Def.Variant.Variants {
					errorName := fmt.Sprintf("%s.%s", pallet.Name, errorVariant.Name)

					registryErrorType, ok := reg[errorName]
					assert.True(t, ok, fmt.Sprintf("Error '%s' not found in registry", errorName))

					testAsserter.assertRegistryItemContainsAllTypes(t, meta, registryErrorType.Fields, errorVariant.Fields)
				}
			}
		})
	}
}

func TestFactory_CreateErrorRegistry_NoPalletWithErrors(t *testing.T) {
	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			Pallets: []types.PalletMetadataV14{
				{
					HasErrors: false,
				},
			},
		},
	}

	factory := NewFactory()

	reg, err := factory.CreateErrorRegistry(testMeta)
	assert.NoError(t, err)
	assert.Empty(t, reg)
}

func TestFactory_CreateErrorRegistry_ErrorsTypeNotFound(t *testing.T) {
	testModuleName := "TestModule"
	errorLookupTypeID := 123

	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			Pallets: []types.PalletMetadataV14{
				{
					Name:      "TestModule",
					HasErrors: true,
					Errors: types.ErrorMetadataV14{
						Type: types.Si1LookupTypeID{
							UCompact: types.NewUCompactFromUInt(uint64(errorLookupTypeID)),
						},
					},
				},
			},
			// EfficientLookup map is empty causing an error.
			EfficientLookup: map[int64]*types.Si1Type{},
		},
	}

	factory := NewFactory()

	reg, err := factory.CreateErrorRegistry(testMeta)
	assert.Equal(t, fmt.Sprintf("errors type %d not found for module '%s'", errorLookupTypeID, testModuleName), err.Error())
	assert.Empty(t, reg)
}

func TestFactory_CreateErrorRegistry_ErrorsTypeNotAVariant(t *testing.T) {
	testModuleName := "TestModule"
	errorLookupTypeID := 123

	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			Pallets: []types.PalletMetadataV14{
				{
					Name:      "TestModule",
					HasErrors: true,
					Errors: types.ErrorMetadataV14{
						Type: types.Si1LookupTypeID{
							UCompact: types.NewUCompactFromUInt(uint64(errorLookupTypeID)),
						},
					},
				},
			},
			EfficientLookup: map[int64]*types.Si1Type{
				int64(errorLookupTypeID): {
					Def: types.Si1TypeDef{
						// Error type definition not a variant causing an error.
						IsVariant: false,
					},
				},
			},
		},
	}

	factory := NewFactory()

	reg, err := factory.CreateErrorRegistry(testMeta)
	assert.Equal(t, fmt.Sprintf("errors type %d for module '%s' is not a variant", errorLookupTypeID, testModuleName), err.Error())
	assert.Empty(t, reg)
}

func TestFactory_CreateErrorRegistry_GetTypeFieldsError(t *testing.T) {
	errorLookupTypeID := 123

	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			Pallets: []types.PalletMetadataV14{
				{
					Name:      "TestModule",
					HasErrors: true,
					Errors: types.ErrorMetadataV14{
						Type: types.Si1LookupTypeID{
							UCompact: types.NewUCompactFromUInt(uint64(errorLookupTypeID)),
						},
					},
				},
			},
			EfficientLookup: map[int64]*types.Si1Type{
				int64(errorLookupTypeID): {
					Def: types.Si1TypeDef{
						IsVariant: true,
						Variant: types.Si1TypeDefVariant{
							Variants: []types.Si1Variant{
								{
									Name: "ErrorVariant1",
									Fields: []types.Si1Field{
										{
											HasName: true,
											Name:    "ErrorVariant1Field",
											Type: types.Si1LookupTypeID{
												// This lookup type ID is not added in the lookup map which should
												// cause an error.
												UCompact: types.NewUCompactFromUInt(uint64(456)),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	factory := NewFactory()

	reg, err := factory.CreateErrorRegistry(testMeta)
	assert.Equal(t, "couldn't get fields for error 'TestModule.ErrorVariant1': type not found for field 'ErrorVariant1Field'", err.Error())
	assert.Empty(t, reg)
}

func TestFactory_CreateCallRegistryWithLiveMetadata(t *testing.T) {
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
				assert.True(t, ok, fmt.Sprintf("Calls type %d not found", pallet.Calls.Type.Int64()))

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

func TestFactory_CreateCallRegistry_NoPalletWithCalls(t *testing.T) {
	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			Pallets: []types.PalletMetadataV14{
				{
					HasCalls: false,
				},
			},
		},
	}

	factory := NewFactory()

	reg, err := factory.CreateCallRegistry(testMeta)
	assert.NoError(t, err)
	assert.Empty(t, reg)
}

func TestFactory_CreateCallRegistry_CallsTypeNotFound(t *testing.T) {
	testModuleName := "TestModule"
	callLookupTypeID := 123

	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			Pallets: []types.PalletMetadataV14{
				{
					Name:     "TestModule",
					HasCalls: true,
					Calls: types.FunctionMetadataV14{
						Type: types.Si1LookupTypeID{
							UCompact: types.NewUCompactFromUInt(uint64(callLookupTypeID)),
						},
					},
				},
			},
			// EfficientLookup map is empty causing an error.
			EfficientLookup: map[int64]*types.Si1Type{},
		},
	}

	factory := NewFactory()

	reg, err := factory.CreateCallRegistry(testMeta)
	assert.Equal(t, fmt.Sprintf("calls type %d not found for module '%s'", callLookupTypeID, testModuleName), err.Error())
	assert.Empty(t, reg)
}

func TestFactory_CreateCallRegistry_CallTypeNotAVariant(t *testing.T) {
	testModuleName := "TestModule"
	callLookupTypeID := 123

	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			Pallets: []types.PalletMetadataV14{
				{
					Name:     "TestModule",
					HasCalls: true,
					Calls: types.FunctionMetadataV14{
						Type: types.Si1LookupTypeID{
							UCompact: types.NewUCompactFromUInt(uint64(callLookupTypeID)),
						},
					},
				},
			},
			EfficientLookup: map[int64]*types.Si1Type{
				int64(callLookupTypeID): {
					Def: types.Si1TypeDef{
						// Calls type definition not a variant causing an error.
						IsVariant: false,
					},
				},
			},
		},
	}

	factory := NewFactory()

	reg, err := factory.CreateCallRegistry(testMeta)
	assert.Equal(t, fmt.Sprintf("calls type %d for module '%s' is not a variant", callLookupTypeID, testModuleName), err.Error())
	assert.Empty(t, reg)
}

func TestFactory_CreateCallRegistry_GetTypeFieldsError(t *testing.T) {
	callLookupTypeID := 123

	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			Pallets: []types.PalletMetadataV14{
				{
					Name:     "TestModule",
					HasCalls: true,
					Calls: types.FunctionMetadataV14{
						Type: types.Si1LookupTypeID{
							UCompact: types.NewUCompactFromUInt(uint64(callLookupTypeID)),
						},
					},
				},
			},
			EfficientLookup: map[int64]*types.Si1Type{
				int64(callLookupTypeID): {
					Def: types.Si1TypeDef{
						IsVariant: true,
						Variant: types.Si1TypeDefVariant{
							Variants: []types.Si1Variant{
								{
									Name: "CallVariant1",
									Fields: []types.Si1Field{
										{
											HasName: true,
											Name:    "CallVariant1Field",
											Type: types.Si1LookupTypeID{
												// This lookup type ID is not added in the lookup map which should
												// cause an error.
												UCompact: types.NewUCompactFromUInt(uint64(456)),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	factory := NewFactory()

	reg, err := factory.CreateCallRegistry(testMeta)
	assert.Equal(t, "couldn't get fields for call 'TestModule.CallVariant1': type not found for field 'CallVariant1Field'", err.Error())
	assert.Empty(t, reg)
}

func TestFactory_CreateEventRegistryWithLiveMetadata(t *testing.T) {
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

func TestFactory_CreateEventRegistry_NoPalletWithEvents(t *testing.T) {
	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			Pallets: []types.PalletMetadataV14{
				{
					HasEvents: false,
				},
			},
		},
	}

	factory := NewFactory()

	reg, err := factory.CreateEventRegistry(testMeta)
	assert.NoError(t, err)
	assert.Empty(t, reg)
}

func TestFactory_CreateEventRegistry_EventsTypeNotFound(t *testing.T) {
	testModuleName := "TestModule"
	eventLookupTypeID := 123

	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			Pallets: []types.PalletMetadataV14{
				{
					Name:      "TestModule",
					HasEvents: true,
					Events: types.EventMetadataV14{
						Type: types.Si1LookupTypeID{
							UCompact: types.NewUCompactFromUInt(uint64(eventLookupTypeID)),
						},
					},
				},
			},
			// EfficientLookup map is empty causing an error.
			EfficientLookup: map[int64]*types.Si1Type{},
		},
	}

	factory := NewFactory()

	reg, err := factory.CreateEventRegistry(testMeta)
	assert.Equal(t, fmt.Sprintf("events type %d not found for module '%s'", eventLookupTypeID, testModuleName), err.Error())
	assert.Empty(t, reg)
}

func TestFactory_CreateEventRegistry_EventTypeNotAVariant(t *testing.T) {
	testModuleName := "TestModule"
	callLookupTypeID := 123

	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			Pallets: []types.PalletMetadataV14{
				{
					Name:      "TestModule",
					HasEvents: true,
					Events: types.EventMetadataV14{
						Type: types.Si1LookupTypeID{
							UCompact: types.NewUCompactFromUInt(uint64(callLookupTypeID)),
						},
					},
				},
			},
			EfficientLookup: map[int64]*types.Si1Type{
				int64(callLookupTypeID): {
					Def: types.Si1TypeDef{
						// Events type definition not a variant causing an error.
						IsVariant: false,
					},
				},
			},
		},
	}

	factory := NewFactory()

	reg, err := factory.CreateEventRegistry(testMeta)
	assert.Equal(t, fmt.Sprintf("events type %d for module '%s' is not a variant", callLookupTypeID, testModuleName), err.Error())
	assert.Empty(t, reg)
}

func TestFactory_CreateEventRegistry_GetTypeFieldError(t *testing.T) {
	callLookupTypeID := 123

	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			Pallets: []types.PalletMetadataV14{
				{
					Name:      "TestModule",
					HasEvents: true,
					Events: types.EventMetadataV14{
						Type: types.Si1LookupTypeID{
							UCompact: types.NewUCompactFromUInt(uint64(callLookupTypeID)),
						},
					},
				},
			},
			EfficientLookup: map[int64]*types.Si1Type{
				int64(callLookupTypeID): {
					Def: types.Si1TypeDef{
						IsVariant: true,
						Variant: types.Si1TypeDefVariant{
							Variants: []types.Si1Variant{
								{
									Name: "EventVariant1",
									Fields: []types.Si1Field{
										{
											HasName: true,
											Name:    "EventVariant1Field",
											Type: types.Si1LookupTypeID{
												// This lookup type ID is not added in the lookup map which should
												// cause an error.
												UCompact: types.NewUCompactFromUInt(uint64(456)),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	factory := NewFactory()

	reg, err := factory.CreateEventRegistry(testMeta)
	assert.Equal(t, "couldn't get fields for event 'TestModule.EventVariant1': type not found for field 'EventVariant1Field'", err.Error())
	assert.Empty(t, reg)
}

func TestFactory_getTypeFields(t *testing.T) {
	fieldLookUpID := 123

	testFieldName := "TestFieldName"
	testFields := []types.Si1Field{
		{
			HasName: true,
			Name:    types.Text(testFieldName),
			Type: types.Si1LookupTypeID{
				UCompact: types.NewUCompactFromUInt(uint64(fieldLookUpID)),
			},
		},
	}

	compactFieldTypeLookupID := 456

	testFieldTypeDef := types.Si1TypeDef{
		IsCompact: true,
		Compact: types.Si1TypeDefCompact{
			Type: types.Si1LookupTypeID{
				UCompact: types.NewUCompactFromUInt(uint64(compactFieldTypeLookupID)),
			},
		},
	}

	compactFieldTypeDef := types.Si1TypeDef{
		IsPrimitive: true,
	}

	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			EfficientLookup: map[int64]*types.Si1Type{
				int64(fieldLookUpID): {
					Def: testFieldTypeDef,
				},
				int64(compactFieldTypeLookupID): {
					Def: compactFieldTypeDef,
				},
			},
		},
	}

	factory := NewFactory().(*factory)

	res, err := factory.getTypeFields(testMeta, testFields)
	assert.NoError(t, err)
	assert.Len(t, res, 1)

	assert.Equal(t, testFieldName, res[0].Name)
	assert.Equal(t, &PrimitiveFieldType[types.UCompact]{}, res[0].FieldType)
	assert.Equal(t, int64(fieldLookUpID), res[0].LookupIndex)
}

func TestFactory_getTypeFields_FieldTypeError(t *testing.T) {
	fieldLookUpID := 123

	testFieldName := "TestFieldName"
	testFields := []types.Si1Field{
		{
			HasName: true,
			Name:    types.Text(testFieldName),
			Type: types.Si1LookupTypeID{
				UCompact: types.NewUCompactFromUInt(uint64(fieldLookUpID)),
			},
		},
	}

	compositeFieldTypeLookupID := 456

	testFieldTypeDef := types.Si1TypeDef{
		IsComposite: true,
		Composite: types.Si1TypeDefComposite{
			Fields: []types.Si1Field{
				{
					HasName: true,
					Name:    "CompositeField1",
					Type: types.Si1LookupTypeID{
						// This lookup ID is not in the efficient lookup map which should cause an error.
						UCompact: types.NewUCompactFromUInt(uint64(compositeFieldTypeLookupID)),
					},
				},
			},
		},
	}

	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			EfficientLookup: map[int64]*types.Si1Type{
				int64(fieldLookUpID): {
					Def: testFieldTypeDef,
				},
			},
		},
	}

	factory := NewFactory().(*factory)

	res, err := factory.getTypeFields(testMeta, testFields)
	assert.Equal(t, "couldn't get field type for 'TestFieldName': couldn't get composite fields: type not found for field 'CompositeField1'", err.Error())
	assert.Nil(t, res)
}

func TestFactory_getTypeFields_FieldTypeNotFoundError(t *testing.T) {
	fieldLookUpID := 123

	testFieldName := "TestFieldName"
	testFields := []types.Si1Field{
		{
			HasName: true,
			Name:    types.Text(testFieldName),
			Type: types.Si1LookupTypeID{
				UCompact: types.NewUCompactFromUInt(uint64(fieldLookUpID)),
			},
		},
	}

	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			// Efficient lookup map is empty to ensure that an error is returned.
			EfficientLookup: map[int64]*types.Si1Type{},
		},
	}

	factory := NewFactory().(*factory)

	res, err := factory.getTypeFields(testMeta, testFields)
	assert.Equal(t, fmt.Sprintf("type not found for field '%s'", testFieldName), err.Error())
	assert.Nil(t, res)
}

func TestFactory_getFieldType_UnsupportedTypeError(t *testing.T) {
	testFieldName := "TestFieldName"

	testFieldTypeDef := types.Si1TypeDef{
		IsHistoricMetaCompat: true,
	}

	testMeta := &types.Metadata{}

	factory := NewFactory().(*factory)

	res, err := factory.getFieldType(testMeta, testFieldName, testFieldTypeDef)
	assert.Equal(t, "unsupported field type definition", err.Error())
	assert.Nil(t, res)
}

func TestFactory_getFieldType_Compact(t *testing.T) {
	testFieldName := "TestFieldName"

	compactFieldTypeLookupID := 456

	testFieldTypeDef := types.Si1TypeDef{
		IsCompact: true,
		Compact: types.Si1TypeDefCompact{
			Type: types.Si1LookupTypeID{
				UCompact: types.NewUCompactFromUInt(uint64(compactFieldTypeLookupID)),
			},
		},
	}

	compactFieldTypeDef := types.Si1TypeDef{
		IsPrimitive: true,
	}

	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			EfficientLookup: map[int64]*types.Si1Type{
				int64(compactFieldTypeLookupID): {
					Def: compactFieldTypeDef,
				},
			},
		},
	}

	factory := NewFactory().(*factory)

	res, err := factory.getFieldType(testMeta, testFieldName, testFieldTypeDef)
	assert.NoError(t, err)
	assert.Equal(t, &PrimitiveFieldType[types.UCompact]{}, res)
}

func TestFactory_getFieldType_Compact_TypeNotFoundError(t *testing.T) {
	testFieldName := "TestFieldName"

	compactFieldTypeLookupID := 456

	testFieldTypeDef := types.Si1TypeDef{
		IsCompact: true,
		Compact: types.Si1TypeDefCompact{
			Type: types.Si1LookupTypeID{
				UCompact: types.NewUCompactFromUInt(uint64(compactFieldTypeLookupID)),
			},
		},
	}

	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			EfficientLookup: map[int64]*types.Si1Type{},
		},
	}

	factory := NewFactory().(*factory)

	res, err := factory.getFieldType(testMeta, testFieldName, testFieldTypeDef)
	assert.Equal(t, "type not found for compact field", err.Error())
	assert.Nil(t, res)
}

func TestFactory_getFieldType_Composite(t *testing.T) {
	testFieldName := "TestFieldName"

	compositeFieldTypeLookupID1 := 123
	compositeFieldTypeLookupID2 := 456

	compositeFieldName1 := "CompositeFieldName1"
	compositeFieldName2 := "CompositeFieldName2"

	compositeFields := []types.Si1Field{
		{
			HasName: true,
			Name:    types.Text(compositeFieldName1),
			Type: types.Si1LookupTypeID{
				UCompact: types.NewUCompactFromUInt(uint64(compositeFieldTypeLookupID1)),
			},
		},
		{
			HasName: true,
			Name:    types.Text(compositeFieldName2),
			Type: types.Si1LookupTypeID{
				UCompact: types.NewUCompactFromUInt(uint64(compositeFieldTypeLookupID2)),
			},
		},
	}
	testFieldTypeDef := types.Si1TypeDef{
		IsComposite: true,
		Composite: types.Si1TypeDefComposite{
			Fields: compositeFields,
		},
	}

	compositeFieldTypeDef1 := types.Si1TypeDef{
		IsPrimitive: true,
		Primitive: types.Si1TypeDefPrimitive{
			Si0TypeDefPrimitive: types.Si0TypeDefPrimitive(types.IsU8),
		},
	}

	compositeFieldTypeDef2 := types.Si1TypeDef{
		IsPrimitive: true,
		Primitive: types.Si1TypeDefPrimitive{
			Si0TypeDefPrimitive: types.Si0TypeDefPrimitive(types.IsI8),
		},
	}

	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			EfficientLookup: map[int64]*types.Si1Type{
				int64(compositeFieldTypeLookupID1): {
					Def: compositeFieldTypeDef1,
				},
				int64(compositeFieldTypeLookupID2): {
					Def: compositeFieldTypeDef2,
				},
			},
		},
	}

	factory := NewFactory().(*factory)

	res, err := factory.getFieldType(testMeta, testFieldName, testFieldTypeDef)
	assert.NoError(t, err)

	compositeFieldType, ok := res.(*CompositeFieldType)
	assert.True(t, ok)
	assert.Len(t, compositeFieldType.Fields, 2)

	assert.Equal(t, &PrimitiveFieldType[types.U8]{}, compositeFieldType.Fields[0].FieldType)
	assert.Equal(t, compositeFieldName1, compositeFieldType.Fields[0].Name)
	assert.Equal(t, int64(compositeFieldTypeLookupID1), compositeFieldType.Fields[0].LookupIndex)

	assert.Equal(t, &PrimitiveFieldType[types.I8]{}, compositeFieldType.Fields[1].FieldType)
	assert.Equal(t, compositeFieldName2, compositeFieldType.Fields[1].Name)
	assert.Equal(t, int64(compositeFieldTypeLookupID2), compositeFieldType.Fields[1].LookupIndex)
}

func TestFactory_getFieldType_Composite_FieldError(t *testing.T) {
	testFieldName := "TestFieldName"

	compositeFieldTypeLookupID1 := 123
	compositeFieldTypeLookupID2 := 456

	compositeFieldName1 := "CompositeFieldName1"
	compositeFieldName2 := "CompositeFieldName2"

	compositeFields := []types.Si1Field{
		{
			HasName: true,
			Name:    types.Text(compositeFieldName1),
			Type: types.Si1LookupTypeID{
				UCompact: types.NewUCompactFromUInt(uint64(compositeFieldTypeLookupID1)),
			},
		},
		{
			HasName: true,
			Name:    types.Text(compositeFieldName2),
			Type: types.Si1LookupTypeID{
				UCompact: types.NewUCompactFromUInt(uint64(compositeFieldTypeLookupID2)),
			},
		},
	}
	testFieldTypeDef := types.Si1TypeDef{
		IsComposite: true,
		Composite: types.Si1TypeDefComposite{
			Fields: compositeFields,
		},
	}

	compositeFieldTypeDef1 := types.Si1TypeDef{
		IsPrimitive: true,
		Primitive: types.Si1TypeDefPrimitive{
			Si0TypeDefPrimitive: types.Si0TypeDefPrimitive(types.IsU8),
		},
	}

	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			EfficientLookup: map[int64]*types.Si1Type{
				int64(compositeFieldTypeLookupID1): {
					Def: compositeFieldTypeDef1,
				},
				// Omitting the entry for composite field 2 should cause and error when parsing the composite fields.
			},
		},
	}

	factory := NewFactory().(*factory)

	res, err := factory.getFieldType(testMeta, testFieldName, testFieldTypeDef)
	assert.Equal(t, fmt.Sprintf("couldn't get composite fields: type not found for field '%s'", compositeFieldName2), err.Error())
	assert.Nil(t, res)
}

func TestFactory_getFieldType_Variant(t *testing.T) {
	testFieldName := "TestField"

	variantName1 := "Variant1"
	variantName2 := "Variant2"

	variantFieldName := "VariantFieldName"
	variantFieldLookupID := 123
	variantFieldTypeDef := types.Si1TypeDef{
		IsPrimitive: true,
		Primitive: types.Si1TypeDefPrimitive{
			Si0TypeDefPrimitive: types.Si0TypeDefPrimitive(types.IsU8),
		},
	}

	testFieldTypeDef := types.Si1TypeDef{
		IsVariant: true,
		Variant: types.Si1TypeDefVariant{
			Variants: []types.Si1Variant{
				{
					Name:   types.Text(variantName1),
					Fields: nil,
					Index:  0,
				},
				{
					Name: types.Text(variantName2),
					Fields: []types.Si1Field{
						{
							HasName: true,
							Name:    types.Text(variantFieldName),
							Type: types.Si1LookupTypeID{
								UCompact: types.NewUCompactFromUInt(uint64(variantFieldLookupID)),
							},
						},
					},
					Index: 1,
				},
			},
		},
	}

	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			EfficientLookup: map[int64]*types.Si1Type{
				int64(variantFieldLookupID): {
					Def: variantFieldTypeDef,
				},
			},
		},
	}

	factory := NewFactory().(*factory)

	res, err := factory.getFieldType(testMeta, testFieldName, testFieldTypeDef)
	assert.NoError(t, err)

	variantFieldType, ok := res.(*VariantFieldType)
	assert.True(t, ok)
	assert.Len(t, variantFieldType.FieldTypeMap, 2)

	assert.Equal(t, &PrimitiveFieldType[byte]{}, variantFieldType.FieldTypeMap[0])

	compositeVariant, ok := variantFieldType.FieldTypeMap[1].(*CompositeFieldType)
	assert.True(t, ok)
	assert.Len(t, compositeVariant.Fields, 1)

	assert.Equal(t, variantFieldName, compositeVariant.Fields[0].Name)
	assert.Equal(t, &PrimitiveFieldType[types.U8]{}, compositeVariant.Fields[0].FieldType)
	assert.Equal(t, int64(variantFieldLookupID), compositeVariant.Fields[0].LookupIndex)
}

func TestFactory_getFieldType_Primitive(t *testing.T) {
	testFieldName := "TestFieldName"

	testFieldTypeDef := types.Si1TypeDef{
		IsPrimitive: true,
		Primitive: types.Si1TypeDefPrimitive{
			Si0TypeDefPrimitive: types.Si0TypeDefPrimitive(types.IsU8),
		},
	}

	testMeta := &types.Metadata{}

	factory := NewFactory().(*factory)

	res, err := factory.getFieldType(testMeta, testFieldName, testFieldTypeDef)
	assert.NoError(t, err)
	assert.Equal(t, &PrimitiveFieldType[types.U8]{}, res)
}

func TestFactory_getFieldType_Array(t *testing.T) {
	testFieldName := "TestFieldName"

	arrayItemTypeLookupID := 456
	arrayLen := 32

	testFieldTypeDef := types.Si1TypeDef{
		IsArray: true,
		Array: types.Si1TypeDefArray{
			Len: types.U32(arrayLen),
			Type: types.Si1LookupTypeID{
				UCompact: types.NewUCompactFromUInt(uint64(arrayItemTypeLookupID)),
			},
		},
	}

	arrayItemTypeDef := types.Si1TypeDef{
		IsPrimitive: true,
		Primitive: types.Si1TypeDefPrimitive{
			Si0TypeDefPrimitive: types.Si0TypeDefPrimitive(types.IsU8),
		},
	}

	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			EfficientLookup: map[int64]*types.Si1Type{
				int64(arrayItemTypeLookupID): {
					Def: arrayItemTypeDef,
				},
			},
		},
	}

	factory := NewFactory().(*factory)

	res, err := factory.getFieldType(testMeta, testFieldName, testFieldTypeDef)
	assert.NoError(t, err)

	arrayFieldType, ok := res.(*ArrayFieldType)
	assert.True(t, ok)

	assert.Equal(t, uint(arrayLen), arrayFieldType.Length)
	assert.Equal(t, &PrimitiveFieldType[types.U8]{}, arrayFieldType.ItemType)
}

func TestFactory_getFieldType_Array_TypeNotFoundError(t *testing.T) {
	testFieldName := "TestFieldName"

	arrayItemTypeLookupID := 456
	arrayLen := 32

	testFieldTypeDef := types.Si1TypeDef{
		IsArray: true,
		Array: types.Si1TypeDefArray{
			Len: types.U32(arrayLen),
			Type: types.Si1LookupTypeID{
				UCompact: types.NewUCompactFromUInt(uint64(arrayItemTypeLookupID)),
			},
		},
	}

	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			// The lookup map does not contain the array item type lookup ID which should cause an error.
			EfficientLookup: map[int64]*types.Si1Type{},
		},
	}

	factory := NewFactory().(*factory)

	res, err := factory.getFieldType(testMeta, testFieldName, testFieldTypeDef)
	assert.Equal(t, "type not found for array field", err.Error())
	assert.Nil(t, res)
}

func TestFactory_getFieldType_Slice(t *testing.T) {
	testFieldName := "TestFieldName"

	sliceItemTypeLookupID := 456

	testFieldTypeDef := types.Si1TypeDef{
		IsSequence: true,
		Sequence: types.Si1TypeDefSequence{
			Type: types.Si1LookupTypeID{
				UCompact: types.NewUCompactFromUInt(uint64(sliceItemTypeLookupID)),
			},
		},
	}

	sliceItemTypeDef := types.Si1TypeDef{
		IsPrimitive: true,
		Primitive: types.Si1TypeDefPrimitive{
			Si0TypeDefPrimitive: types.Si0TypeDefPrimitive(types.IsU256),
		},
	}

	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			EfficientLookup: map[int64]*types.Si1Type{
				int64(sliceItemTypeLookupID): {
					Def: sliceItemTypeDef,
				},
			},
		},
	}

	factory := NewFactory().(*factory)

	res, err := factory.getFieldType(testMeta, testFieldName, testFieldTypeDef)
	assert.NoError(t, err)

	sliceFieldType, ok := res.(*SliceFieldType)
	assert.True(t, ok)

	assert.Equal(t, &PrimitiveFieldType[types.U256]{}, sliceFieldType.ItemType)
}

func TestFactory_getFieldType_Slice_TypeNotFoundError(t *testing.T) {
	testFieldName := "TestFieldName"

	sliceItemTypeLookupID := 456

	testFieldTypeDef := types.Si1TypeDef{
		IsSequence: true,
		Sequence: types.Si1TypeDefSequence{
			Type: types.Si1LookupTypeID{
				UCompact: types.NewUCompactFromUInt(uint64(sliceItemTypeLookupID)),
			},
		},
	}

	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			// The lookup map does not contain the array item type lookup ID which should cause an error.
			EfficientLookup: map[int64]*types.Si1Type{},
		},
	}

	factory := NewFactory().(*factory)

	res, err := factory.getFieldType(testMeta, testFieldName, testFieldTypeDef)
	assert.Equal(t, "type not found for vector field", err.Error())
	assert.Nil(t, res)
}

func TestFactory_getFieldType_Tuple(t *testing.T) {
	testFieldName := "TestFieldName"

	tupleItemLookupID1 := 123
	tupleItemLookupID2 := 456

	testFieldTypeDef := types.Si1TypeDef{
		IsTuple: true,
		Tuple: []types.Si1LookupTypeID{
			{
				UCompact: types.NewUCompactFromUInt(uint64(tupleItemLookupID1)),
			},
			{
				UCompact: types.NewUCompactFromUInt(uint64(tupleItemLookupID2)),
			},
		},
	}

	tupleItem1TypeDef := types.Si1TypeDef{
		IsPrimitive: true,
		Primitive: types.Si1TypeDefPrimitive{
			Si0TypeDefPrimitive: types.Si0TypeDefPrimitive(types.IsChar),
		},
	}

	tupleItem2TypeDef := types.Si1TypeDef{
		IsPrimitive: true,
		Primitive: types.Si1TypeDefPrimitive{
			Si0TypeDefPrimitive: types.Si0TypeDefPrimitive(types.IsI16),
		},
	}

	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			EfficientLookup: map[int64]*types.Si1Type{
				int64(tupleItemLookupID1): {
					Def: tupleItem1TypeDef,
				},
				int64(tupleItemLookupID2): {
					Def: tupleItem2TypeDef,
				},
			},
		},
	}

	factory := NewFactory().(*factory)

	res, err := factory.getFieldType(testMeta, testFieldName, testFieldTypeDef)
	assert.NoError(t, err)

	compositeFieldType, ok := res.(*CompositeFieldType)
	assert.True(t, ok)
	assert.Len(t, compositeFieldType.Fields, 2)
	assert.Equal(t, testFieldName, compositeFieldType.FieldName)

	assert.Equal(t, fmt.Sprintf(tupleItemFieldNameFormat, 0), compositeFieldType.Fields[0].Name)
	assert.Equal(t, &PrimitiveFieldType[byte]{}, compositeFieldType.Fields[0].FieldType)
	assert.Equal(t, int64(tupleItemLookupID1), compositeFieldType.Fields[0].LookupIndex)

	assert.Equal(t, fmt.Sprintf(tupleItemFieldNameFormat, 1), compositeFieldType.Fields[1].Name)
	assert.Equal(t, &PrimitiveFieldType[types.I16]{}, compositeFieldType.Fields[1].FieldType)
	assert.Equal(t, int64(tupleItemLookupID2), compositeFieldType.Fields[1].LookupIndex)
}

func TestFactory_getFieldType_Tuple_NilTuple(t *testing.T) {
	testFieldName := "TestFieldName"

	testFieldTypeDef := types.Si1TypeDef{
		IsTuple: true,
	}

	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			EfficientLookup: map[int64]*types.Si1Type{},
		},
	}

	factory := NewFactory().(*factory)

	res, err := factory.getFieldType(testMeta, testFieldName, testFieldTypeDef)
	assert.NoError(t, err)
	assert.Equal(t, &PrimitiveFieldType[[]any]{}, res)
}

func TestFactory_getFieldType_BitSequence(t *testing.T) {
	testFieldName := "TestFieldName"

	bitStoreLookupID := 123
	bitOrderLookupID := 456

	testFieldTypeDef := types.Si1TypeDef{
		IsBitSequence: true,
		BitSequence: types.Si1TypeDefBitSequence{
			BitStoreType: types.Si1LookupTypeID{
				UCompact: types.NewUCompactFromUInt(uint64(bitStoreLookupID)),
			},
			BitOrderType: types.Si1LookupTypeID{
				UCompact: types.NewUCompactFromUInt(uint64(bitOrderLookupID)),
			},
		},
	}

	bitStoreTypeDef := types.Si1TypeDef{
		IsPrimitive: true,
		Primitive: types.Si1TypeDefPrimitive{
			Si0TypeDefPrimitive: types.Si0TypeDefPrimitive(types.IsI64),
		},
	}

	bitOrderTypeDef := types.Si1TypeDef{
		IsPrimitive: true,
		Primitive: types.Si1TypeDefPrimitive{
			Si0TypeDefPrimitive: types.Si0TypeDefPrimitive(types.IsI256),
		},
	}

	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			EfficientLookup: map[int64]*types.Si1Type{
				int64(bitStoreLookupID): {
					Def: bitStoreTypeDef,
				},
				int64(bitOrderLookupID): {
					Def: bitOrderTypeDef,
				},
			},
		},
	}

	factory := NewFactory().(*factory)

	res, err := factory.getFieldType(testMeta, testFieldName, testFieldTypeDef)
	assert.NoError(t, err)

	bitSequenceType, ok := res.(*BitSequenceType)
	assert.True(t, ok)

	assert.Equal(t, &PrimitiveFieldType[types.I64]{}, bitSequenceType.BitStoreType)
	assert.Equal(t, &PrimitiveFieldType[types.I256]{}, bitSequenceType.BitOrderType)
}

func TestFactory_getFieldType_BitSequence_BitStoreTypeNotFound(t *testing.T) {
	testFieldName := "TestFieldName"

	bitStoreLookupID := 123
	bitOrderLookupID := 456

	testFieldTypeDef := types.Si1TypeDef{
		IsBitSequence: true,
		BitSequence: types.Si1TypeDefBitSequence{
			BitStoreType: types.Si1LookupTypeID{
				UCompact: types.NewUCompactFromUInt(uint64(bitStoreLookupID)),
			},
			BitOrderType: types.Si1LookupTypeID{
				UCompact: types.NewUCompactFromUInt(uint64(bitOrderLookupID)),
			},
		},
	}

	bitOrderTypeDef := types.Si1TypeDef{
		IsPrimitive: true,
		Primitive: types.Si1TypeDefPrimitive{
			Si0TypeDefPrimitive: types.Si0TypeDefPrimitive(types.IsI256),
		},
	}

	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			EfficientLookup: map[int64]*types.Si1Type{
				int64(bitOrderLookupID): {
					Def: bitOrderTypeDef,
				},
			},
		},
	}

	factory := NewFactory().(*factory)

	res, err := factory.getFieldType(testMeta, testFieldName, testFieldTypeDef)
	assert.Equal(t, "bit store type not found", err.Error())
	assert.Nil(t, res)
}

func TestFactory_getFieldType_BitSequence_BitStoreFieldTypeError(t *testing.T) {
	testFieldName := "TestFieldName"

	bitStoreLookupID := 123
	bitOrderLookupID := 456

	testFieldTypeDef := types.Si1TypeDef{
		IsBitSequence: true,
		BitSequence: types.Si1TypeDefBitSequence{
			BitStoreType: types.Si1LookupTypeID{
				UCompact: types.NewUCompactFromUInt(uint64(bitStoreLookupID)),
			},
			BitOrderType: types.Si1LookupTypeID{
				UCompact: types.NewUCompactFromUInt(uint64(bitOrderLookupID)),
			},
		},
	}

	bitStoreTypeDef := types.Si1TypeDef{
		IsComposite: true,
		Composite: types.Si1TypeDefComposite{
			Fields: []types.Si1Field{
				{
					Name: "BitStoreCompositeField1",
					Type: types.Si1LookupTypeID{
						// This type is not present in the efficient lookup map and should cause an error.
						UCompact: types.NewUCompactFromUInt(uint64(123456)),
					},
				},
			},
		},
	}

	bitOrderTypeDef := types.Si1TypeDef{
		IsPrimitive: true,
		Primitive: types.Si1TypeDefPrimitive{
			Si0TypeDefPrimitive: types.Si0TypeDefPrimitive(types.IsI256),
		},
	}

	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			EfficientLookup: map[int64]*types.Si1Type{
				int64(bitStoreLookupID): {
					Def: bitStoreTypeDef,
				},
				int64(bitOrderLookupID): {
					Def: bitOrderTypeDef,
				},
			},
		},
	}

	factory := NewFactory().(*factory)

	res, err := factory.getFieldType(testMeta, testFieldName, testFieldTypeDef)
	assert.Equal(t, "couldn't get bit store field type: couldn't get composite fields: type not found for field 'BitStoreCompositeField1'", err.Error())
	assert.Nil(t, res)
}

func TestFactory_getFieldType_BitSequence_BitOrderTypeNotFound(t *testing.T) {
	testFieldName := "TestFieldName"

	bitStoreLookupID := 123
	bitOrderLookupID := 456

	testFieldTypeDef := types.Si1TypeDef{
		IsBitSequence: true,
		BitSequence: types.Si1TypeDefBitSequence{
			BitStoreType: types.Si1LookupTypeID{
				UCompact: types.NewUCompactFromUInt(uint64(bitStoreLookupID)),
			},
			BitOrderType: types.Si1LookupTypeID{
				UCompact: types.NewUCompactFromUInt(uint64(bitOrderLookupID)),
			},
		},
	}

	bitStoreTypeDef := types.Si1TypeDef{
		IsPrimitive: true,
		Primitive: types.Si1TypeDefPrimitive{
			Si0TypeDefPrimitive: types.Si0TypeDefPrimitive(types.IsI64),
		},
	}

	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			EfficientLookup: map[int64]*types.Si1Type{
				int64(bitStoreLookupID): {
					Def: bitStoreTypeDef,
				},
			},
		},
	}

	factory := NewFactory().(*factory)

	res, err := factory.getFieldType(testMeta, testFieldName, testFieldTypeDef)
	assert.Equal(t, "bit order type not found", err.Error())
	assert.Nil(t, res)
}

func TestFactory_getFieldType_BitSequence_BitOrderFieldTypeError(t *testing.T) {
	testFieldName := "TestFieldName"

	bitStoreLookupID := 123
	bitOrderLookupID := 456

	testFieldTypeDef := types.Si1TypeDef{
		IsBitSequence: true,
		BitSequence: types.Si1TypeDefBitSequence{
			BitStoreType: types.Si1LookupTypeID{
				UCompact: types.NewUCompactFromUInt(uint64(bitStoreLookupID)),
			},
			BitOrderType: types.Si1LookupTypeID{
				UCompact: types.NewUCompactFromUInt(uint64(bitOrderLookupID)),
			},
		},
	}

	bitStoreTypeDef := types.Si1TypeDef{
		IsPrimitive: true,
		Primitive: types.Si1TypeDefPrimitive{
			Si0TypeDefPrimitive: types.Si0TypeDefPrimitive(types.IsI64),
		},
	}

	bitOrderTypeDef := types.Si1TypeDef{
		IsComposite: true,
		Composite: types.Si1TypeDefComposite{
			Fields: []types.Si1Field{
				{
					Name: "BitOrderCompositeField1",
					Type: types.Si1LookupTypeID{
						// This type is not present in the efficient lookup map and should cause an error.
						UCompact: types.NewUCompactFromUInt(uint64(123456)),
					},
				},
			},
		},
	}

	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			EfficientLookup: map[int64]*types.Si1Type{
				int64(bitStoreLookupID): {
					Def: bitStoreTypeDef,
				},
				int64(bitOrderLookupID): {
					Def: bitOrderTypeDef,
				},
			},
		},
	}

	factory := NewFactory().(*factory)

	res, err := factory.getFieldType(testMeta, testFieldName, testFieldTypeDef)
	assert.Equal(t, "couldn't get bit order field type: couldn't get composite fields: type not found for field 'BitOrderCompositeField1'", err.Error())
	assert.Nil(t, res)
}

func TestFactory_getVariantFieldType_CompositeVariantTypeFieldError(t *testing.T) {
	variantName1 := "Variant1"
	variantName2 := "Variant2"

	variantFieldName := "VariantFieldName"
	variantFieldLookupID := 123
	variantFieldTypeDef := types.Si1TypeDef{
		IsComposite: true,
		Composite: types.Si1TypeDefComposite{
			Fields: []types.Si1Field{
				{
					Name: "CompositeVariantField",
					Type: types.Si1LookupTypeID{
						// This lookup ID is not in the efficient lookup map which should cause an error.
						UCompact: types.NewUCompactFromUInt(uint64(456)),
					},
				},
			},
		},
	}

	testFieldTypeDef := types.Si1TypeDef{
		IsVariant: true,
		Variant: types.Si1TypeDefVariant{
			Variants: []types.Si1Variant{
				{
					Name:   types.Text(variantName1),
					Fields: nil,
					Index:  0,
				},
				{
					Name: types.Text(variantName2),
					Fields: []types.Si1Field{
						{
							HasName: true,
							Name:    types.Text(variantFieldName),
							Type: types.Si1LookupTypeID{
								UCompact: types.NewUCompactFromUInt(uint64(variantFieldLookupID)),
							},
						},
					},
					Index: 1,
				},
			},
		},
	}

	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			EfficientLookup: map[int64]*types.Si1Type{
				int64(variantFieldLookupID): {
					Def: variantFieldTypeDef,
				},
			},
		},
	}

	factory := NewFactory().(*factory)

	res, err := factory.getVariantFieldType(testMeta, testFieldTypeDef)
	assert.Equal(t, "couldn't get field types for variant '1': couldn't get field type for 'VariantFieldName': couldn't get composite fields: type not found for field 'CompositeVariantField'", err.Error())
	assert.Nil(t, res)
}

func TestFactory_getCompactFieldType_CompactTuple(t *testing.T) {
	testFieldName := "TestFieldName"

	tupleItemLookupID1 := 111
	tupleItemLookupID2 := 222

	compactFieldTypeDef := types.Si1TypeDef{
		IsTuple: true,
		Tuple: []types.Si1LookupTypeID{
			{
				UCompact: types.NewUCompactFromUInt(uint64(tupleItemLookupID1)),
			},
			{
				UCompact: types.NewUCompactFromUInt(uint64(tupleItemLookupID2)),
			},
		},
	}

	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			EfficientLookup: map[int64]*types.Si1Type{
				int64(tupleItemLookupID1): {
					Def: types.Si1TypeDef{
						IsPrimitive: true,
						Primitive: types.Si1TypeDefPrimitive{
							Si0TypeDefPrimitive: types.Si0TypeDefPrimitive(types.IsU8),
						},
					},
				},
				int64(tupleItemLookupID2): {
					Def: types.Si1TypeDef{
						IsPrimitive: true,
						Primitive: types.Si1TypeDefPrimitive{
							Si0TypeDefPrimitive: types.Si0TypeDefPrimitive(types.IsU32),
						},
					},
				},
			},
		},
	}

	factory := NewFactory().(*factory)

	res, err := factory.getCompactFieldType(testMeta, testFieldName, compactFieldTypeDef)
	assert.NoError(t, err)

	compositeFieldType, ok := res.(*CompositeFieldType)
	assert.True(t, ok)
	assert.Len(t, compositeFieldType.Fields, 2)

	assert.Equal(t, fmt.Sprintf(tupleItemFieldNameFormat, 0), compositeFieldType.Fields[0].Name)
	assert.Equal(t, &PrimitiveFieldType[types.UCompact]{}, compositeFieldType.Fields[0].FieldType)
	assert.Equal(t, int64(tupleItemLookupID1), compositeFieldType.Fields[0].LookupIndex)
	assert.Equal(t, fmt.Sprintf(tupleItemFieldNameFormat, 1), compositeFieldType.Fields[1].Name)
	assert.Equal(t, &PrimitiveFieldType[types.UCompact]{}, compositeFieldType.Fields[1].FieldType)
	assert.Equal(t, int64(tupleItemLookupID2), compositeFieldType.Fields[1].LookupIndex)
}

func TestFactory_getCompactFieldType_CompactComposite(t *testing.T) {
	testFieldName := "TestFieldName"

	compositeFieldName1 := "CompositeFieldName1"
	compositeFieldName2 := "CompositeFieldName2"

	compositeFieldLookupID1 := 111
	compositeFieldLookupID2 := 222

	compactFieldTypeDef := types.Si1TypeDef{
		IsComposite: true,
		Composite: types.Si1TypeDefComposite{
			Fields: []types.Si1Field{
				{
					HasName: true,
					Name:    types.Text(compositeFieldName1),
					Type: types.Si1LookupTypeID{
						UCompact: types.NewUCompactFromUInt(uint64(compositeFieldLookupID1)),
					},
				},
				{
					HasName: true,
					Name:    types.Text(compositeFieldName2),
					Type: types.Si1LookupTypeID{
						UCompact: types.NewUCompactFromUInt(uint64(compositeFieldLookupID2)),
					},
				},
			},
		},
	}

	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			EfficientLookup: map[int64]*types.Si1Type{
				int64(compositeFieldLookupID1): {
					Def: types.Si1TypeDef{
						IsPrimitive: true,
						Primitive: types.Si1TypeDefPrimitive{
							Si0TypeDefPrimitive: types.Si0TypeDefPrimitive(types.IsU8),
						},
					},
				},
				int64(compositeFieldLookupID2): {
					Def: types.Si1TypeDef{
						IsPrimitive: true,
						Primitive: types.Si1TypeDefPrimitive{
							Si0TypeDefPrimitive: types.Si0TypeDefPrimitive(types.IsU32),
						},
					},
				},
			},
		},
	}

	factory := NewFactory().(*factory)

	res, err := factory.getCompactFieldType(testMeta, testFieldName, compactFieldTypeDef)
	assert.NoError(t, err)

	compositeFieldType, ok := res.(*CompositeFieldType)
	assert.True(t, ok)
	assert.Len(t, compositeFieldType.Fields, 2)

	assert.Equal(t, compositeFieldName1, compositeFieldType.Fields[0].Name)
	assert.Equal(t, &PrimitiveFieldType[types.UCompact]{}, compositeFieldType.Fields[0].FieldType)
	assert.Equal(t, int64(compositeFieldLookupID1), compositeFieldType.Fields[0].LookupIndex)
	assert.Equal(t, compositeFieldName2, compositeFieldType.Fields[1].Name)
	assert.Equal(t, &PrimitiveFieldType[types.UCompact]{}, compositeFieldType.Fields[1].FieldType)
	assert.Equal(t, int64(compositeFieldLookupID2), compositeFieldType.Fields[1].LookupIndex)
}

func TestFactory_getArrayFieldType(t *testing.T) {
	testFieldName := "TestFieldName"

	arrayLen := 32

	arrayItemTypeDef := types.Si1TypeDef{
		IsPrimitive: true,
		Primitive: types.Si1TypeDefPrimitive{
			Si0TypeDefPrimitive: types.Si0TypeDefPrimitive(types.IsU8),
		},
	}

	testMeta := &types.Metadata{}

	factory := NewFactory().(*factory)

	res, err := factory.getArrayFieldType(uint(arrayLen), testMeta, testFieldName, arrayItemTypeDef)
	assert.NoError(t, err)

	arrayFieldType, ok := res.(*ArrayFieldType)
	assert.True(t, ok)

	assert.Equal(t, uint(arrayLen), arrayFieldType.Length)
	assert.Equal(t, &PrimitiveFieldType[types.U8]{}, arrayFieldType.ItemType)
}

func TestFactory_getArrayFieldType_ItemFieldTypeError(t *testing.T) {
	testFieldName := "TestFieldName"

	arrayLen := 32

	compositeLookupID := 123

	arrayItemTypeDef := types.Si1TypeDef{
		IsComposite: true,
		Composite: types.Si1TypeDefComposite{
			Fields: []types.Si1Field{
				{
					HasName: true,
					Name:    "CompositeField1",
					Type: types.Si1LookupTypeID{
						// This lookup ID is not present in the efficient lookup map which should cause an error.
						UCompact: types.NewUCompactFromUInt(uint64(compositeLookupID)),
					},
				},
			},
		},
	}

	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			EfficientLookup: map[int64]*types.Si1Type{},
		},
	}

	factory := NewFactory().(*factory)

	res, err := factory.getArrayFieldType(uint(arrayLen), testMeta, testFieldName, arrayItemTypeDef)
	assert.Equal(t, "couldn't get array item field type: couldn't get composite fields: type not found for field 'CompositeField1'", err.Error())
	assert.Nil(t, res)
}

func TestFactory_getSliceFieldType(t *testing.T) {
	testFieldName := "TestFieldName"

	sliceItemTypeDef := types.Si1TypeDef{
		IsPrimitive: true,
		Primitive: types.Si1TypeDefPrimitive{
			Si0TypeDefPrimitive: types.Si0TypeDefPrimitive(types.IsU8),
		},
	}

	testMeta := &types.Metadata{}

	factory := NewFactory().(*factory)

	res, err := factory.getSliceFieldType(testMeta, testFieldName, sliceItemTypeDef)
	assert.NoError(t, err)

	sliceFieldType, ok := res.(*SliceFieldType)
	assert.True(t, ok)

	assert.Equal(t, &PrimitiveFieldType[types.U8]{}, sliceFieldType.ItemType)
}

func TestFactory_getSliceFieldType_ItemFieldTypeError(t *testing.T) {
	testFieldName := "TestFieldName"

	compositeLookupID := 123

	sliceItemTypeDef := types.Si1TypeDef{
		IsComposite: true,
		Composite: types.Si1TypeDefComposite{
			Fields: []types.Si1Field{
				{
					HasName: true,
					Name:    "CompositeField1",
					Type: types.Si1LookupTypeID{
						// This lookup ID is not present in the efficient lookup map which should cause an error.
						UCompact: types.NewUCompactFromUInt(uint64(compositeLookupID)),
					},
				},
			},
		},
	}

	testMeta := &types.Metadata{}

	factory := NewFactory().(*factory)

	res, err := factory.getSliceFieldType(testMeta, testFieldName, sliceItemTypeDef)
	assert.Equal(t, "couldn't get slice item field type: couldn't get composite fields: type not found for field 'CompositeField1'", err.Error())
	assert.Nil(t, res)
}

func TestFactory_getTupleType(t *testing.T) {
	testFieldName := "TestFieldName"

	tupleItemLookupID1 := 123
	tupleItemLookupID2 := 456

	tupleTypeDef := []types.Si1LookupTypeID{
		{
			UCompact: types.NewUCompactFromUInt(uint64(tupleItemLookupID1)),
		},
		{
			UCompact: types.NewUCompactFromUInt(uint64(tupleItemLookupID2)),
		},
	}

	tupleItem1TypeDef := types.Si1TypeDef{
		IsPrimitive: true,
		Primitive: types.Si1TypeDefPrimitive{
			Si0TypeDefPrimitive: types.Si0TypeDefPrimitive(types.IsU8),
		},
	}

	tupleItem2TypeDef := types.Si1TypeDef{
		IsPrimitive: true,
		Primitive: types.Si1TypeDefPrimitive{
			Si0TypeDefPrimitive: types.Si0TypeDefPrimitive(types.IsU32),
		},
	}

	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			EfficientLookup: map[int64]*types.Si1Type{
				int64(tupleItemLookupID1): {
					Def: tupleItem1TypeDef,
				},
				int64(tupleItemLookupID2): {
					Def: tupleItem2TypeDef,
				},
			},
		},
	}

	factory := NewFactory().(*factory)

	res, err := factory.getTupleType(testMeta, testFieldName, tupleTypeDef)
	assert.NoError(t, err)

	compositeFieldType, ok := res.(*CompositeFieldType)
	assert.True(t, ok)
	assert.Len(t, compositeFieldType.Fields, 2)

	assert.Equal(t, fmt.Sprintf(tupleItemFieldNameFormat, 0), compositeFieldType.Fields[0].Name)
	assert.Equal(t, &PrimitiveFieldType[types.U8]{}, compositeFieldType.Fields[0].FieldType)
	assert.Equal(t, int64(tupleItemLookupID1), compositeFieldType.Fields[0].LookupIndex)
	assert.Equal(t, fmt.Sprintf(tupleItemFieldNameFormat, 1), compositeFieldType.Fields[1].Name)
	assert.Equal(t, &PrimitiveFieldType[types.U32]{}, compositeFieldType.Fields[1].FieldType)
	assert.Equal(t, int64(tupleItemLookupID2), compositeFieldType.Fields[1].LookupIndex)
}

func TestFactory_getTupleType_TupleItemNotFound(t *testing.T) {
	testFieldName := "TestFieldName"

	tupleItemLookupID1 := 123
	tupleItemLookupID2 := 456

	tupleTypeDef := []types.Si1LookupTypeID{
		{
			UCompact: types.NewUCompactFromUInt(uint64(tupleItemLookupID1)),
		},
		{
			UCompact: types.NewUCompactFromUInt(uint64(tupleItemLookupID2)),
		},
	}

	tupleItem1TypeDef := types.Si1TypeDef{
		IsPrimitive: true,
		Primitive: types.Si1TypeDefPrimitive{
			Si0TypeDefPrimitive: types.Si0TypeDefPrimitive(types.IsU8),
		},
	}

	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			EfficientLookup: map[int64]*types.Si1Type{
				int64(tupleItemLookupID1): {
					Def: tupleItem1TypeDef,
				},
				// Lookup ID for tuple item 2 is missing which should cause an error.
			},
		},
	}

	factory := NewFactory().(*factory)

	res, err := factory.getTupleType(testMeta, testFieldName, tupleTypeDef)
	assert.Equal(t, "type definition for tuple item 1 not found", err.Error())
	assert.Nil(t, res)
}

func TestFactory_getTupleType_TupleItemFieldTypeError(t *testing.T) {
	testFieldName := "TestFieldName"

	tupleItemLookupID1 := 123
	tupleItemLookupID2 := 456

	tupleTypeDef := []types.Si1LookupTypeID{
		{
			UCompact: types.NewUCompactFromUInt(uint64(tupleItemLookupID1)),
		},
		{
			UCompact: types.NewUCompactFromUInt(uint64(tupleItemLookupID2)),
		},
	}

	tupleItem1TypeDef := types.Si1TypeDef{
		IsPrimitive: true,
		Primitive: types.Si1TypeDefPrimitive{
			Si0TypeDefPrimitive: types.Si0TypeDefPrimitive(types.IsU8),
		},
	}

	tupleItem2TypeDef := types.Si1TypeDef{
		IsComposite: true,
		Composite: types.Si1TypeDefComposite{
			Fields: []types.Si1Field{
				{
					HasName: true,
					Name:    "CompositeField1",
					Type: types.Si1LookupTypeID{
						// This lookup ID is not in the efficient lookup map which should cause an error.
						UCompact: types.NewUCompactFromUInt(uint64(789)),
					},
				},
			},
		},
	}

	testMeta := &types.Metadata{
		AsMetadataV14: types.MetadataV14{
			EfficientLookup: map[int64]*types.Si1Type{
				int64(tupleItemLookupID1): {
					Def: tupleItem1TypeDef,
				},
				int64(tupleItemLookupID2): {
					Def: tupleItem2TypeDef,
				},
			},
		},
	}

	factory := NewFactory().(*factory)

	res, err := factory.getTupleType(testMeta, testFieldName, tupleTypeDef)
	assert.Equal(t, "couldn't get tuple field type: couldn't get composite fields: type not found for field 'CompositeField1'", err.Error())
	assert.Nil(t, res)
}

func Test_getPrimitiveType_UnsupportedTypeError(t *testing.T) {
	primitiveTypeDef := types.Si0TypeDefPrimitive(32)

	res, err := getPrimitiveType(primitiveTypeDef)
	assert.NotNil(t, err)
	assert.Nil(t, res)
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
