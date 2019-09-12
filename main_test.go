package gsrpc_test

import (
	"fmt"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client"
)

func Example() {
	fmt.Println("Hello world")
}

func Example_simpleConnect() {
	api, err := gsrpc.NewSubstrateApi("ws://127.0.0.1:9944")
	if err != nil {
		panic(err)
	}

	fmt.Println(api.RPC.Chain.GetBlockHashLatest())
	//chain := api.RPC.System.Chain()
	//name := api.RPC.System.Name()
	//version := api.RPC.System.Version()

	// Output: 0xa310
}

// listen to new blocks
// func Example_listenToNewBlocks(t *testing.T) {
//	api := NewSubstrateApi("ws://127.0.0.1:9944")
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
