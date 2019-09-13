package chain

import (
	"github.com/centrifuge/go-substrate-rpc-client/config"
	"os"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/types"

	"github.com/centrifuge/go-substrate-rpc-client/client"
	"github.com/stretchr/testify/assert"
)

var chain *Chain

func TestMain(m *testing.M) {
	cl, err := client.Connect(config.NewDefaultConfig().RPCURL)
	if err != nil {
		panic(err)
	}

	chain = NewChain(&cl)

	os.Exit(m.Run())
}

func TestChain_GetBlockHash(t *testing.T) {
	res, err := chain.GetBlockHash(1)
	assert.NoError(t, err)
	hex, err := types.Hex(res)
	assert.NoError(t, err)
	assert.True(t, len(hex) > 0)
}

func TestChain_GetBlockHashLatest(t *testing.T) {
	res, err := chain.GetBlockHashLatest()
	assert.NoError(t, err)
	hex, err := types.Hex(res)
	assert.NoError(t, err)
	assert.True(t, len(hex) > 0)
}
