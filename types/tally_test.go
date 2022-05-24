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
	"math/big"
	"testing"

	. "github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

var (
	testTally = Tally{
		Votes: NewU128(*big.NewInt(123)),
		Total: NewU128(*big.NewInt(456)),
	}
)

func TestTally_EncodeDecode(t *testing.T) {
	assertRoundTripFuzz[Tally](t, 100)
	assertDecodeNilData[Tally](t)
	assertEncodeEmptyObj[Tally](t, 32)
}

func TestTally_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{testTally, MustHexDecodeString("0x7b000000000000000000000000000000c8010000000000000000000000000000")},
	})
}

func TestTally_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{MustHexDecodeString("0x7b000000000000000000000000000000c8010000000000000000000000000000"), testTally},
	})
}
