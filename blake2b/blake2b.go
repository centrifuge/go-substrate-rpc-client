package blake2b

import (
	"hash"

	"golang.org/x/crypto/blake2b"
)

type blake2b128concat struct {
	// the underlying blake2b_128 hasher
	hasher hash.Hash

	// the key we need to append to the finalized hash digest (see Sum())
	key []byte
}

func (bh *blake2b128concat) Write(p []byte) (n int, err error) {
	// save the key for later use in Sum()
	bh.key = append(bh.key, p...)

	return bh.hasher.Write(p)
}

func (bh *blake2b128concat) Sum(b []byte) []byte {
	// append key to final hash digest
	return append(b, append(bh.hasher.Sum(nil), bh.key...)...)
}

func (bh *blake2b128concat) Reset() {
	bh.hasher.Reset()
}

func (bh *blake2b128concat) Size() int {
	return bh.hasher.Size()
}

func (bh *blake2b128concat) BlockSize() int {
	return bh.hasher.BlockSize()
}

func New128(key []byte) (hash.Hash, error) {
	return blake2b.New(16, key)
}

func New128Concat(key []byte) (hash.Hash, error) {
	inner, err := blake2b.New(16, key)
	if err != nil {
		return nil, err
	}

	hasher := blake2b128concat{
		hasher: inner,
		key:    key,
	}

	return &hasher, nil
}

func New256(key []byte) (hash.Hash, error) {
	return blake2b.New256(key)
}
