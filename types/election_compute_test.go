package types_test

import (
	"bytes"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	. "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

func TestOptionElectionCompute_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewOptionElectionCompute(NewElectionCompute(byte(0))))
	assertRoundtrip(t, NewOptionElectionComputeEmpty())
}

func TestOptionElectionCompute_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{NewOptionElectionCompute(NewElectionCompute(byte(0))), MustHexDecodeString("0x0100")},
		{NewOptionElectionCompute(NewElectionCompute(byte(1))), MustHexDecodeString("0x0101")},
		{NewOptionElectionCompute(NewElectionCompute(byte(2))), MustHexDecodeString("0x0102")},
		{NewOptionBytesEmpty(), MustHexDecodeString("0x00")},
	})
}

func TestOptionElectionCompute_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{MustHexDecodeString("0x0100"), NewOptionElectionCompute(NewElectionCompute(byte(0)))},
		{MustHexDecodeString("0x0101"), NewOptionElectionCompute(NewElectionCompute(byte(1)))},
		{MustHexDecodeString("0x0102"), NewOptionElectionCompute(NewElectionCompute(byte(2)))},
		{MustHexDecodeString("0x00"), NewOptionBytesEmpty()},
	})
}

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
