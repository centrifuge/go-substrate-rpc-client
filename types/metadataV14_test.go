package types_test

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/centrifuge/go-substrate-rpc-client/v3/types"
)

// var exampleMetadataV14 = Metadata{
// 	MagicNumber:   0x6174656d,
// 	Version:       13,
// 	AsMetadataV14: exampleRuntimeMetadataV14,
// }

// // var exampleRuntimeMetadataV14 = MetadataV14{
// // 	Pallets: []PalletMetadataV14{examplePalletMetadataV14Empty, examplePalletMetadataV141, examplePalletMetadataV142},
// // }

// // // var examplePalletMetadataV14Empty = PalletMetadataV14{
// // // 	Name:       "EmptyModule_14",
// // // 	HasStorage: false,
// // // 	Storage:    StorageMetadataV14{},
// // // 	HasCalls:   false,
// // // 	Calls:      nil,
// // // 	HasEvents:  false,
// // // 	Events:     nil,
// // // 	Constants:  nil,
// // // 	Errors:     nil,
// // // 	Index:      0,
// // // }

// // // var examplePalletMetadataV141 = PalletMetadataV14{
// // // 	Name:       "Module1_14",
// // // 	HasStorage: true,
// // // 	Storage:    exampleStorageMetadataV14,
// // // 	HasCalls:   true,
// // // 	Calls:      []FunctionMetadataV4{exampleFunctionMetadataV4},
// // // 	HasEvents:  true,
// // // 	Events:     []EventMetadataV4{exampleEventMetadataV4},
// // // 	Constants:  []ModuleConstantMetadataV6{exampleModuleConstantMetadataV6},
// // // 	Errors:     []ErrorMetadataV8{exampleErrorMetadataV8},
// // // 	Index:      1,
// // // }

// // // var examplePalletMetadataV142 = PalletMetadataV14{
// // // 	Name:       "Module2_14",
// // // 	HasStorage: true,
// // // 	Storage:    exampleStorageMetadataV14,
// // // 	HasCalls:   true,
// // // 	Calls:      []FunctionMetadataV4{exampleFunctionMetadataV4},
// // // 	HasEvents:  true,
// // // 	Events:     []EventMetadataV4{exampleEventMetadataV4},
// // // 	Constants:  []ModuleConstantMetadataV6{exampleModuleConstantMetadataV6},
// // // 	Errors:     []ErrorMetadataV8{exampleErrorMetadataV8},
// // // 	Index:      2,
// // // }

// // // var exampleStorageMetadataV14 = StorageMetadataV14{
// // // 	Prefix: "myStoragePrefix_14",
// // // 	Items: []StorageFunctionMetadataV14{exampleStorageFunctionMetadataV14Type, exampleStorageFunctionMetadataV14Map,
// // // 		exampleStorageFunctionMetadataV14DoubleMap, exampleStorageFunctionMetadataV14NMap},
// // // }

// // // var exampleStorageFunctionMetadataV14Type = StorageFunctionMetadataV14{
// // // 	Name:          "myStorageFunc_14",
// // // 	Modifier:      StorageFunctionModifierV0{IsOptional: true},
// // // 	Type:          StorageFunctionTypeV14{IsType: true, AsType: "U8"},
// // // 	Fallback:      []byte{23, 14},
// // // 	Documentation: []Text{"My", "storage func", "doc"},
// // // }

// // // var exampleStorageFunctionMetadataV14Map = StorageFunctionMetadataV14{
// // // 	Name:          "myStorageFunc2_14",
// // // 	Modifier:      StorageFunctionModifierV0{IsOptional: true},
// // // 	Type:          StorageFunctionTypeV14{IsMap: true, AsMap: exampleMapTypeV10},
// // // 	Fallback:      []byte{23, 14},
// // // 	Documentation: []Text{"My", "storage func", "doc"},
// // // }

// // // var exampleStorageFunctionMetadataV14DoubleMap = StorageFunctionMetadataV14{
// // // 	Name:          "myStorageFunc3_14",
// // // 	Modifier:      StorageFunctionModifierV0{IsOptional: true},
// // // 	Type:          StorageFunctionTypeV14{IsDoubleMap: true, AsDoubleMap: exampleDoubleMapTypeV10},
// // // 	Fallback:      []byte{23, 14},
// // // 	Documentation: []Text{"My", "storage func", "doc"},
// // // }

// // // var exampleStorageFunctionMetadataV14NMap = StorageFunctionMetadataV14{
// // // 	Name:          "myStorageFunc4_14",
// // // 	Modifier:      StorageFunctionModifierV0{IsOptional: true},
// // // 	Type:          StorageFunctionTypeV14{IsNMap: true, AsNMap: exampleNMapTypeV14},
// // // 	Fallback:      []byte{23, 14},
// // // 	Documentation: []Text{"My", "storage func", "doc"},
// // // }

// var exampleNMapTypeV14 = NMapTypeV14{
// 	Hashers: []StorageHasherV10{{IsBlake2_256: true}, {IsBlake2_128Concat: true}, {IsIdentity: true}},
// 	Keys:    []Type{"myKey1", "myKey2", "myKey3"},
// 	Value:   "and a value",
// }

// func TestMetadataV14_ExistsModuleMetadata(t *testing.T) {
// 	assert.True(t, exampleMetadataV14.ExistsModuleMetadata("EmptyModule_14"))
// 	assert.False(t, exampleMetadataV14.ExistsModuleMetadata("NotExistModule"))
// }

// func TestMetadataV14_FindEventNamesForEventID(t *testing.T) {
// 	module, event, err := exampleMetadataV14.FindEventNamesForEventID(EventID([2]byte{1, 0}))

// 	assert.NoError(t, err)
// 	assert.Equal(t, examplePalletMetadataV141.Name, module)
// 	assert.Equal(t, exampleEventMetadataV4.Name, event)
// }

// func TestMetadataV14_FindEventNamesForUnknownModule(t *testing.T) {
// 	_, _, err := exampleMetadataV14.FindEventNamesForEventID(EventID([2]byte{1, 18}))

// 	assert.Error(t, err)
// }

// func TestMetadataV14_TestFindStorageEntryMetadata(t *testing.T) {
// 	_, err := exampleMetadataV14.FindStorageEntryMetadata("myStoragePrefix_14", "myStorageFunc2_14")
// 	assert.NoError(t, err)
// }

// func TestMetadataV14_TestFindCallIndex(t *testing.T) {
// 	callIndex, err := exampleMetadataV14.FindCallIndex("Module2_14.my function")
// 	assert.NoError(t, err)
// 	assert.Equal(t, examplePalletMetadataV142.Index, callIndex.SectionIndex)
// 	assert.Equal(t, uint8(0), callIndex.MethodIndex)
// }

// func TestMetadataV14_TestFindCallIndexWithUnknownModule(t *testing.T) {
// 	_, err := exampleMetadataV14.FindCallIndex("UnknownModule.my function")
// 	assert.Error(t, err)
// }

// func TestMetadataV14_TestFindCallIndexWithUnknownFunction(t *testing.T) {
// 	_, err := exampleMetadataV14.FindCallIndex("Module2_14.unknownFunction")
// 	assert.Error(t, err)
// }

// func TestNewMetadataV14_Decode(t *testing.T) {
// 	metadata := NewMetadataV14()
// 	err := DecodeFromBytes(MustHexDecodeString(ExamplaryMetadataV14SubstrateString), metadata)
// 	assert.EqualValues(t, metadata.Version, 13)
// 	assert.NoError(t, err)
// 	data, err := EncodeToBytes(metadata)
// 	assert.NoError(t, err)
// 	assert.Equal(t, ExamplaryMetadataV14SubstrateString, HexEncodeToString(data))
// }

func Test_ParseMetadata(t *testing.T) {
	var meta Metadata
	err := DecodeFromHexString(MetadataV14Data, &meta)
	if err != nil {
		t.Fatal(err)
	}
	v14 := meta.AsMetadataV14
	d, _ := json.Marshal(v14)
	fmt.Println(string(d))
}

func TestMetadataV14FindCallIndex(t *testing.T) {
	var meta Metadata
	err := DecodeFromHexString(MetadataV14Data, &meta)
	if err != nil {
		t.Fatal(err)
	}
	callIdx, err := meta.FindCallIndex("Balances.transfer")
	if err != nil {
		panic(err)
	}
	fmt.Println(callIdx)
}
func TestMetadataV14FindEventNamesForEventID(t *testing.T) {
	var meta Metadata
	err := DecodeFromHexString(MetadataV14Data, &meta)
	if err != nil {
		t.Fatal(err)
	}
	id := EventID{}
	id[0] = 5
	id[1] = 2
	mod, event, err := meta.FindEventNamesForEventID(id)
	if err != nil {
		panic(err)
	}
	fmt.Println(mod, event)
}

func TestMetadataV14FindStorageEntryMetadata(t *testing.T) {
	var meta Metadata
	err := DecodeFromHexString(MetadataV14Data, &meta)
	if err != nil {
		t.Fatal(err)
	}
	st, err := meta.FindStorageEntryMetadata("System", "Account")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(st)
}

func TestMetadataV14ExistsModuleMetadata(t *testing.T) {
	var meta Metadata
	err := DecodeFromHexString(MetadataV14Data, &meta)
	if err != nil {
		t.Fatal(err)
	}
	s := meta.ExistsModuleMetadata("System")

	fmt.Println(s)
}
