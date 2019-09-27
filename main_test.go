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

package gsrpc_test

import (
	"fmt"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client"
	"github.com/centrifuge/go-substrate-rpc-client/config"
)

func Example_simpleConnect() {
	api, err := gsrpc.NewSubstrateAPI(config.NewDefaultConfig().RPCURL)
	if err != nil {
		panic(err)
	}

	hash, err := api.RPC.Chain.GetBlockHashLatest()
	if err != nil {
		panic(err)
	}

	fmt.Println(hash.Hex())
}

// TODO: add example for listening to new blocks
// func Example_listenToNewBlocks() {
// 	api, err := gsrpc.NewSubstrateAPI(config.NewDefaultConfig().RPCURL)
// 	if err != nil {
// 		panic(err)
// 	}

// 	heads, errs, close, err := api.RPC.System.SubscribeNewHead()
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer close()

// 	// see https://godoc.org/github.com/ethereum/go-ethereum/rpc for more details

// 	count := 0

// 	for {
// 		select {
// 		case head := <-heads:
// 			fmt.Printf("#%v: Got header %v\n", count, head.Number)
// 			count++
// 		case err := <-errs:
// 			fmt.Errorf("Got error: %v;", err)
// 		}
// 	}
// }

//// TODO: implement more: https://polkadot.js.org/api/examples/promise/03_listen_to_balance_change/
