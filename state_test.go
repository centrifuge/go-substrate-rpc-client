package substrate

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/scale"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
)

func TestState_GetMetaData(t *testing.T) {
	s := State{nonetwork: true}
	res, err := s.MetaData([]byte{})
	assert.NoError(t, err)
	assert.Equal(t, "system", res.Metadata.Modules[0].Name)
	// fmt.Println(res)
}

func TestState_Storage(t *testing.T) {
	t.SkipNow()
	c, _ := Connect("ws://127.0.0.1:9944")
	s := NewStateRPC(c)
	b, _ := hexutil.Decode(AlicePubKey)
	h, _ := hexutil.Decode("0x142d4b3d1946e4956b4bd5a5bfc906142e921b51415ceccb3c82b3bd3ff3daf1")

	m, _ := s.MetaData(h)
	key, _ := NewStorageKey(*m, "System", "AccountNonce", b)
	res, _ := s.Storage(key, nil)

	buf := bytes.NewBuffer(res)
	tempDec := scale.NewDecoder(buf)
	var nonce uint64
	tempDec.Decode(&nonce)
	fmt.Println(nonce)
}
