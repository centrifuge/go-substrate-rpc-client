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

var testExtrinsicStatus1 = ExtrinsicStatus{IsFuture: true}
var testExtrinsicStatus2 = ExtrinsicStatus{IsReady: true}
var testExtrinsicStatus3 = ExtrinsicStatus{IsFinalized: true, AsFinalized: NewHash([]byte{0xab})}
var testExtrinsicStatus4 = ExtrinsicStatus{IsUsurped: true, AsUsurped: NewHash([]byte{0xcd})}
var testExtrinsicStatus5 = ExtrinsicStatus{IsBroadcast: true, AsBroadcast: []Text{"This", "is", "broadcast"}}
var testExtrinsicStatus6 = ExtrinsicStatus{IsDropped: true}
var testExtrinsicStatus7 = ExtrinsicStatus{IsInvalid: true}

func TestExtrinsicStatus_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, testExtrinsicStatus1)
	assertRoundtrip(t, testExtrinsicStatus2)
	assertRoundtrip(t, testExtrinsicStatus3)
	assertRoundtrip(t, testExtrinsicStatus4)
	assertRoundtrip(t, testExtrinsicStatus5)
	assertRoundtrip(t, testExtrinsicStatus6)
	assertRoundtrip(t, testExtrinsicStatus7)
}

func TestExtrinsicStatus_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{testExtrinsicStatus1, []byte{0x00}},
		{testExtrinsicStatus2, []byte{0x01}},
		{testExtrinsicStatus3, MustHexDecodeString("0x02ab00000000000000000000000000000000000000000000000000000000000000")}, //nolint:lll
		{testExtrinsicStatus4, MustHexDecodeString("0x03cd00000000000000000000000000000000000000000000000000000000000000")}, //nolint:lll
		{testExtrinsicStatus5, MustHexDecodeString("0x040c10546869730869732462726f616463617374")},
		{testExtrinsicStatus6, []byte{0x05}},
		{testExtrinsicStatus7, []byte{0x06}},
	})
}

func TestExtrinsicStatus_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{[]byte{0x00}, testExtrinsicStatus1},
		{[]byte{0x01}, testExtrinsicStatus2},
		{MustHexDecodeString("0x02ab00000000000000000000000000000000000000000000000000000000000000"), testExtrinsicStatus3}, //nolint:lll
		{MustHexDecodeString("0x03cd00000000000000000000000000000000000000000000000000000000000000"), testExtrinsicStatus4}, //nolint:lll
		{MustHexDecodeString("0x040c10546869730869732462726f616463617374"), testExtrinsicStatus5},
		{[]byte{0x05}, testExtrinsicStatus6},
		{[]byte{0x06}, testExtrinsicStatus7},
	})
}
