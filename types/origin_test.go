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

// newOrigin creates a new Origin type. This function is not exported by purpose – Origin should be ignored and not be
// allowed to be constructed.
func newOrigin() Origin {
	return Origin(0x00)
}

func TestOrigin_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, newOrigin())
}

func TestOrigin_EncodedLength(t *testing.T) {
	assertEncodedLength(t, []encodedLengthAssert{
		{newOrigin(), 0},
	})
}

func TestOrigin_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{newOrigin(), MustHexDecodeString("0x")},
	})
}

func TestOrigin_Hash(t *testing.T) {
	assertHash(t, []hashAssert{
		{newOrigin(), MustHexDecodeString(
			"0x0e5751c026e543b2e8ab2eb06099daa1d1e5df47778f7787faab45cdf12fe3a8")},
	})
}

func TestOrigin_Hex(t *testing.T) {
	assertEncodeToHex(t, []encodeToHexAssert{
		{newOrigin(), ""},
	})
}

func TestOrigin_String(t *testing.T) {
	assertString(t, []stringAssert{
		{newOrigin(), ""},
	})
}

func TestOrigin_Eq(t *testing.T) {
	assertEq(t, []eqAssert{
		{newOrigin(), newOrigin(), true},
		{newOrigin(), NewBytes([]byte{}), false},
		{newOrigin(), NewBool(true), false},
		{newOrigin(), NewBool(false), false},
	})
}
