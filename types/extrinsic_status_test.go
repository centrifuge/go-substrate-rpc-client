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
	"encoding/json"
	"testing"

	. "github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/stretchr/testify/assert"
)

var testExtrinsicStatus0 = ExtrinsicStatus{IsFuture: true}
var testExtrinsicStatus1 = ExtrinsicStatus{IsReady: true}
var testExtrinsicStatus2 = ExtrinsicStatus{IsBroadcast: true, AsBroadcast: []Text{"This", "is", "broadcast"}}
var testExtrinsicStatus3 = ExtrinsicStatus{IsInBlock: true, AsInBlock: NewHash([]byte{0xaa})}
var testExtrinsicStatus4 = ExtrinsicStatus{IsRetracted: true, AsRetracted: NewHash([]byte{0xbb})}
var testExtrinsicStatus5 = ExtrinsicStatus{IsFinalityTimeout: true, AsFinalityTimeout: NewHash([]byte{0xcc})}
var testExtrinsicStatus6 = ExtrinsicStatus{IsFinalized: true, AsFinalized: NewHash([]byte{0xdd})}
var testExtrinsicStatus7 = ExtrinsicStatus{IsUsurped: true, AsUsurped: NewHash([]byte{0xee})}
var testExtrinsicStatus8 = ExtrinsicStatus{IsDropped: true}
var testExtrinsicStatus9 = ExtrinsicStatus{IsInvalid: true}

func TestExtrinsicStatus_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, testExtrinsicStatus0)
	assertRoundtrip(t, testExtrinsicStatus1)
	assertRoundtrip(t, testExtrinsicStatus2)
	assertRoundtrip(t, testExtrinsicStatus3)
	assertRoundtrip(t, testExtrinsicStatus4)
	assertRoundtrip(t, testExtrinsicStatus5)
	assertRoundtrip(t, testExtrinsicStatus6)
	assertRoundtrip(t, testExtrinsicStatus7)
	assertRoundtrip(t, testExtrinsicStatus8)
	assertRoundtrip(t, testExtrinsicStatus9)
}

func TestExtrinsicStatus_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{testExtrinsicStatus0, []byte{0x00}},
		{testExtrinsicStatus1, []byte{0x01}},
		{testExtrinsicStatus2, MustHexDecodeString("0x020c10546869730869732462726f616463617374")},
		{testExtrinsicStatus3, MustHexDecodeString("0x03aa00000000000000000000000000000000000000000000000000000000000000")}, //nolint:lll
		{testExtrinsicStatus4, MustHexDecodeString("0x04bb00000000000000000000000000000000000000000000000000000000000000")}, //nolint:lll
		{testExtrinsicStatus5, MustHexDecodeString("0x05cc00000000000000000000000000000000000000000000000000000000000000")}, //nolint:lll
		{testExtrinsicStatus6, MustHexDecodeString("0x06dd00000000000000000000000000000000000000000000000000000000000000")}, //nolint:lll
		{testExtrinsicStatus7, MustHexDecodeString("0x07ee00000000000000000000000000000000000000000000000000000000000000")}, //nolint:lll
		{testExtrinsicStatus8, []byte{0x08}},
		{testExtrinsicStatus9, []byte{0x09}},
	})
}

func TestExtrinsicStatus_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{[]byte{0x00}, testExtrinsicStatus0},
		{[]byte{0x01}, testExtrinsicStatus1},
		{MustHexDecodeString("0x020c10546869730869732462726f616463617374"), testExtrinsicStatus2},
		{MustHexDecodeString("0x03aa00000000000000000000000000000000000000000000000000000000000000"), testExtrinsicStatus3}, //nolint:lll
		{MustHexDecodeString("0x04bb00000000000000000000000000000000000000000000000000000000000000"), testExtrinsicStatus4}, //nolint:lll
		{MustHexDecodeString("0x05cc00000000000000000000000000000000000000000000000000000000000000"), testExtrinsicStatus5}, //nolint:lll
		{MustHexDecodeString("0x06dd00000000000000000000000000000000000000000000000000000000000000"), testExtrinsicStatus6}, //nolint:lll
		{MustHexDecodeString("0x07ee00000000000000000000000000000000000000000000000000000000000000"), testExtrinsicStatus7}, //nolint:lll
		{[]byte{0x08}, testExtrinsicStatus8},
		{[]byte{0x09}, testExtrinsicStatus9},
	})
}

var testExtrinsicStatusTestCases = []struct {
	encoded []byte
	decoded ExtrinsicStatus
}{
	{
		[]byte("\"future\""),
		ExtrinsicStatus{IsFuture: true},
	}, {
		[]byte("\"ready\""),
		ExtrinsicStatus{IsReady: true},
	}, {
		[]byte("{\"broadcast\":[\"hello\",\"world\"]}"),
		ExtrinsicStatus{IsBroadcast: true, AsBroadcast: []Text{"hello", "world"}},
	}, {
		[]byte("{\"inBlock\":\"0x95e3b7f86541d06306691a2fe8cbd935d0bdd28ea14fe515e2db0fa87847f8f8\"}"),
		ExtrinsicStatus{IsInBlock: true, AsInBlock: NewHash(MustHexDecodeString(
			"0x95e3b7f86541d06306691a2fe8cbd935d0bdd28ea14fe515e2db0fa87847f8f8"))},
	}, {
		[]byte("{\"retracted\":\"0x95e3b7f86541d06306691a2fe8cbd935d0bdd28ea14fe515e2db0fa87847f8f8\"}"),
		ExtrinsicStatus{IsRetracted: true, AsRetracted: NewHash(MustHexDecodeString(
			"0x95e3b7f86541d06306691a2fe8cbd935d0bdd28ea14fe515e2db0fa87847f8f8"))},
	}, {
		[]byte("{\"finalityTimeout\":\"0x95e3b7f86541d06306691a2fe8cbd935d0bdd28ea14fe515e2db0fa87847f8f8\"}"),
		ExtrinsicStatus{IsFinalityTimeout: true, AsFinalityTimeout: NewHash(MustHexDecodeString(
			"0x95e3b7f86541d06306691a2fe8cbd935d0bdd28ea14fe515e2db0fa87847f8f8"))},
	}, {
		[]byte("{\"finalized\":\"0x95e3b7f86541d06306691a2fe8cbd935d0bdd28ea14fe515e2db0fa87847f8f8\"}"),
		ExtrinsicStatus{IsFinalized: true, AsFinalized: NewHash(MustHexDecodeString(
			"0x95e3b7f86541d06306691a2fe8cbd935d0bdd28ea14fe515e2db0fa87847f8f8"))},
	}, {
		[]byte("{\"usurped\":\"0x95e3b7f86541d06306691a2fe8cbd935d0bdd28ea14fe515e2db0fa87847f8ab\"}"),
		ExtrinsicStatus{IsUsurped: true, AsUsurped: NewHash(MustHexDecodeString(
			"0x95e3b7f86541d06306691a2fe8cbd935d0bdd28ea14fe515e2db0fa87847f8ab"))},
	}, {
		[]byte("\"dropped\""),
		ExtrinsicStatus{IsDropped: true},
	}, {
		[]byte("\"invalid\""),
		ExtrinsicStatus{IsInvalid: true},
	},
}

func TestExtrinsicStatus_UnmarshalJSON(t *testing.T) {
	for _, test := range testExtrinsicStatusTestCases {
		var actual ExtrinsicStatus
		err := json.Unmarshal(test.encoded, &actual)
		assert.NoError(t, err)
		assert.Equal(t, test.decoded, actual)
	}
}

func TestExtrinsicStatus_MarshalJSON(t *testing.T) {
	for _, test := range testExtrinsicStatusTestCases {
		actual, err := json.Marshal(test.decoded)
		assert.NoError(t, err)
		assert.Equal(t, test.encoded, actual)
	}
}
