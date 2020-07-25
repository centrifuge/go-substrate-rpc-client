package blake2b

import (
	"hash"

	"golang.org/x/crypto/blake2b"
)

type state struct {
	// the underlying blake2b_128 hasher
	hasher hash.Hash

	// the key we need to append to the finalized hash digest (see Sum())
	key []byte
}

func (s *state) Write(p []byte) (n int, err error) {
	// save the key for later use in Sum()
	s.key = append(s.key, p...)

	return s.hasher.Write(p)
}

func (s *state) Sum(b []byte) []byte {
	// append key to final hash digest
	return append(b, append(s.hasher.Sum(nil), s.key...)...)
}

func (s *state) Reset() {
	s.hasher.Reset()
}

func (s *state) Size() int {
	return s.hasher.Size()
}

func (s *state) BlockSize() int {
	return s.hasher.BlockSize()
}

func New128Concat(key []byte) (hash.Hash, error) {
	inner, err := blake2b.New(16, key)
	if err != nil {
		return nil, err
	}

	hasher := state{
		hasher: inner,
		key:    key,
	}

	return &hasher, nil
}

func New128(key []byte) (hash.Hash, error) {
	return blake2b.New(16, key)
}

func New256(key []byte) (hash.Hash, error) {
	return blake2b.New256(key)
}
