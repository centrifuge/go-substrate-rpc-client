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
	"fmt"
	"testing"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client"
	"github.com/centrifuge/go-substrate-rpc-client/config"
	"github.com/centrifuge/go-substrate-rpc-client/signature"
	"github.com/centrifuge/go-substrate-rpc-client/types"
)

func TestState_SubscribeExtrinsic(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping end-to-end test in short mode.")
	}

	api, err := gsrpc.NewSubstrateAPI(config.Get().RPCURL)
	if err != nil {
		panic(err)
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}

	bob, err := types.NewAddressFromHexAccountID("0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48")
	if err != nil {
		panic(err)
	}

	c, err := types.NewCall(meta, "Balances.transfer", bob, types.UCompact(6969))
	if err != nil {
		panic(err)
	}

	ext := types.NewExtrinsic(c)
	if err != nil {
		panic(err)
	}

	// blockHash, err := api.RPC.Chain.GetBlockHashLatest()
	// if err != nil {
	// 	panic(err)
	// }

	era := types.ExtrinsicEra{IsMortalEra: false}

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		panic(err)
	}

	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		panic(err)
	}

	key, err := types.CreateStorageKey(meta, "System", "AccountNonce", signature.TestKeyringPairAlice.PublicKey)
	if err != nil {
		panic(err)
	}

	var nonce uint32
	err = api.RPC.State.GetStorageLatest(key, &nonce)
	if err != nil {
		panic(err)
	}

	o := types.SignatureOptions{
		// BlockHash:   blockHash,
		BlockHash:   genesisHash, // BlockHash needs to == GenesisHash if era is immortal. // TODO: add an error?
		Era:         era,
		GenesisHash: genesisHash,
		Nonce:       types.UCompact(nonce),
		SpecVersion: rv.SpecVersion,
		Tip:         0,
	}

	err = ext.Sign(signature.TestKeyringPairAlice, o)
	if err != nil {
		panic(err)
	}

	extEnc, err := types.EncodeToHexString(ext)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", extEnc)
}
