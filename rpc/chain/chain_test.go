package chain

import (
	"os"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client"
	"github.com/stretchr/testify/assert"
)

var chain *Chain

func TestMain(m *testing.M) {
	client, err := substrate.Connect("ws://127.0.0.1:9944")
	if err != nil {
		panic(err)
	}

	chain = NewChain(&client)

	os.Exit(m.Run())
}

func TestChain_GetBlockHash(t *testing.T) {
	res, err := chain.GetBlockHash(1)
	assert.NoError(t, err)
	assert.False(t, res.IsEmpty())
}

func TestChain_GetBlockHashLatest(t *testing.T) {
	res, err := chain.GetBlockHashLatest()
	assert.NoError(t, err)
	assert.False(t, res.IsEmpty())
}
