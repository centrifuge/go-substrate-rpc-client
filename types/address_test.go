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
	"encoding/binary"
	"testing"

	. "github.com/centrifuge/go-substrate-rpc-client/types"
)

func TestAddress_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, NewAddressFromAccountID([]byte{128, 42}))
	assertRoundtrip(t, NewAddressFromAccountIndex(421))
}

func TestAddress_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{NewAddressFromAccountID([]byte{
			1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8,
			1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8,
		}), []byte{
			255,
			1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8,
			1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8,
		}},
		{NewAddressFromAccountIndex(binary.BigEndian.Uint32([]byte{
			17, 18, 19, 20,
		})), []byte{
			253, 20, 19, 18, 17, // order is reversed because scale uses little endian
		}},
		{NewAddressFromAccountIndex(uint32(binary.BigEndian.Uint16([]byte{
			21, 22,
		}))), []byte{
			252, 22, 21, // order is reversed because scale uses little endian
		}},
		{NewAddressFromAccountIndex(uint32(23)), []byte{23}},
	})
}

func TestAddress_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{[]byte{
			255,
			1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8,
			1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8,
		}, NewAddressFromAccountID([]byte{
			1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8,
			1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8,
		})},
		{[]byte{
			253, 20, 19, 18, 17, // order is reversed because scale uses little endian
		}, NewAddressFromAccountIndex(binary.BigEndian.Uint32([]byte{
			17, 18, 19, 20,
		}))},
		{[]byte{
			252, 22, 21, // order is reversed because scale uses little endian
		}, NewAddressFromAccountIndex(uint32(binary.BigEndian.Uint16([]byte{
			21, 22,
		})))},
		{[]byte{23}, NewAddressFromAccountIndex(uint32(23))},
	})
}
