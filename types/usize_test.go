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

func TestUSize_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, USize(0))
	assertRoundtrip(t, USize(12))
}

func TestUSize_EncodedLength(t *testing.T) {
	assertEncodedLength(t, []encodedLengthAssert{{USize(13), 4}})
}

func TestUSize_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{USize(29), MustHexDecodeString("0x1d000000")},
	})
}

func TestUSize_Hash(t *testing.T) {
	assertHash(t, []hashAssert{
		{USize(29), MustHexDecodeString("0x60ebb66f09bc7fdd21772ab1ed0efb1fd1208e3f5cd20d2d9a29a2a79b6f953f")},
	})
}

func TestUSize_Hex(t *testing.T) {
	assertEncodeToHex(t, []encodeToHexAssert{
		{USize(29), "0x1d000000"},
	})
}

func TestUSize_String(t *testing.T) {
	assertString(t, []stringAssert{
		{USize(29), "29"},
	})
}

func TestUSize_Eq(t *testing.T) {
	assertEq(t, []eqAssert{
		{USize(23), USize(23), true},
		{USize(23), NewBool(false), false},
	})
}
