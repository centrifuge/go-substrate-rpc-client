package blake2b_test

import (
	"testing"

	. "github.com/Snowfork/go-substrate-rpc-client/blake2b"
	"github.com/stretchr/testify/assert"
)

func Test_128Concat(t *testing.T) {
	key := []byte("abc")

	h128, _ := New128(nil)
	h128.Write(key)
	h128Concat, _ := New128Concat(nil)
	h128Concat.Write(key)

	assert.Equal(t, append(h128.Sum(nil), key...), h128Concat.Sum(nil))
}

func Test_128Concat_MAC(t *testing.T) {
	key := []byte("abc")

	h128, _ := New128(key)
	h128Concat, _ := New128Concat(key)

	assert.Equal(t, append(h128.Sum(nil), key...), h128Concat.Sum(nil))
}
