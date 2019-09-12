package chain

import (
	"encoding/hex"
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client"
	"github.com/centrifuge/go-substrate-rpc-client/types"
)

type Chain struct {
	client *substrate.Client
}

func NewChain(c *substrate.Client) *Chain {
	return &Chain{c}
}

func (c *Chain) GetBlockHash(blockNumber uint64) (types.Hash, error) {
	return c.getBlockHash(&blockNumber)
}

func (c *Chain) GetBlockHashLatest() (types.Hash, error) {
	return c.getBlockHash(nil)
}

func (c *Chain) getBlockHash(blockNumber *uint64) (types.Hash, error) {
	var res string
	var err error

	if blockNumber == nil {
		err = (*c.client).Call(&res, "chain_getBlockHash")
	} else {
		err = (*c.client).Call(&res, "chain_getBlockHash", *blockNumber)
	}

	if err != nil {
		return types.Hash{}, err
	}

	bz, err := hex.DecodeString(res[2:])
	if err != nil {
		return types.Hash{}, err
	}

	if len(bz) != 32 {
		return types.Hash{}, fmt.Errorf("Required result to be 32 bytes, but got %v", len(bz))
	}

	var bz32 [32]byte
	copy(bz32[:], bz)

	hash := types.NewHash(bz32)

	return hash, nil
}
