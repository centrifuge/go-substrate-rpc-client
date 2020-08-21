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

func TestAccountInfoV4_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewAccountInfoV4([]byte{1, 2, 3}, 13))
}

func TestAccountInfoV4_EncodedLength(t *testing.T) {
	assertEncodedLength(t, []encodedLengthAssert{
		{NewAccountInfoV4([]byte{1, 2, 3}, 13), 12},
	})
}

func TestAccountInfoV4_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{NewAccountInfoV4([]byte{1, 2, 3}, 13), MustHexDecodeString("0x0c0102030d00000000000000")},
	})
}

func TestAccountInfoV4_Hash(t *testing.T) {
	assertHash(t, []hashAssert{
		{NewAccountInfoV4([]byte{1, 2, 3}, 13), MustHexDecodeString(
			"0x4fac0dfeb9b4efd2518c762e7d097fafaffaf8d56a2e784f9fc9919c22277804")},
	})
}

func TestAccountInfoV4_Hex(t *testing.T) {
	assertEncodeToHex(t, []encodeToHexAssert{
		{NewAccountInfoV4([]byte{1, 2, 3}, 13), "0x0c0102030d00000000000000"},
	})
}

func TestAccountInfoV4_String(t *testing.T) {
	assertString(t, []stringAssert{
		{NewAccountInfoV4([]byte{1, 2, 3}, 13), "{[1 2 3] 13}"},
	})
}

func TestAccountInfoV4_Eq(t *testing.T) {
	assertEq(t, []eqAssert{
		{NewAccountInfoV4([]byte{1, 2, 3}, 13), NewAccountInfoV4([]byte{1, 2, 3}, 13), true},
		{NewAccountInfoV4([]byte{1, 2, 3}, 13), NewAccountInfoV4([]byte{1, 2, 2}, 13), false},
		{NewAccountInfoV4([]byte{1, 2, 3}, 13), NewBool(false), false},
	})
}
