// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
// Copyright (C) 2019  Centrifuge GmbH
//
// This file is part of Go Substrate RPC Client (GSRPC).
//
// GSRPC is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// GSRPC is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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

var exampleMethodIDX = MethodIDX{
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

	err := DecodeFromBytes(mustDecodeHexString(ExamplaryMetadataV4String), metadata)
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

func TestMethodIDX_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, exampleMethodIDX)
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
