package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var exampleRuntimeVersion = RuntimeVersion{
	APIs:             []RuntimeVersionAPI{exampleRuntimeVersionAPI},
	AuthoringVersion: 13,
	ImplName:         "My impl",
	ImplVersion:      21,
	SpecName:         "My spec",
	SpecVersion:      39,
}

var exampleRuntimeVersionAPI = RuntimeVersionAPI{
	APIID:   "0x37e397fc7c91f5e4",
	Version: 23,
}

func TestRuntimeVersion_Encode_Decode(t *testing.T) {
	enc, err := EncodeToBytes(exampleRuntimeVersion)
	assert.NoError(t, err)

	var output RuntimeVersion
	err = DecodeFromBytes(enc, &output)
	assert.NoError(t, err)

	assert.Equal(t, exampleRuntimeVersion, output)
}

func TestRuntimeVersionAPI_Encode_Decode(t *testing.T) {
	enc, err := EncodeToBytes(exampleRuntimeVersionAPI)
	assert.NoError(t, err)

	var output RuntimeVersionAPI
	err = DecodeFromBytes(enc, &output)
	assert.NoError(t, err)

	assert.Equal(t, exampleRuntimeVersionAPI, output)
}
