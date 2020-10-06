package types_test

import (
	"testing"

	. "github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/stretchr/testify/assert"
)

var exampleMetadataV12 = Metadata{
	MagicNumber:   0x6174656d,
	Version:       12,
	IsMetadataV12: true,
	AsMetadataV12: exampleRuntimeMetadataV12,
}

var exampleRuntimeMetadataV12 = MetadataV12{
	Modules: []ModuleMetadataV12{exampleModuleMetadataV12Empty, exampleModuleMetadataV121, exampleModuleMetadataV122},
}

var exampleModuleMetadataV12Empty = ModuleMetadataV12{
	Name:       "EmptyModule",
	HasStorage: false,
	Storage:    StorageMetadataV10{},
	HasCalls:   false,
	Calls:      nil,
	HasEvents:  false,
	Events:     nil,
	Constants:  nil,
	Errors:     nil,
	Index:      0,
}

var exampleModuleMetadataV121 = ModuleMetadataV12{
	Name:       "Module1",
	HasStorage: true,
	Storage:    exampleStorageMetadataV10,
	HasCalls:   true,
	Calls:      []FunctionMetadataV4{exampleFunctionMetadataV4},
	HasEvents:  true,
	Events:     []EventMetadataV4{exampleEventMetadataV4},
	Constants:  []ModuleConstantMetadataV6{exampleModuleConstantMetadataV6},
	Errors:     []ErrorMetadataV8{exampleErrorMetadataV8},
	Index:      1,
}

var exampleModuleMetadataV122 = ModuleMetadataV12{
	Name:       "Module2",
	HasStorage: true,
	Storage:    exampleStorageMetadataV10,
	HasCalls:   true,
	Calls:      []FunctionMetadataV4{exampleFunctionMetadataV4},
	HasEvents:  true,
	Events:     []EventMetadataV4{exampleEventMetadataV4},
	Constants:  []ModuleConstantMetadataV6{exampleModuleConstantMetadataV6},
	Errors:     []ErrorMetadataV8{exampleErrorMetadataV8},
	Index:      2,
}

func TestNewMetadataV12_Decode(t *testing.T) {
	tests := []struct {
		name, hexData string
	}{
		{
			"PolkadotV12", ExamplaryMetadataV12PolkadotString,
		},
	}

	for _, s := range tests {
		t.Run(s.name, func(t *testing.T) {
			metadata := NewMetadataV12()
			err := DecodeFromBytes(MustHexDecodeString(s.hexData), metadata)
			assert.True(t, metadata.IsMetadataV12)
			assert.NoError(t, err)
			data, err := EncodeToBytes(metadata)
			assert.NoError(t, err)
			assert.Equal(t, s.hexData, HexEncodeToString(data))
		})

	}
}

func TestMetadataV12_ExistsModuleMetadata(t *testing.T) {
	assert.True(t, exampleMetadataV12.ExistsModuleMetadata("EmptyModule"))
	assert.False(t, exampleMetadataV12.ExistsModuleMetadata("NotExistModule"))
}

func TestMetadataV12_FindEventNamesForEventID(t *testing.T) {
	module, event, err := exampleMetadataV12.FindEventNamesForEventID(EventID([2]byte{1, 0}))

	assert.NoError(t, err)
	assert.Equal(t, exampleModuleMetadataV121.Name, module)
	assert.Equal(t, exampleEventMetadataV4.Name, event)
}

func TestMetadataV12_FindEventNamesForUnknownModule(t *testing.T) {
	_, _, err := exampleMetadataV12.FindEventNamesForEventID(EventID([2]byte{1, 18}))

	assert.Error(t, err)
}

func TestMetadataV12_TestFindStorageEntryMetadata(t *testing.T) {
	_, err := exampleMetadataV12.FindStorageEntryMetadata("myStoragePrefix", "myStorageFunc2")
	assert.NoError(t, err)
}

func TestMetadataV12_TestFindCallIndex(t *testing.T) {
	callIndex, err := exampleMetadataV12.FindCallIndex("Module2.my function")
	assert.NoError(t, err)
	assert.Equal(t, exampleModuleMetadataV122.Index, callIndex.SectionIndex)
	assert.Equal(t, uint8(0), callIndex.MethodIndex)
}

func TestMetadataV12_TestFindCallIndexWithUnknownModule(t *testing.T) {
	_, err := exampleMetadataV12.FindCallIndex("UnknownModule.my function")
	assert.Error(t, err)
}

func TestMetadataV12_TestFindCallIndexWithUnknownFunction(t *testing.T) {
	_, err := exampleMetadataV12.FindCallIndex("Module2.unknownFunction")
	assert.Error(t, err)
}
