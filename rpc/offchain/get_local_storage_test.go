package offchain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOffchain_LocalStorageGet(t *testing.T) {
	data, err := offchain.LocalStorageGet(Persistent, []byte("foo"))
	assert.NoError(t, err)
	assert.Equal(t, mockSrv.storageValueHex, data.Hex())
}
