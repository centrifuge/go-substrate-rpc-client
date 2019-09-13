package state

import (
	"os"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/client"
)

var state *State

func TestMain(m *testing.M) {
	// FIXME: due to a size limit, websocket connections don't work with getMetadata right now.
	// 	Related issue: https://github.com/ethereum/go-ethereum/issues/16846
	//  Should get fixed with https://github.com/ethereum/go-ethereum/pull/19866 , released in 1.9.1
	//cl, err := client.Connect("ws://127.0.0.1:9944")
	cl, err := client.Connect("http://127.0.0.1:9933")
	if err != nil {
		panic(err)
	}

	state = NewState(&cl)

	os.Exit(m.Run())
}
