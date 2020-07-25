package blake2b_test

import (
	"testing"

	. "github.com/Snowfork/go-substrate-rpc-client/blake2b"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/blake2b"
)

func Test_128Concat(t *testing.T) {
	key := []byte("abc")

	h128, _ := blake2b.New(16, nil)
	h128.Write(key)

	h128Concat, _ := New128Concat()
	h128Concat.Write(key)
	assert.Equal(t, append(h128.Sum(nil), key...), h128Concat.Sum(nil))
}
