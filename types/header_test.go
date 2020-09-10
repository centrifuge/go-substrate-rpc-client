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

var exampleHeader = Header{
	ParentHash:     Hash{1, 2, 3, 4, 5},
	Number:         42,
	StateRoot:      Hash{2, 3, 4, 5, 6},
	ExtrinsicsRoot: Hash{3, 4, 5, 6, 7},
	Digest: Digest{
		{IsOther: true, AsOther: Bytes{4, 5}},
		{IsAuthoritiesChange: true, AsAuthoritiesChange: []AuthorityID{{5, 6}}},
		{IsChangesTrieRoot: true, AsChangesTrieRoot: Hash{6, 7}},
		{IsSealV0: true, AsSealV0: SealV0{7, Signature{8, 9, 10}}},
		{IsConsensus: true, AsConsensus: Consensus{9, Bytes{10, 11, 12}}},
		{IsSeal: true, AsSeal: Seal{11, Bytes{12, 13, 14}}},
		{IsPreRuntime: true, AsPreRuntime: PreRuntime{13, Bytes{14, 15, 16}}},
	},
}

func TestHeader_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, exampleHeader)
}

func TestHeader_EncodedLength(t *testing.T) {
	assertEncodedLength(t, []encodedLengthAssert{{exampleHeader, 269}})
}

func TestHeader_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{exampleHeader, MustHexDecodeString("0x0102030405000000000000000000000000000000000000000000000000000000a8020304050600000000000000000000000000000000000000000000000000000003040506070000000000000000000000000000000000000000000000000000001c000804050104050600000000000000000000000000000000000000000000000000000000000002060700000000000000000000000000000000000000000000000000000000000003070000000000000008090a0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004090000000c0a0b0c050b0000000c0c0d0e060d0000000c0e0f10")}, //nolint:lll
	})
}

func TestHeader_Hash(t *testing.T) {
	assertHash(t, []hashAssert{
		{exampleHeader, MustHexDecodeString("0x4172ec8c10ce00ea6152bb9fb4c63b98305953ebdfa229e47d6b7425842528e4")},
	})
}

func TestHeader_Hex(t *testing.T) {
	assertEncodeToHex(t, []encodeToHexAssert{
		{exampleHeader, "0x0102030405000000000000000000000000000000000000000000000000000000a8020304050600000000000000000000000000000000000000000000000000000003040506070000000000000000000000000000000000000000000000000000001c000804050104050600000000000000000000000000000000000000000000000000000000000002060700000000000000000000000000000000000000000000000000000000000003070000000000000008090a0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004090000000c0a0b0c050b0000000c0c0d0e060d0000000c0e0f10"}, //nolint:lll
	})
}

func TestHeader_String(t *testing.T) {
	assertString(t, []stringAssert{
		{exampleHeader, "{[1 2 3 4 5 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] 42 [2 3 4 5 6 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] [3 4 5 6 7 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] [{true [4 5] false [] false [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] false {0 [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]} false {0 []} false {0 []} false {0 []}} {false [] true [[5 6 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]] false [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] false {0 [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]} false {0 []} false {0 []} false {0 []}} {false [] false [] true [6 7 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] false {0 [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]} false {0 []} false {0 []} false {0 []}} {false [] false [] false [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] true {7 [8 9 10 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]} false {0 []} false {0 []} false {0 []}} {false [] false [] false [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] false {0 [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]} true {9 [10 11 12]} false {0 []} false {0 []}} {false [] false [] false [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] false {0 [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]} false {0 []} true {11 [12 13 14]} false {0 []}} {false [] false [] false [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] false {0 [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]} false {0 []} false {0 []} true {13 [14 15 16]}}]}"}, //nolint:lll
	})
}

func TestHeader_Eq(t *testing.T) {
	assertEq(t, []eqAssert{
		{exampleHeader, exampleHeader, true},
		{exampleHeader, NewBytes(hash64), false},
		{exampleHeader, NewBool(false), false},
	})
}
