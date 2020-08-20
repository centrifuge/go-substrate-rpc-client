// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
//
// Copyright 2019 Centrifuge GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hash

import (
	"hash"

	"golang.org/x/crypto/blake2b"
)

type blake2b128Concat struct {
	data []byte
}

func NewBlake2b128Concat(b []byte) hash.Hash {
	return &blake2b128Concat{data: b}
}

// Write (via the embedded io.Writer interface) adds more data to the running hash.
// It never returns an error.
func (bc *blake2b128Concat) Write(p []byte) (n int, err error) {
	bc.data = append(bc.data, p...)
	return len(p), nil
}

// Sum appends the current hash to b and returns the resulting slice.
// It does not change the underlying hash state.
func (bc *blake2b128Concat) Sum(b []byte) []byte {
	res := make([]byte, 0, 8)

	h, err := blake2b.New(16, nil)
	if err != nil {
		panic(err)
	}
	_, err = h.Write(bc.data)
	if err != nil {
		panic(err)
	}
	res = append(res, h.Sum(nil)...)
	res = append(res, bc.data...)

	return append(b, res...)
}

// Reset resets the Hash to its initial state.
func (bc *blake2b128Concat) Reset() {
	bc.data = make([]byte, 0)
}

// Size returns the number of bytes Sum will return.
func (bc *blake2b128Concat) Size() int {
	return len(bc.Sum(nil))
}

// BlockSize returns the hash's underlying block size.
// The Write method must be able to accept any amount
// of data, but it may operate more efficiently if all writes
// are a multiple of the block size.
func (bc *blake2b128Concat) BlockSize() int {
	return blake2b.BlockSize
}
