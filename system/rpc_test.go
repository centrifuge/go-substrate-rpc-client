package system

import (
	"fmt"
	"github.com/centrifuge/go-substrate-rpc-client"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBlockHash(t *testing.T) {
	c, _ := substrate.Connect("ws://127.0.0.1:9944")
	h, err := BlockHash(c, 0)
	assert.NoError(t, err)

	fmt.Printf("%s", hexutil.Encode(h))
}
