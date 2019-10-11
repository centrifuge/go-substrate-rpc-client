// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
//
// Copyright 2019 Centrifuge GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/centrifuge/go-substrate-rpc-client/types"
)

var exampleMetadata = Metadata{
	MagicNumber: 13121231,
	Version:     4,
	Metadata:    exampleRuntimeMetadataV4,
}

var exampleRuntimeMetadataV4 = RuntimeMetadataV4{
	Modules: []ModuleMetadata{exampleModuleMetadata},
}

var exampleCallIndex = CallIndex{
	SectionIndex: 123,
	MethodIndex:  234,
}

var exampleModuleMetadata = ModuleMetadata{
	Name:       "myModule",
	Prefix:     "modulePrefix",
	HasStorage: true,
	Storage:    []StorageFunctionMetadata{exampleStorageFunctionMetadata},
	HasCalls:   true,
	Calls:      []FunctionMetadata{exampleFunctionMetadata},
	HasEvents:  true,
	Events:     []EventMetadata{exampleEventMetadata},
}

var exampleStorageFunctionMetadata = StorageFunctionMetadata{
	Name:          "myStorageFunc",
	Modifier:      3,
	Type:          2,
	Plane:         "",
	Map:           TypMap{},
	DMap:          exampleTypDoubleMap,
	Fallback:      []byte{23, 14},
	Documentation: []string{"My", "storage func", "doc"},
}

var exampleFunctionMetadata = FunctionMetadata{
	Name:          "my function",
	Args:          []FunctionArgumentMetadata{exampleFunctionArgumentMetadata},
	Documentation: []string{"My", "doc"},
}

var exampleEventMetadata = EventMetadata{
	Name:          "myEvent",
	Args:          []string{"arg1", "arg2"},
	Documentation: []string{"My", "doc"},
}

var exampleTypMap = TypMap{
	Hasher:   13,
	Key:      "my key",
	Value:    "and my value",
	IsLinked: false,
}

var exampleTypDoubleMap = TypDoubleMap{
	Hasher:     4,
	Key:        "myKey",
	Key2:       "otherKey",
	Value:      "and a value",
	Key2Hasher: "and a hasher",
}

var exampleFunctionArgumentMetadata = FunctionArgumentMetadata{Name: "myFunctionName", Type: "myType"}

func TestMetadata_Decode(t *testing.T) {
	metadata := NewMetadata()

	err := DecodeFromBytes(MustHexDecodeString(ExamplaryMetadataV4String), metadata)
	assert.NoError(t, err)

	assert.Equal(t, ExamplaryMetadataV4, *metadata)
	// assert.Equal(t, decodedMetadata, fmt.Sprintf("%v", metadata))
}

func TestMetadata_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, exampleMetadata)
}

func TestRuntimeMetadataV4_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, exampleRuntimeMetadataV4)
}

func TestCallIndex_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, exampleCallIndex)
}

func TestModuleMetadata_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, exampleModuleMetadata)
}

func TestStorageFunctionMetadata_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, exampleStorageFunctionMetadata)
}

func TestFunctionMetadata_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, exampleFunctionMetadata)
}

func TestEventMetadata_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, exampleEventMetadata)
}

func TestTypMap_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, exampleTypMap)
}

func TestTypDoubleMap_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, exampleTypDoubleMap)
}

func TestFunctionArgumentMetadata_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, exampleFunctionArgumentMetadata)
}
