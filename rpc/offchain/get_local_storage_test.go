package offchain

import (
	"crypto/rand"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
	"github.com/stretchr/testify/assert"
)

func TestOffchain_LocalStorageGetSet(t *testing.T) {
	key := make([]byte, 20)
	n, err := rand.Read(key)
	assert.NoError(t, err)
	assert.Equal(t, 20, n)

	data, err := offchain.LocalStorageGet(Persistent, key)
	assert.NoError(t, err)
	assert.Empty(t, data)

	err = offchain.LocalStorageSet(Persistent, key, key)
	assert.NoError(t, err)

	data, err = offchain.LocalStorageGet(Persistent, key)
	assert.NoError(t, err)

	got := make([]byte, 0, 20)
	err = types.DecodeFromHexString(data.Hex(), &got)
	assert.NoError(t, err)
	assert.Equal(t, key, got)
}
