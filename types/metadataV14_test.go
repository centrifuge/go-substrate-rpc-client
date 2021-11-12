package types_test

import (
	"fmt"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
	. "github.com/centrifuge/go-substrate-rpc-client/v3/types"
	"github.com/stretchr/testify/assert"
)

func TestMetadataV14_TestFindCallIndexWithUnknownFunction(t *testing.T) {
	var metadata Metadata
	err := DecodeFromHexString(MetadataV14Data, &metadata)
	assert.EqualValues(t, metadata.Version, 14)
	assert.NoError(t, err)

	_, err = metadata.FindCallIndex("Module2_14.unknownFunction")
	assert.Error(t, err)
}

// Test that decoding the example metadata v14 doesn't fail
func TestMetadataV14_Decode(t *testing.T) {
	// Verify that we can succcessfully decode metadata v14
	var metadata Metadata
	err := DecodeFromHexString(MetadataV14Data, &metadata)
	assert.EqualValues(t, metadata.Version, 14)
	assert.NoError(t, err)
}

// Test that decoding the example metadata v14 doesn't fail
func TestMetadataV14_Debug_Extrinsic(t *testing.T) {
	// Verify that we can succcessfully decode metadata v14
	var metadata Metadata
	err := DecodeFromHexString(MetadataV14Data, &metadata)
	assert.EqualValues(t, metadata.Version, 14)
	assert.NoError(t, err)

	encoded, err := EncodeToBytes(metadata.AsMetadataV14.Extrinsic)
	assert.NoError(t, err)

	var decodedExtrinsic ExtrinsicV14
	err = DecodeFromBytes(encoded, &decodedExtrinsic)
	assert.NoError(t, err)

	assert.Equal(t, metadata.AsMetadataV14.Extrinsic, decodedExtrinsic)
}

// Test that decoding the example metadata v14 doesn't fail
func TestMetadataV14_Debug_Lookup(t *testing.T) {
	// Verify that we can succcessfully decode metadata v14
	var metadata Metadata
	err := DecodeFromHexString(MetadataV14Data, &metadata)
	assert.EqualValues(t, metadata.Version, 14)
	assert.NoError(t, err)

	encoded, err := EncodeToBytes(metadata.AsMetadataV14.Lookup)
	assert.NoError(t, err)

	var decodedLookup PortableRegistry
	err = DecodeFromBytes(encoded, &decodedLookup)
	assert.NoError(t, err, "oopsi")

	assert.Equal(t, metadata.AsMetadataV14.Lookup, decodedLookup)

}

// Test that decoding the example metadata v14 doesn't fail
func TestMetadataV14_Debug_Type(t *testing.T) {
	// Verify that we can succcessfully decode metadata v14
	var metadata Metadata
	err := DecodeFromHexString(MetadataV14Data, &metadata)
	assert.EqualValues(t, metadata.Version, 14)
	assert.NoError(t, err)

	encoded, err := EncodeToBytes(metadata.AsMetadataV14.Type)
	assert.NoError(t, err)

	var decodedTypz Si1LookupTypeID
	err = DecodeFromBytes(encoded, &decodedTypz)
	assert.NoError(t, err)

	assert.Equal(t, metadata.AsMetadataV14.Type, decodedTypz)
}

// Test that decoding the example metadata v14 doesn't fail
func TestMetadataV14_Debug_Pallets(t *testing.T) {
	// Verify that we can succcessfully decode metadata v14
	var metadata Metadata
	err := DecodeFromHexString(MetadataV14Data, &metadata)
	assert.EqualValues(t, metadata.Version, 14)
	assert.NoError(t, err)

	encoded, err := EncodeToBytes(metadata.AsMetadataV14.Pallets)
	assert.NoError(t, err)

	var decodedPallets PalletMetadataV14
	err = DecodeFromBytes(encoded, &decodedPallets)
	assert.NoError(t, err)

	assert.Equal(t, metadata.AsMetadataV14.Pallets, decodedPallets)
}

// Verify that (Decode . Encode) outputs the input.
func TestMetadataV14_Decode_Rountrip(t *testing.T) {
	// Decode the metadata
	var metadata Metadata
	err := DecodeFromHexString(MetadataV14Data, &metadata)
	assert.EqualValues(t, metadata.Version, 14)
	assert.NoError(t, err)

	// Now encode it
	encoded, err := EncodeToHexString(metadata)
	assert.NoError(t, err)

	// Verify it equals the original metadata
	assert.Equal(t, MetadataV14Data, encoded)

	// Verify that decoding the encoded metadata
	// equals the decoded original metadata
	var decodedMetadata Metadata
	err = DecodeFromHexString(encoded, &decodedMetadata)
	assert.NoError(t, err)

	assert.EqualValues(t, metadata, decodedMetadata)
}

// Verify that decoding the metadata v14 hex string twice
// produces the same output.
func TestMetadataV14_DecodeTwice(t *testing.T) {
	// Verify that we can succcessfully decode metadata v14
	var metadata1 Metadata
	err := DecodeFromHexString(MetadataV14Data, &metadata1)
	assert.EqualValues(t, metadata1.Version, 14)
	assert.NoError(t, err)

	// Decode it again
	var metadata2 Metadata
	err = DecodeFromHexString(MetadataV14Data, &metadata2)
	assert.EqualValues(t, metadata2.Version, 14)
	assert.NoError(t, err)

	// Verify they are the same value
	assert.EqualValues(t, metadata1, metadata2)
}

// Verify that we can find the index of a valid call
func TestMetadataV14FindCallIndex(t *testing.T) {
	var meta Metadata
	err := DecodeFromHexString(MetadataV14Data, &meta)
	assert.NoError(t, err)
	index, err := meta.FindCallIndex("Balances.transfer")
	assert.NoError(t, err)
	assert.Equal(t, index, CallIndex{SectionIndex: 5, MethodIndex: 0})
}

// Verify that we get an error when querying for an invalid
// call with FindCallIndex.
func TestMetadataV14FindCallIndexNonExistent(t *testing.T) {
	var meta Metadata
	err := DecodeFromHexString(MetadataV14Data, &meta)
	assert.NoError(t, err)
	_, err = meta.FindCallIndex("Doesnt.Exist")
	assert.Error(t, err)
}

// TODO(nuno): make verifications more meaningful
func TestMetadataV14FindEventNamesForEventID(t *testing.T) {
	var meta Metadata
	err := DecodeFromHexString(MetadataV14Data, &meta)
	if err != nil {
		t.Fatal(err)
	}
	id := EventID{5, 2}
	_, _, err = meta.FindEventNamesForEventID(id)
	assert.NoError(t, err)
}

// TODO(nuno): make verifications more meaningful
func TestMetadataV14FindStorageEntryMetadata(t *testing.T) {
	var meta Metadata
	err := DecodeFromHexString(MetadataV14Data, &meta)
	assert.NoError(t, err)

	_, err = meta.FindStorageEntryMetadata("System", "Account")
	assert.NoError(t, err)
}

// Verify FindStorageEntryMetadata returns an err when
// the given module can't be found.
func TestMetadataV14FindStorageEntryMetadata_InvalidModule(t *testing.T) {
	var meta Metadata
	err := DecodeFromHexString(MetadataV14Data, &meta)
	assert.NoError(t, err)

	_, err = meta.FindStorageEntryMetadata("SystemZ", "Account")
	fmt.Println(err)
	assert.NoError(t, err)
}

// Verify FindStorageEntryMetadata returns an err when
// it doesn't find a storage within an existing module.
func TestMetadataV14FindStorageEntryMetadata_InvalidStorage(t *testing.T) {
	var meta Metadata
	err := DecodeFromHexString(MetadataV14Data, &meta)
	assert.NoError(t, err)

	_, err = meta.FindStorageEntryMetadata("System", "Accountz")
	assert.Error(t, err)
}

func TestMetadataV14ExistsModuleMetadata(t *testing.T) {
	var meta Metadata
	err := DecodeFromHexString(MetadataV14Data, &meta)
	if err != nil {
		t.Fatal(err)
	}
	res := meta.ExistsModuleMetadata("System")
	assert.True(t, res)
}

func TestMetadataV14_Pallet_Empty(t *testing.T) {
	var pallet = PalletMetadataV14{
		Name:       NewText("System"),
		HasStorage: false,
		HasCalls:   false,
		HasEvents:  false,
		Constants:  nil,
		HasErrors:  false,
		Index:      42,
	}

	encoded, err := EncodeToBytes(pallet)
	assert.NoError(t, err)

	var encodedPallets PalletMetadataV14
	err = DecodeFromBytes(encoded, &encodedPallets)
	assert.NoError(t, err)

	// Verify they are the same value
	assert.EqualValues(t, encodedPallets, pallet)
}

func TestMetadataV14_Pallet_Filled(t *testing.T) {
	var pallet = PalletMetadataV14{
		Name:       NewText("System"),
		HasStorage: true,
		Storage: StorageMetadataV14{
			Prefix: "Pre-fix",
			Items: []StorageEntryMetadataV14{
				{
					Name:     "StorageName",
					Modifier: StorageFunctionModifierV0{IsOptional: true},
					Type: StorageEntryTypeV14{
						IsPlainType: false,
						IsMap:       true,
						AsMap: MapTypeV14{
							Hashers: []StorageHasherV10{
								{IsBlake2_128: true}, {IsBlake2_256: true},
							},
							Key:   NewSi1LookupTypeIDFromUInt(3),
							Value: NewSi1LookupTypeIDFromUInt(4),
						},
					},
				},
				{
					Name: "Account",
					Modifier: types.StorageFunctionModifierV0{
						IsOptional: false,
						IsDefault:  true,
						IsRequired: false,
					},
					Type: types.StorageEntryTypeV14{
						IsPlainType: false,
						IsMap:       true,
						AsMap: types.MapTypeV14{
							Hashers: []types.StorageHasherV10{
								{
									IsBlake2_128:       false,
									IsBlake2_256:       false,
									IsBlake2_128Concat: true,
									IsTwox128:          false,
									IsTwox256:          false,
									IsTwox64Concat:     false,
									IsIdentity:         false,
								},
							},
						},
					},
					Fallback: types.Bytes{
						0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
					},
					Documentation: []types.Text{" The full account information for a particular account ID."},
				},
			},
		},
		HasCalls:  true,
		Calls:     FunctionMetadataV14{Type: NewSi1LookupTypeIDFromUInt(24)},
		HasEvents: true,
		Events:    EventMetadataV14{Type: NewSi1LookupTypeIDFromUInt(72)},
		Constants: []ConstantMetadataV14{
			{
				Name:  NewText("Yellow"),
				Type:  NewSi1LookupTypeIDFromUInt(83),
				Value: []byte("Valuez"),
				Docs:  []Text{"README", "Contribute"},
			},
		},
		HasErrors: true,
		Errors:    ErrorMetadataV14{Type: NewSi1LookupTypeIDFromUInt(57)},
		Index:     42,
	}

	encoded, err := EncodeToBytes(pallet)
	assert.NoError(t, err)

	var encodedPallets PalletMetadataV14
	err = DecodeFromBytes(encoded, &encodedPallets)
	assert.NoError(t, err)

	// Verify they are the same
	assert.Equal(t, encodedPallets, pallet)
}

func TestSi1Type(t *testing.T) {
	type Si1Type struct {
		Path   Si1Path
		Params []Si1TypeParameter
		Def    Si1TypeDef
		Docs   []Text
	}

	// Replicate the first Si1Type we get from rpc json, marsh it, and aside encode it, and decode it
	var ti = Si1Type{
		Path: []Text{"sp_core", "crypto", "AccountId32"},
		Def: Si1TypeDef{
			IsComposite: true,
			Composite: Si1TypeDefComposite{
				Fields: []Si1Field{
					{
						Type:        NewSi1LookupTypeIDFromUInt(1),
						HasTypeName: true,
						TypeName:    NewText("[u8; 32]"),
					},
				},
			},
		},
	}

	// Verify that (decode . encode) equals the original value
	encoded, err := EncodeToHexString(ti)
	assert.NoError(t, err)

	var decoded Si1Type
	DecodeFromHexString(encoded, &decoded)

	assert.Equal(t, ti, decoded)
}
