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

package teste2e

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/extrinsic"
	"testing"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/config"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

func TestChain_Events(t *testing.T) {
	targetURL := config.Default().RPCURL // Replace with desired endpoint
	api, err := gsrpc.NewSubstrateAPI(targetURL)
	assert.NoError(t, err)

	meta, err := api.RPC.State.GetMetadataLatest()
	assert.NoError(t, err)

	key, err := types.CreateStorageKey(meta, "System", "Events", nil)
	assert.NoError(t, err)

	blockNUmber := uint64(0) // Replace with desired block to parse events
	bh, err := api.RPC.Chain.GetBlockHash(blockNUmber)
	assert.NoError(t, err)

	raw, err := api.RPC.State.GetStorageRaw(key, bh)
	assert.NoError(t, err)

	events := types.EventRecords{}
	err = types.EventRecordsRaw(*raw).DecodeEventRecords(meta, &events)
	assert.NoError(t, err)
}

func TestChain_SubmitExtrinsic(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping end-to-end test in short mode.")
	}

	from, ok := signature.LoadKeyringPairFromEnv()
	if !ok {
		from = signature.TestKeyringPairAlice
	}

	api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
	assert.NoError(t, err)

	meta, err := api.RPC.State.GetMetadataLatest()
	assert.NoError(t, err)

	bob, err := types.NewMultiAddressFromHexAccountID("0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48")
	assert.NoError(t, err)

	c, err := types.NewCall(meta, "Balances.transfer", bob, types.NewUCompactFromUInt(6969))
	assert.NoError(t, err)

	ext := extrinsic.NewExtrinsic(c)

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	assert.NoError(t, err)

	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	assert.NoError(t, err)

	key, err := types.CreateStorageKey(meta, "System", "Account", from.PublicKey)
	assert.NoError(t, err)

	var accountInfo types.AccountInfo
	ok, err = api.RPC.State.GetStorageLatest(key, &accountInfo)
	assert.NoError(t, err)
	assert.True(t, ok)

	nonce := uint32(accountInfo.Nonce)
	var txn types.Hash
	for {
		extI := ext

		err = extI.Sign(from, meta, extrinsic.WithEra(types.ExtrinsicEra{IsImmortalEra: true}, genesisHash),
			extrinsic.WithNonce(types.NewUCompactFromUInt(uint64(nonce))),
			extrinsic.WithTip(types.NewUCompactFromUInt(0)),
			extrinsic.WithSpecVersion(rv.SpecVersion),
			extrinsic.WithTransactionVersion(rv.TransactionVersion),
			extrinsic.WithGenesisHash(genesisHash),
		)
		assert.NoError(t, err)

		txn, err = api.RPC.Author.SubmitExtrinsic(extI)
		if err != nil {
			nonce++
			continue
		}

		break
	}
	assert.NotEmpty(t, txn)
}
