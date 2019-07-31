package substrate

import (
	"bytes"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/scale"
	"github.com/centrifuge/go-substrate-rpc-client/testrpc"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
)

var testServer *testrpc.Server
var testClient Client

func TestMain(m *testing.M) {
	testServer = new(testrpc.Server)
	testServer.Init()
	var err error
	testClient, err = Connect("ws://localhost:8080")
	if err != nil {
		panic(err)
	}
	m.Run()
	// TODO shutdown server
}

func TestState_GetMetaData(t *testing.T) {
	s := NewStateRPC(testClient)
	res, err := s.MetaData([]byte{})
	assert.NoError(t, err)
	assert.Equal(t, "system", res.Metadata.Modules[0].Name)
}

func TestState_Storage(t *testing.T) {
	s := NewStateRPC(testClient)
	b, _ := hexutil.Decode(AlicePubKey)
	h, _ := hexutil.Decode("0x142d4b3d1946e4956b4bd5a5bfc906142e921b51415ceccb3c82b3bd3ff3daf1")

	m, _ := s.MetaData(h)
	key, err := NewStorageKey(*m, "System", "AccountNonce", b)
	assert.NoError(t, err)
	res, err := s.Storage(key, nil)
	assert.NoError(t, err)

	buf := bytes.NewBuffer(res)
	tempDec := scale.NewDecoder(buf)
	var nonce uint64
	err = tempDec.Decode(&nonce)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0xffffffffffffff1), nonce)
}
