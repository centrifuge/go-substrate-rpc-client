package offchain

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOffchain_LocalStorageGetSet(t *testing.T) {
	key := make([]byte, 20)
	n, err := rand.Read(key)
	assert.NoError(t, err)
	assert.Equal(t, 20, n)

	value := []byte{0, 1, 2}

	data, err := offchain.LocalStorageGet(Persistent, key)
	assert.NoError(t, err)
	assert.Empty(t, data)

	err = offchain.LocalStorageSet(Persistent, key, value)
	assert.NoError(t, err)

	data, err = offchain.LocalStorageGet(Persistent, key)
	assert.NoError(t, err)

	assert.Equal(t, value, []byte(*data))
}
