package types

import (
	"hash"

	"golang.org/x/crypto/blake2b"
)

type blake2b128concat struct {
	hasher hash.Hash

	// the key we need to append to the hash digest
	key []byte
}

func (bh blake2b128concat) Write(p []byte) (n int, err error) {
	// save the key for later use in Sum()
	bh.key = p
	return bh.hasher.Write(p)
}

func (bh blake2b128concat) Sum(b []byte) []byte {
	// append key to final hash digest
	return append(b, append(bh.hasher.Sum(nil), bh.key...)...)
}

func (bh blake2b128concat) Reset() {
	bh.hasher.Reset()
}

func (bh blake2b128concat) Size() int {
	return bh.hasher.Size()
}

func (bh blake2b128concat) BlockSize() int {
	return bh.hasher.BlockSize()
}

func Blake2b128ConcatNew() (hash.Hash, error) {
	inner, err := blake2b.New(16, nil)
	if err != nil {
		return nil, err
	}

	hasher := blake2b128concat{
		hasher: inner,
		key:    nil,
	}

	return hasher, nil
}
