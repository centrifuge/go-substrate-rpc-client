package state

import (
	"os"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/client"
	"github.com/stretchr/testify/assert"
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

func TestState_GetMetadataLatest(t *testing.T) {
	metadata, err := state.GetMetadataLatest()
	assert.NoError(t, err)
	assert.Equal(t, "system", metadata.Metadata.Modules[0].Name)
}

// TODO make test dynamic
//func TestState_GetMetadata(t *testing.T) {
//	bz, _ := hex.DecodeString("cc9ea640d4d4f4dd260b1cbb65cb275df995b056710265f9becdd7e6e1a7b9e0")
//
//	var bz32 [32]byte
//	copy(bz32[:], bz)
//
//	metadata, err := state.GetMetadata(types.NewHash(bz32))
//	assert.NoError(t, err)
//	assert.Equal(t, "system", metadata.Metadata.Modules[0].Name)
//}
