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
	"math/big"
	"os"
	"testing"
	"time"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client"
	"github.com/centrifuge/go-substrate-rpc-client/signature"
	"github.com/centrifuge/go-substrate-rpc-client/types"
)

func TestChain_ProcessEvents(t *testing.T) {
	//os.Setenv("TEST_PRIV_KEY", "0xe22b94b0b31792a48c173066cc48cf2c0df646bc990add739430cd5cf99a4c36")
	//from, ok := signature.LoadKeyringPairFromEnv()
	//if !ok {
	//	t.Skip("skipping end-to-end that requires a private key because TEST_PRIV_KEY is not set or empty")
	//}

	//api, err := gsrpc.NewSubstrateAPI("wss://kusama-rpc.polkadot.io/")
	api, err := gsrpc.NewSubstrateAPI("wss://fullnode.flint.centrifuge.io")
	//api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}

	key, err := types.CreateStorageKey(meta, "System", "Events", nil, nil)
	if err != nil {
		panic(err)
	}

	hash, err := api.RPC.Chain.GetBlockHash(1516955)
	//hash, err := api.RPC.Chain.GetBlockHash(3116641)
	if err != nil {
		panic(err)
	}

	var records types.EventRecordsRaw
	fmt.Printf("%x\n", key)
	_, err = api.RPC.State.GetStorage(key, &records, hash)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%x\n", records)

	e := types.EventRecords{}
	err = records.DecodeEventRecords(meta, &e)
	if err != nil {
		panic(err)
	}

}

func TestChain_SubmitExtrinsic(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping end-to-end test in short mode.")
	}

	os.Setenv("TEST_PRIV_KEY", "0xe22b94b0b31792a48c173066cc48cf2c0df646bc990add739430cd5cf99a4c36")
	from, ok := signature.LoadKeyringPairFromEnv()
	if !ok {
		t.Skip("skipping end-to-end that requires a private key because TEST_PRIV_KEY is not set or empty")
	}

	api, err := gsrpc.NewSubstrateAPI("wss://fullnode.fulvous.centrifuge.io")
	//api, err := gsrpc.NewSubstrateAPI(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}

	bob, err := types.NewAddressFromHexAccountID("0x9efc9f132428d21268710181fe4315e1a02d838e0e5239fe45599f54310a7c34")
	if err != nil {
		panic(err)
	}

	c, err := types.NewCall(meta, "Balances.transfer", bob, types.NewUCompact(big.NewInt(10000)))
	if err != nil {
		panic(err)
	}

	ext := types.NewExtrinsic(c)

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

	key, err := types.CreateStorageKey(meta, "System", "Account", from.PublicKey, nil)
	if err != nil {
		panic(err)
	}

	var accountInfo types.AccountInfo
	ok, err = api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil || !ok {
		panic(err)
	}

	nonce := uint32(accountInfo.Nonce)

	for i := uint32(0); i < 4; i++ {
		o := types.SignatureOptions{
			// BlockHash:   blockHash,
			BlockHash:   genesisHash, // BlockHash needs to == GenesisHash if era is immortal. // TODO: add an error?
			Era:         era,
			GenesisHash: genesisHash,
			Nonce:       types.NewUCompact(new(big.Int).SetUint64(uint64(nonce + i))),
			SpecVersion: rv.SpecVersion,
			Tip:         types.NewUCompact(big.NewInt(0)),
		}

		extI := ext

		err = extI.Sign(from, o)
		if err != nil {
			panic(err)
		}

		extEnc, err := types.EncodeToHexString(extI)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Extrinsic: %#v\n", extEnc)

		_, err = api.RPC.Author.SubmitExtrinsic(extI)
		if err != nil {
			panic(err)
		}
	}

	for i := 0; ; i++ {
		xts, err := api.RPC.Author.PendingExtrinsics()
		if err != nil {
			panic(err)
		}

		fmt.Printf("Pending extrinsics: %#v\n", xts)

		if i >= 2 {
			break
		}

		time.Sleep(1 * time.Second)
	}
}
