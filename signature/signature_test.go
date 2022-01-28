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
	"bytes"
	"testing"

	"github.com/ComposableFi/go-substrate-rpc-client/v4/scale"
	"github.com/ComposableFi/go-substrate-rpc-client/v4/types"
)

func TestDecodeSubstrateMessage(t *testing.T) {
	hex := "0xeee4318d5f09f8e6bb227855f3d1f265ce2bd6303316c560b9c1bc31be8149e761010000000000000000000014013845a96faf999274e973735201a2a0e70c7ebec67d016a1e1d5a68cd3205a41a708588694d948865da4f640cb58ca097a2dcb51baa33606453271136ca1ae7d1000001e05c7d9481cb6d7ae52c1a1b7608e1f55cf3e63a177306bca813a38c375da18f48f875874ef4a4931cbc9d7509ef9c3fa8144ea4407a65c2d31376341e981f560001b85e8d525f60a5f26a2350bd0b632d4c9f505bdd0d826f314cb9524a88a8b36f59e678dfc266ee8676db39a993b350ce5bbc33ab514bcac710c0ba4eedc10c7e0001626fe670748b10bbdecdacfc8253b7819afca627ee9a2787149d664aff76daa94a3bc509172cf444eca4a93025c8f90e5a56e4eff0b16f11511bebbbdbaecce300"
	signedCommitment := &types.SignedCommitment{}

	bz, err := types.HexDecodeString(hex)
	if err != nil {
		t.Error(err)
	}
	errr := scale.NewDecoder(bytes.NewReader(bz)).Decode(signedCommitment)

	if errr != nil {
		t.Errorf("error decoding %+v ", errr.Error())
	}

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
