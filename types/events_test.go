package types

import (
	"bytes"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/scale"
	"github.com/stretchr/testify/assert"
)

func TestVoteThreshold_Decoder(t *testing.T) {
	// SuperMajorityAgainst
	decoder := scale.NewDecoder(bytes.NewReader([]byte{1}))
	vt := VoteThreshold(0)
	err := decoder.Decode(&vt)
	assert.NoError(t, err)
	assert.Equal(t, vt, SuperMajorityAgainst)

	// Error
	decoder = scale.NewDecoder(bytes.NewReader([]byte{3}))
	err = decoder.Decode(&vt)
	assert.Error(t, err)
}

func TestVoteThreshold_Encode(t *testing.T) {
	vt := SuperMajorityAgainst
	var buf bytes.Buffer
	encoder := scale.NewEncoder(&buf)
	assert.NoError(t, encoder.Encode(vt))
	assert.Equal(t, buf.Len(), 1)
	assert.Equal(t, buf.Bytes(), []byte{1})
}

func TestDispatchResult_Decode(t *testing.T) {
	// ok
	decoder := scale.NewDecoder(bytes.NewReader([]byte{0}))
	var res DispatchResult
	err := decoder.Decode(&res)
	assert.NoError(t, err)
	assert.True(t, res.Ok)

	// Dispatch Error
	decoder = scale.NewDecoder(bytes.NewReader([]byte{1, 3, 1, 1}))
	res = DispatchResult{}
	assert.NoError(t, decoder.Decode(&res))
	assert.False(t, res.Ok)
	assert.True(t, res.Error.HasModule)
	assert.Equal(t, res.Error.Module, byte(1))
	assert.Equal(t, res.Error.Error, byte(1))

	// decoder error
	decoder = scale.NewDecoder(bytes.NewReader([]byte{1, 3, 1}))
	res = DispatchResult{}
	assert.Error(t, decoder.Decode(&res))
}
