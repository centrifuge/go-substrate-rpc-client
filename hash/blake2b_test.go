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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlake2_128Concat(t *testing.T) {
	h, err := NewBlake2b128Concat(nil)
	assert.NoError(t, err)
	n, err := h.Write([]byte("abc"))
	assert.NoError(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, []byte{
		0xcf, 0x4a, 0xb7, 0x91, 0xc6, 0x2b, 0x8d, 0x2b, 0x21, 0x9, 0xc9, 0x2, 0x75, 0x28, 0x78, 0x16, 0x61, 0x62, 0x63,
	}, h.Sum(nil))
	assert.Equal(t, 128, h.BlockSize())
	assert.Equal(t, 19, h.Size())
	h.Reset()
	assert.Equal(t, 16, h.Size())
}
