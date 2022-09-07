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
	. "github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	. "github.com/centrifuge/go-substrate-rpc-client/v4/types/test_utils"
	fuzz "github.com/google/gofuzz"
)

var (
	testTranche = Tranche{
		FirstVal:  323,
		SecondVal: [16]U8{4, 5, 6, 3, 1, 3, 2, 4},
	}
)

func TestTranche_EncodeDecode(t *testing.T) {
	AssertRoundTripFuzz[Tranche](t, 100)
	AssertDecodeNilData[Tranche](t)
	AssertEncodeEmptyObj[Tranche](t, 24)
}

func TestTranche_Encode(t *testing.T) {
	AssertEncode(t, []EncodingAssert{
		{testTranche, MustHexDecodeString("0x430100000000000004050603010302040000000000000000")},
	})
}

func TestTranche_Decode(t *testing.T) {
	AssertDecode(t, []DecodingAssert{
		{MustHexDecodeString("0x430100000000000004050603010302040000000000000000"), testTranche},
	})
}

var (
	testCurrencyID1 = CurrencyID{
		IsNative: true,
	}
	testCurrencyID2 = CurrencyID{
		IsUsd: true,
	}
	testCurrencyID3 = CurrencyID{
		IsTranche: true,
		Tranche:   testTranche,
	}
	testCurrencyID4 = CurrencyID{
		IsKSM: true,
	}
	testCurrencyID5 = CurrencyID{
		IsKUSD: true,
	}
	testCurrencyID6 = CurrencyID{
		IsPermissioned:       true,
		PermissionedCurrency: PermissionedCurrency{},
	}

	currencyIDFuzzOpts = []FuzzOpt{
		WithFuzzFuncs(func(cID *CurrencyID, c fuzz.Continue) {
			switch c.Intn(6) {
			case 0:
				cID.IsNative = true
			case 1:
				cID.IsUsd = true
			case 2:
				cID.IsTranche = true

				c.Fuzz(&cID.Tranche)
			case 3:
				cID.IsKSM = true
			case 4:
				cID.IsKUSD = true
			case 5:
				cID.IsPermissioned = true

				c.Fuzz(&cID.PermissionedCurrency)
			}
		}),
	}
)

func TestCurrencyID_EncodeDecode(t *testing.T) {
	AssertRoundTripFuzz[CurrencyID](t, 1000, currencyIDFuzzOpts...)
	AssertDecodeNilData[CurrencyID](t)
	AssertEncodeEmptyObj[CurrencyID](t, 0)
}

func TestCurrencyID_Encode(t *testing.T) {
	AssertEncode(t, []EncodingAssert{
		{testCurrencyID1, MustHexDecodeString("0x00")},
		{testCurrencyID2, MustHexDecodeString("0x01")},
		{testCurrencyID3, MustHexDecodeString("0x02430100000000000004050603010302040000000000000000")},
		{testCurrencyID4, MustHexDecodeString("0x03")},
		{testCurrencyID5, MustHexDecodeString("0x04")},
		{testCurrencyID6, MustHexDecodeString("0x05")},
	})
}

func TestCurrencyID_Decode(t *testing.T) {
	AssertDecode(t, []DecodingAssert{
		{MustHexDecodeString("0x00"), testCurrencyID1},
		{MustHexDecodeString("0x01"), testCurrencyID2},
		{MustHexDecodeString("0x02430100000000000004050603010302040000000000000000"), testCurrencyID3},
		{MustHexDecodeString("0x03"), testCurrencyID4},
		{MustHexDecodeString("0x04"), testCurrencyID5},
		{MustHexDecodeString("0x05"), testCurrencyID6},
	})
}

var (
	testPrice = Price{
		CurrencyID: testCurrencyID4,
		Amount:     NewU128(*big.NewInt(123)),
	}
)

func TestPrice_EncodeDecode(t *testing.T) {
	AssertRoundTripFuzz[Price](t, 100, currencyIDFuzzOpts...)
	AssertDecodeNilData[Price](t)
	AssertEncodeEmptyObj[Price](t, 16)
}

func TestPrice_Encode(t *testing.T) {
	AssertEncode(t, []EncodingAssert{
		{testPrice, MustHexDecodeString("0x037b000000000000000000000000000000")},
	})
}

func TestPrice_Decode(t *testing.T) {
	AssertDecode(t, []DecodingAssert{
		{MustHexDecodeString("0x037b000000000000000000000000000000"), testPrice},
	})
}

var (
	testSale = Sale{
		Seller: newTestAccountID(),
		Price:  testPrice,
	}
)

func TestSale_EncodeDecode(t *testing.T) {
	AssertRoundTripFuzz[Sale](t, 100, currencyIDFuzzOpts...)
	AssertDecodeNilData[Sale](t)
	AssertEncodeEmptyObj[Sale](t, 48)
}

func TestSale_Encode(t *testing.T) {
	AssertEncode(t, []EncodingAssert{
		{
			testSale,
			MustHexDecodeString("0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20037b000000000000000000000000000000"),
		},
	})
}

func TestSale_Decode(t *testing.T) {
	AssertDecode(t, []DecodingAssert{
		{
			MustHexDecodeString("0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20037b000000000000000000000000000000"),
			testSale,
		},
	})
}
