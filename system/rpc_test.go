// +build tests

package system

import (
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client"
	"github.com/centrifuge/go-substrate-rpc-client/testrpc"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
)

var testServer *testrpc.Server
var testClient substrate.Client
var rpcURL string

func TestMain(m *testing.M) {
	testServer = new(testrpc.Server)
	var err error
	if rpcURL == "" {
		rpcURL, err = testServer.Init(testrpc.GetTestMetaData(), nil)
		if err != nil {
			panic(err)
		}
	}

	testClient, err = substrate.Connect(rpcURL)
	if err != nil {
		panic(err)
	}
	m.Run()
}

func TestBlockHash(t *testing.T) {
	testServer.AddStorageKey("0xa8e78ad25e03ac0281ec709fd3f128efb7e112239d0a7c3e1c86375109bff334", "0xa8e78ad25e03ac0281ec709fd3f128efb7e112239d0a7c3e1c86375109bff338")
	h, err := BlockHash(testClient, 0)
	assert.NoError(t, err)
	assert.Equal(t, hexutil.Encode(h), "0xa8e78ad25e03ac0281ec709fd3f128efb7e112239d0a7c3e1c86375109bff338")
}
