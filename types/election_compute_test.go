package types

import (
	"bytes"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/scale"
	"github.com/stretchr/testify/assert"
)

func TestElectionComputeEncodeDecode(t *testing.T) {
	// encode
	ec := OnChain
	var buf bytes.Buffer
	encoder := scale.NewEncoder(&buf)
	assert.NoError(t, encoder.Encode(ec))
	assert.Equal(t, buf.Len(), 1)
	assert.Equal(t, buf.Bytes(), []byte{0})

	//decode
	decoder := scale.NewDecoder(bytes.NewReader(buf.Bytes()))
	ec0 := ElectionCompute(0)
	err := decoder.Decode(&ec0)
	assert.NoError(t, err)
	assert.Equal(t, ec0, OnChain)

	//decode error
	decoder = scale.NewDecoder(bytes.NewReader([]byte{5}))
	ec0 = ElectionCompute(0)
	err = decoder.Decode(&ec0)
	assert.Error(t, err)
}
