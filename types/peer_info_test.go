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

var testPeerInfo = PeerInfo{
	PeerID:          "abc12345",
	Roles:           "some roles",
	ProtocolVersion: 123,
	BestHash:        NewHash([]byte{0xab, 0xcd}),
	BestNumber:      1337,
}

func TestPeerInfo_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, testPeerInfo)
}

func TestPeerInfo_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{testPeerInfo, MustHexDecodeString("0x20616263313233343528736f6d6520726f6c65737b000000abcd000000000000000000000000000000000000000000000000000000000000e514")}, //nolint:lll
	})
}

func TestPeerInfo_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{MustHexDecodeString("0x20616263313233343528736f6d6520726f6c65737b000000abcd000000000000000000000000000000000000000000000000000000000000e514"), testPeerInfo}, //nolint:lll
	})
}
