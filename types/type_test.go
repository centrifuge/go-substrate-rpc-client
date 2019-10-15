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

func TestType_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, Type(""))
	assertRoundtrip(t, Type("My nice type"))
}

func TestType_EncodedLength(t *testing.T) {
	assertEncodedLength(t, []encodedLengthAssert{
		{Type(""), 1},
		{Type("My nice type"), 13},
		{Type("ښ"), 3},
	})
}

func TestType_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{Type(""), MustHexDecodeString("0x00")},
		{Type("My nice type"), MustHexDecodeString("0x304d79206e6963652074797065")},
		{Type("ښ"), MustHexDecodeString("0x08da9a")},
	})
}

func TestType_Hash(t *testing.T) {
	assertHash(t, []hashAssert{
		{Type(""), MustHexDecodeString("0x03170a2e7597b7b7e3d84c05391d139a62b157e78786d8c082f29dcf4c111314")},
		{Type("My nice type"), MustHexDecodeString(
			"0x21b1f717069b923d0d6dbbcef60497b18e45443d9d4b42b06b168c3b5c914646")},
	})
}

func TestType_Hex(t *testing.T) {
	assertEncodeToHex(t, []encodeToHexAssert{
		{Type(""), "0x00"},
		{Type("My nice type"), "0x304d79206e6963652074797065"},
		{Type("ښ"), "0x08da9a"},
	})
}

func TestType_Type(t *testing.T) {
	assertString(t, []stringAssert{
		{Type(""), ""},
		{Type("My nice type"), "My nice type"},
		{Type("ښ"), "ښ"},
	})
}

func TestType_Eq(t *testing.T) {
	assertEq(t, []eqAssert{
		{Type("My nice type"), Type("My nice type"), true},
		{Type(""), Type("23"), false},
		{Type("My nice type"), NewU8(23), false},
		{Type("My nice type"), NewBool(false), false},
	})
}
