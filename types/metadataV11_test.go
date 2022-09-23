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

	. "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	. "github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/stretchr/testify/assert"
)

var exampleMetadataV11 = Metadata{
	MagicNumber:   0x6174656d,
	Version:       11,
	AsMetadataV11: MetadataV11{MetadataV10: exampleRuntimeMetadataV10},
}

func TestNewMetadataV11_Decode(t *testing.T) {
	tests := []struct {
		name, hexData string
	}{
		{
			"SubstrateV11", ExamplaryMetadataV11SubstrateString,
		},

		{
			"PolkadotV11", ExamplaryMetadataV11PolkadotString,
		},
	}

	for _, s := range tests {
		t.Run(s.name, func(t *testing.T) {
			metadata := NewMetadataV11()
			err := Decode(MustHexDecodeString(s.hexData), metadata)
			assert.NoError(t, err)
			data, err := Encode(metadata)
			assert.NoError(t, err)
			assert.Equal(t, s.hexData, HexEncodeToString(data))
		})

	}
}

func TestMetadataV11_ExistsModuleMetadata(t *testing.T) {
	assert.True(t, exampleMetadataV11.ExistsModuleMetadata("EmptyModule"))
	assert.False(t, exampleMetadataV11.ExistsModuleMetadata("NotExistModule"))
}
