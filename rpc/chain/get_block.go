// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
//
// Copyright 2019 Centrifuge GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package chain

import (
	"github.com/centrifuge/go-substrate-rpc-client/types"
)

func (c *Chain) GetBlock(blockHash types.Hash) (*types.SignedBlock, error) {
	return c.getBlock(&blockHash)
}

func (c *Chain) GetBlockLatest() (*types.SignedBlock, error) {
	return c.getBlock(nil)
}

func (c *Chain) getBlock(blockHash *types.Hash) (*types.SignedBlock, error) {
	var SignedBlock types.SignedBlock
	var err error
	if blockHash == nil {
		err = (*c.client).Call(&SignedBlock, "chain_getBlock")
	} else {
		hexHash, err := types.Hex(*blockHash)
		if err != nil {
			return nil, err
		}
		err = (*c.client).Call(&SignedBlock, "chain_getBlock", hexHash)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	return &SignedBlock, err
}
