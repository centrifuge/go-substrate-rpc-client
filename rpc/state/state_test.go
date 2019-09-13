package state

import (
	"os"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/client"
	"github.com/centrifuge/go-substrate-rpc-client/config"
)

var state *State

func TestMain(m *testing.M) {
	cl, err := client.Connect(config.NewDefaultConfig().RPCURL)
	if err != nil {
		panic(err)
	}

	state = NewState(&cl)

	os.Exit(m.Run())
}
