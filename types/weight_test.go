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

func TestWeight_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewWeight(0))
	assertRoundtrip(t, NewWeight(12))
}

func TestWeight_EncodedLength(t *testing.T) {
	assertEncodedLength(t, []encodedLengthAssert{{NewWeight(13), 8}})
}

func TestWeight_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{NewWeight(29), MustHexDecodeString("0x1d00000000000000")},
	})
}

func TestWeight_Hash(t *testing.T) {
	assertHash(t, []hashAssert{
		{NewWeight(29), MustHexDecodeString("0x83e168a13a013e6d47b0778f046aaa05d6c01d6857d044d9e9b658a6d85eb865")},
	})
}

func TestWeight_Hex(t *testing.T) {
	assertEncodeToHex(t, []encodeToHexAssert{
		{NewWeight(29), "0x1d00000000000000"},
	})
}

func TestWeight_String(t *testing.T) {
	assertString(t, []stringAssert{
		{NewWeight(29), "29"},
	})
}

func TestWeight_Eq(t *testing.T) {
	assertEq(t, []eqAssert{
		{NewWeight(23), NewWeight(23), true},
		{NewWeight(23), NewBool(false), false},
	})
}
