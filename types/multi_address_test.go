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

package types_test

import (
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

func TestNewMultiAddressFromAccountID(t *testing.T) {
	assertRoundtrip(t, types.NewMultiAddressFromAccountID(signature.TestKeyringPairAlice.PublicKey))

	_, err := types.NewMultiAddressFromHexAccountID("123!")
	assert.Error(t, err)

	addr, err := types.NewMultiAddressFromHexAccountID(types.HexEncodeToString(signature.TestKeyringPairAlice.PublicKey))
	assert.NoError(t, err)
	assertRoundtrip(t, addr)
	assertRoundtrip(t, types.MultiAddress{
		IsIndex: true,
		AsIndex: 100,
	})
	assertRoundtrip(t, types.MultiAddress{
		IsRaw: true,
		AsRaw: []byte{1, 2, 3},
	})
	assertRoundtrip(t, types.MultiAddress{
		IsAddress32: true,
		AsAddress32: [32]byte{},
	})
	assertRoundtrip(t, types.MultiAddress{
		IsAddress20: true,
		AsAddress20: [20]byte{},
	})
}
