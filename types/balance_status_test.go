package types

import (
	"bytes"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/scale"
	"github.com/stretchr/testify/assert"
)

func TestBalanceStatusEncodeDecode(t *testing.T) {
	// encode
	bs := Reserved
	var buf bytes.Buffer
	encoder := scale.NewEncoder(&buf)
	assert.NoError(t, encoder.Encode(bs))
	assert.Equal(t, buf.Len(), 1)
	assert.Equal(t, buf.Bytes(), []byte{1})

	//decode
	decoder := scale.NewDecoder(bytes.NewReader(buf.Bytes()))
	bs0 := BalanceStatus(0)
	err := decoder.Decode(&bs0)
	assert.NoError(t, err)
	assert.Equal(t, bs0, Reserved)

	//decode error
	decoder = scale.NewDecoder(bytes.NewReader([]byte{5}))
	bs0 = BalanceStatus(0)
	err = decoder.Decode(&bs0)
	assert.Error(t, err)
}
