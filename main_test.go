// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
// Copyright (C) 2019  Philip Stanislaus, Philip Stehlik, Vimukthi Wickramasinghe
//
// This file is part of Go Substrate RPC Client (GSRPC).
//
// GSRPC is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// GSRPC is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package gsrpc_test

import (
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/config"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client"
)

func Example_simpleConnect() {
	api, err := gsrpc.NewSubstrateAPI(config.NewDefaultConfig().RPCURL)
	if err != nil {
		panic(err)
	}

	fmt.Println(api.RPC.Chain.GetBlockHashLatest())
	//chain := api.RPC.System.Chain()
	//name := api.RPC.System.Name()
	//version := api.RPC.System.Version()

	// TODO Output: 0x26735e602578d1f1741552e32d074a5fdefde5bf5e1ad074ba0ba6a96c9cc014
}

// listen to new blocks
// func Example_listenToNewBlocks(t *testing.T) {
//	api := NewSubstrateAPI("ws://127.0.0.1:9944")
//
//	heads, errs, close, err := api.RPC.System.SubscribeNewHead()
//	if err != nil {
//		panic(err)
//	}
//	defer close()
//
//	// see https://godoc.org/github.com/ethereum/go-ethereum/rpc for more details
//
//	count := 0
//
//	for {
//		select {
//		case head := <-heads:
//			fmt.Printf("#%v: Got header %v\n", count, head.Number)
//			count++
//		case err := <-errs:
//			fmt.Errorf("Got error: %v;", err)
//		}
//	}
//}
//
//// implement more: https://polkadot.js.org/api/examples/promise/03_listen_to_balance_change/
