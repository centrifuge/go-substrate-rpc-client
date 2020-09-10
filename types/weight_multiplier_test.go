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

	. "github.com/centrifuge/go-substrate-rpc-client/types"
)

func TestWeightMultiplier_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewWeightMultiplier(0))
	assertRoundtrip(t, NewWeightMultiplier(12))
	assertRoundtrip(t, NewWeightMultiplier(-12))
}

func TestWeightMultiplier_EncodedLength(t *testing.T) {
	assertEncodedLength(t, []encodedLengthAssert{{NewWeightMultiplier(-13), 8}})
}

func TestWeightMultiplier_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{NewWeightMultiplier(-29), MustHexDecodeString("0xe3ffffffffffffff")},
	})
}

func TestWeightMultiplier_Hash(t *testing.T) {
	assertHash(t, []hashAssert{
		{NewWeightMultiplier(-29), MustHexDecodeString("0x4d42db2aa4a23bde81a3ad3705220affaa457c56a0135080c71db7783fec8f44")},
	})
}

func TestWeightMultiplier_Hex(t *testing.T) {
	assertEncodeToHex(t, []encodeToHexAssert{
		{NewWeightMultiplier(-29), "0xe3ffffffffffffff"},
	})
}

func TestWeightMultiplier_String(t *testing.T) {
	assertString(t, []stringAssert{
		{NewWeightMultiplier(-29), "-29"},
	})
}

func TestWeightMultiplier_Eq(t *testing.T) {
	assertEq(t, []eqAssert{
		{NewWeightMultiplier(23), NewWeightMultiplier(23), true},
		{NewWeightMultiplier(-23), NewWeightMultiplier(23), false},
		{NewWeightMultiplier(23), NewU64(23), false},
		{NewWeightMultiplier(23), NewBool(false), false},
	})
}
