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

package signature_test

import (
	"crypto/rand"
	"testing"

	. "github.com/centrifuge/go-substrate-rpc-client/signature"
	"github.com/stretchr/testify/assert"
)

func TestSignAndVerify(t *testing.T) {
	data := []byte("hello!")

	sig, err := Sign(data, TestKeyringPairAlice.URI)
	assert.NoError(t, err)

	ok, err := Verify(data, sig, TestKeyringPairAlice.URI)
	assert.NoError(t, err)

	assert.True(t, ok)
}

func TestSignAndVerifyLong(t *testing.T) {
	data := make([]byte, 258)
	_, err := rand.Read(data)
	assert.NoError(t, err)

	sig, err := Sign(data, TestKeyringPairAlice.URI)
	assert.NoError(t, err)

	ok, err := Verify(data, sig, TestKeyringPairAlice.URI)
	assert.NoError(t, err)

	assert.True(t, ok)
}
