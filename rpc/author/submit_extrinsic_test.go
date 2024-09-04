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

package author_test

import (
	"fmt"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/extrinsic"
	"testing"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/config"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/stretchr/testify/assert"
)

func TestAuthor_SubmitExtrinsic(t *testing.T) {
	// Instantiate the API
	api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
	assert.NoError(t, err)

	meta, err := api.RPC.State.GetMetadataLatest()
	assert.NoError(t, err)

	// Create a call, transferring 12345 units to Bob
	bob, err := types.NewMultiAddressFromHexAccountID("0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48")
	assert.NoError(t, err)

	amount := types.NewUCompactFromUInt(12345)
	c, err := types.NewCall(meta, "Balances.transfer", bob, amount)
	assert.NoError(t, err)

	for {
		// Create the extrinsic
		ext := extrinsic.NewExtrinsic(c)
		genesisHash, err := api.RPC.Chain.GetBlockHash(0)
		assert.NoError(t, err)

		rv, err := api.RPC.State.GetRuntimeVersionLatest()
		assert.NoError(t, err)

		// Get the nonce for Alice
		key, err := types.CreateStorageKey(meta, "System", "Account", signature.TestKeyringPairAlice.PublicKey)
		assert.NoError(t, err)

		var accountInfo types.AccountInfo
		ok, err := api.RPC.State.GetStorageLatest(key, &accountInfo)
		assert.NoError(t, err)
		assert.True(t, ok)
		nonce := uint32(accountInfo.Nonce)

		fmt.Printf("Sending %v from %#x to %#x with nonce %v\n", amount, signature.TestKeyringPairAlice.PublicKey,
			bob.AsID, nonce)

		// Sign the transaction using Alice's default account
		err = ext.Sign(signature.TestKeyringPairAlice, meta, extrinsic.WithEra(types.ExtrinsicEra{IsImmortalEra: true}, genesisHash),
			extrinsic.WithNonce(types.NewUCompactFromUInt(uint64(nonce))),
			extrinsic.WithTip(types.NewUCompactFromUInt(0)),
			extrinsic.WithSpecVersion(rv.SpecVersion),
			extrinsic.WithTransactionVersion(rv.TransactionVersion),
			extrinsic.WithGenesisHash(genesisHash),
		)
		assert.NoError(t, err)

		res, err := api.RPC.Author.SubmitExtrinsic(ext)
		if err != nil {
			t.Logf("extrinsic submit failed: %v", err)
			continue
		}

		hex, err := codec.Hex(res)
		assert.NoError(t, err)
		assert.NotEmpty(t, hex)
		break
	}
}
