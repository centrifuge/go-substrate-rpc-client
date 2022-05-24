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
	"testing"

	"github.com/ComposableFi/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/require"
)

func TestDecodeSubstrateMessage(t *testing.T) {
	hex := "0x046d68809de63e597f2ff84b4e287c829b93a3aad304b1aebd51ec2d7392d41047e951c723040000690000000000000004d80500000010d260924684152d72c9d81d028382b4847f895ac6081cdcaa93c37e00df3b275b7fbc96c10c6fc86743d3aaba0bf59278822b1610218858edfbaf77a9884ee031004f733874acfbdd9d3e45e9b1247647d02761d1673c303a7982cd6eb5d573eb3036415493a0d655937c64a9536b3113747ec972022b7250c89c30888a45a3908c013251044285c3000d54d0fc0572a0559473e7e223175daa44046b3523d4c4ef7a71f6a4287cc7030d8fd30f504ea0360e0a4c3047e33f8f5d4e4c11f867cbac700047e643aee47dacc161c8cb049597f30a0c74ed596cd78b94f21526c1d744835758ad4492002c5890a0f45990aa175a320a676ce89977cffbe2eb73b2cf16d62900"
	compactCommitment := types.CompactSignedCommitment{}

	// attempt to decode the SignedCommitments
	err := types.DecodeFromHexString(hex, &compactCommitment)
	require.NoError(t, err)

	signedCommitment := compactCommitment.Unpack()

	t.Logf("COMMITMENT %#v \n", signedCommitment.Commitment)

	for _, v := range signedCommitment.Signatures {
		if v.HasValue {
			_, sig := v.Unwrap()
			t.Logf("SIGNATURE %#v \n", types.HexEncodeToString(sig[:]))

		} else {
			t.Logf("NONE\n")

		}
	}
}
