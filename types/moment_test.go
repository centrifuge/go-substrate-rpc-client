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
	"time"

	. "github.com/centrifuge/go-substrate-rpc-client/types"
)

func TestMoment_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewMoment(time.Unix(1575470205, 874000000000)))
	assertRoundtrip(t, NewMoment(time.Unix(0, 0)))
}

func TestMoment_EncodedLength(t *testing.T) {
	assertEncodedLength(t, []encodedLengthAssert{{NewMoment(time.Unix(12345, 0)), 8}})
}

func TestMoment_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{NewMoment(time.Unix(12345, 0)), MustHexDecodeString("0xa85ebc0000000000")},
	})
}

func TestMoment_Hash(t *testing.T) {
	assertHash(t, []hashAssert{
		{NewMoment(time.Unix(12345, 0)), MustHexDecodeString(
			"0x388c23ace057ec800ef437dcc68bfbcd2b1a22fe79aceb623d7d21f7bda56848")},
	})
}

func TestMoment_Hex(t *testing.T) {
	assertEncodeToHex(t, []encodeToHexAssert{
		{NewMoment(time.Unix(12345, 0)), "0xa85ebc0000000000"},
	})
}

func TestMoment_String(t *testing.T) {
	assertString(t, []stringAssert{
		{NewMoment(time.Unix(12345, 0).UTC()), "1970-01-01 03:25:45 +0000 UTC"},
	})
}

func TestMoment_Eq(t *testing.T) {
	assertEq(t, []eqAssert{
		{NewMoment(time.Unix(12345, 0)), NewMoment(time.Unix(12345, 0)), true},
		{NewMoment(time.Unix(12345, 0)), NewBool(false), false},
	})
}
