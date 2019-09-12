package state

import (
	client2 "github.com/centrifuge/go-substrate-rpc-client/client"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var state *State

func TestMain(m *testing.M) {
	client, err := client2.Connect("ws://127.0.0.1:9944")
	if err != nil {
		panic(err)
	}

	state = NewState(&client)

	os.Exit(m.Run())
}

func TestState_GetMetadata(t *testing.T) {
	_, err := state.GetMetadataLatest()
	assert.NoError(t, err)
	// assert.False(t, res.IsEmpty()) // TODO
}
