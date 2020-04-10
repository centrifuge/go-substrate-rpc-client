package types_test

import (
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/scale"
	. "github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/stretchr/testify/assert"
)

func TestNewMetadataV11_Decode(t *testing.T) {
	tests := []struct{
		name, hexData string
	}{
		{
			"SubstrateV11", ExamplaryMetadataV11SubstrateString,
		},

		{
			"PolkadotV11", ExamplaryMetadataV11PolkadotString,
		},
	}

	for _, s := range tests{
		t.Run(s.name, func(t *testing.T) {
			metadata := NewMetadataV11()
			err := DecodeFromBytes(MustHexDecodeString(s.hexData), metadata)
			assert.NoError(t, err)
			data, err := EncodeToBytes(metadata, scale.EncoderOptions{})
			assert.NoError(t, err)
			assert.Equal(t, s.hexData, HexEncodeToString(data))
		})

	}
}

