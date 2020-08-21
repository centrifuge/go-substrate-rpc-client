package types_test

import (
	"testing"

	. "github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/stretchr/testify/assert"
)

var exampleMetadataV11 = Metadata{
	MagicNumber:   0x6174656d,
	Version:       11,
	IsMetadataV11: true,
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
			err := DecodeFromBytes(MustHexDecodeString(s.hexData), metadata)
			assert.NoError(t, err)
			data, err := EncodeToBytes(metadata)
			assert.NoError(t, err)
			assert.Equal(t, s.hexData, HexEncodeToString(data))
		})

	}
}

func TestMetadataV11_ExistsModuleMetadata(t *testing.T) {
	assert.True(t, exampleMetadataV11.ExistsModuleMetadata("EmptyModule"))
	assert.False(t, exampleMetadataV11.ExistsModuleMetadata("NotExistModule"))
}
