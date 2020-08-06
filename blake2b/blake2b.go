// Package blake2b is a thin wrapper over golang.org/x/crypto/blake2b,
// and additionally provides a BLAKE2b-128-concat hashing algorithm.
package blake2b

import (
	"hash"

	"golang.org/x/crypto/blake2b"
)

// State for the blake2b_128_concat hasher
type concatState struct {
	// the underlying blake2b_128 hasher
	hasher hash.Hash

	// the key we need to append to the finalized hash digest (see Sum())
	key []byte
}

func (s *concatState) Write(p []byte) (n int, err error) {
	// save the key for later use in Sum()
	s.key = append(s.key, p...)

	return s.hasher.Write(p)
}

func (s *concatState) Sum(b []byte) []byte {
	// append key to final hash digest
	return append(s.hasher.Sum(b), s.key...)
}

func (s *concatState) Reset() {
	s.key = nil
	s.hasher.Reset()
}

func (s *concatState) Size() int {
	return s.hasher.Size() + len(s.key)
}

func (s *concatState) BlockSize() int {
	return s.hasher.BlockSize()
}

// New128Concat returns a new hash.Hash computing the BLAKE2b-128-concat checksum. A non-nil
// key turns the hash into a MAC. The key must be between zero and 64 bytes long.
func New128Concat(key []byte) (hash.Hash, error) {
	inner, err := blake2b.New(16, key)
	if err != nil {
		return nil, err
	}

	hasher := concatState{
		hasher: inner,
		key:    key,
	}

	return &hasher, nil
}

// New128 returns a new hash.Hash computing the BLAKE2b-128 checksum. A non-nil
// key turns the hash into a MAC. The key must be between zero and 64 bytes long.
func New128(key []byte) (hash.Hash, error) {
	return blake2b.New(16, key)
}

// New256 returns a new hash.Hash computing the BLAKE2b-256 checksum. A non-nil
// key turns the hash into a MAC. The key must be between zero and 64 bytes long.
func New256(key []byte) (hash.Hash, error) {
	return blake2b.New256(key)
}
