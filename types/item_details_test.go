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

	fuzz "github.com/google/gofuzz"

	. "github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

var (
	testInstanceDetails = ItemDetails{
		Owner:    newTestAccountID(),
		Approved: NewOptionAccountID(newTestAccountID()),
		IsFrozen: true,
		Deposit:  NewU128(*big.NewInt(123)),
	}

	instanceDetailsFuzzOpts = combineFuzzOpts(
		optionAccountIDFuzzOpts,
		[]fuzzOpt{
			withFuzzFuncs(func(i *ItemDetails, c fuzz.Continue) {
				c.Fuzz(&i.Owner)
				c.Fuzz(&i.Approved)
				c.Fuzz(&i.IsFrozen)
				c.Fuzz(&i.Deposit)
			}),
		},
	)
)

func TestInstanceDetails_EncodeDecode(t *testing.T) {
	assertRoundTripFuzz[ItemDetails](t, 1000, instanceDetailsFuzzOpts...)
	assertDecodeNilData[ItemDetails](t)
	assertEncodeEmptyObj[ItemDetails](t, 50)
}

func TestInstanceDetails_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{
			testInstanceDetails,
			MustHexDecodeString("0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20010102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20017b000000000000000000000000000000"),
		},
	})
}

func TestInstanceDetails_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{
			MustHexDecodeString("0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20010102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20017b000000000000000000000000000000"),
			testInstanceDetails,
		},
	})
}
