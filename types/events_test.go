package types

import (
	"bytes"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v3/scale"
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

func TestProxyTypeEncodeDecode(t *testing.T) {
	// encode
	pt := Governance
	var buf bytes.Buffer
	encoder := scale.NewEncoder(&buf)
	assert.NoError(t, encoder.Encode(pt))
	assert.Equal(t, buf.Len(), 1)
	assert.Equal(t, buf.Bytes(), []byte{2})

	//decode
	decoder := scale.NewDecoder(bytes.NewReader(buf.Bytes()))
	pt0 := ProxyType(0)
	err := decoder.Decode(&pt0)
	assert.NoError(t, err)
	assert.Equal(t, pt0, Governance)

	//decode error
	decoder = scale.NewDecoder(bytes.NewReader([]byte{5}))
	pt0 = ProxyType(0)
	err = decoder.Decode(&pt0)
	assert.Error(t, err)
}

func TestPaysEncodeDecode(t *testing.T) {
	// encode
	pt := Pays{IsNo: true}
	var buf bytes.Buffer
	encoder := scale.NewEncoder(&buf)
	assert.NoError(t, encoder.Encode(pt))
	assert.Equal(t, buf.Len(), 1)
	assert.Equal(t, buf.Bytes(), []byte{1})

	// decode unsupported
	var pays Pays
	decoder := scale.NewDecoder(bytes.NewReader([]byte{2}))
	err := decoder.Decode(&pays)
	assert.Error(t, err)

	// decode supported
	decoder = scale.NewDecoder(bytes.NewReader(buf.Bytes()))
	err = decoder.Decode(&pays)
	assert.NoError(t, err)
	assert.True(t, pays.IsNo)
}

func TestDispatchClassEncodeDecode(t *testing.T) {
	// encode
	dc := DispatchClass{IsMandatory: true}
	var buf bytes.Buffer
	encoder := scale.NewEncoder(&buf)
	assert.NoError(t, encoder.Encode(dc))
	assert.Equal(t, buf.Len(), 1)
	assert.Equal(t, buf.Bytes(), []byte{2})

	// decode unsupported
	var dcc DispatchClass
	decoder := scale.NewDecoder(bytes.NewReader([]byte{3}))
	err := decoder.Decode(&dcc)
	assert.Error(t, err)

	// decode supported
	decoder = scale.NewDecoder(bytes.NewReader(buf.Bytes()))
	err = decoder.Decode(&dcc)
	assert.NoError(t, err)
	assert.True(t, dcc.IsMandatory)
}
