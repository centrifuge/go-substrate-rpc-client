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

	. "github.com/centrifuge/go-substrate-rpc-client/v3/types"
)

var testDigestItem1 = DigestItem{IsOther: true, AsOther: NewBytes([]byte{0xab})}
var testDigestItem2 = DigestItem{IsChangesTrieRoot: true, AsChangesTrieRoot: NewHash([]byte{0x01, 0x02, 0x03})}

func TestDigestItem_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, testDigestItem1)
	assertRoundtrip(t, testDigestItem2)
	assertRoundtrip(t, DigestItem{
		IsPreRuntime: true,
		AsPreRuntime: PreRuntime{},
	})
	assertRoundtrip(t, DigestItem{
		IsConsensus: true,
		AsConsensus: Consensus{},
	})
	assertRoundtrip(t, DigestItem{
		IsSeal: true,
		AsSeal: Seal{},
	})
	assertRoundtrip(t, DigestItem{
		IsChangesTrieSignal: true,
		AsChangesTrieSignal: ChangesTrieSignal{IsNewConfiguration: true, AsNewConfiguration: NewBytes([]byte{1, 2, 3})},
	})
}
