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

func TestHealth_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, Health{3, false, true})
	assertRoundtrip(t, Health{1, true, true})
	assertRoundtrip(t, Health{0, true, false})
}

func TestHealth_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{Health{3, false, true}, MustHexDecodeString("0x03000000000000000001")},
		{Health{1, true, true}, MustHexDecodeString("0x01000000000000000101")},
		{Health{0, true, false}, MustHexDecodeString("0x00000000000000000100")},
	})
}

func TestHealth_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{MustHexDecodeString("0x03000000000000000001"), Health{3, false, true}},
		{MustHexDecodeString("0x01000000000000000101"), Health{1, true, true}},
		{MustHexDecodeString("0x00000000000000000100"), Health{0, true, false}},
	})
}
