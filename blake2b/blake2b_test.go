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
	h128.Write(key)
	h128Concat, _ := New128Concat(nil)
	h128Concat.Write(key)
	h128Concat.Write(key)

	assert.Equal(t, append(h128.Sum(nil), append(key, key...)...), h128Concat.Sum(nil))
}

func Test_128Concat_MAC(t *testing.T) {
	key := []byte("abc")

	h128, _ := New128(key)
	h128Concat, _ := New128Concat(key)

	assert.Equal(t, append(h128.Sum(nil), key...), h128Concat.Sum(nil))
}

func Test_128Concat_Size(t *testing.T) {
	key := []byte("abc")

	h128, _ := New128(nil)
	h128.Write(key)

	h128Concat, _ := New128Concat(nil)
	h128Concat.Write(key)

	assert.Equal(t, h128.Size()+len(key), h128Concat.Size())
}

func Test_128Concat_Reset(t *testing.T) {
	key := []byte("abc")

	h128, _ := New128(nil)
	h128Concat, _ := New128Concat(nil)
	h128Concat.Write(key)

	h128Concat.Reset()

	assert.Equal(t, h128.Sum(nil), h128Concat.Sum(nil))
}

func Test_128Concat_BlockSize(t *testing.T) {
	key := []byte("abc")

	h128, _ := New128(nil)
	h128Concat, _ := New128Concat(nil)
	h128Concat.Write(key)

	assert.Equal(t, h128.BlockSize(), h128Concat.BlockSize())
}
