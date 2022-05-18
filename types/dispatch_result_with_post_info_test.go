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

	fuzz "github.com/google/gofuzz"

	. "github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

var (
	testDispatchResultWithPostInfo1 = DispatchResultWithPostInfo{
		IsOk: true,
		Ok: PostDispatchInfo{
			ActualWeight: NewOptionWeight(123),
			PaysFee: Pays{
				IsYes: true,
			},
		},
	}
	testDispatchResultWithPostInfo2 = DispatchResultWithPostInfo{
		IsError: true,
		Error: DispatchErrorWithPostInfo{
			PostInfo: PostDispatchInfo{
				ActualWeight: NewOptionWeight(456),
				PaysFee: Pays{
					IsNo: true,
				},
			},
			Error: DispatchError{
				IsOther: true,
			},
		},
	}

	dispatchResultWithPostInfoFuzzOpts = combineFuzzOpts(
		optionWeightFuzzOpts,
		paysFuzzOpts,
		dispatchErrorFuzzOpts,
		[]fuzzOpt{
			withFuzzFuncs(func(d *DispatchResultWithPostInfo, c fuzz.Continue) {
				if c.RandBool() {
					d.IsOk = true
					c.Fuzz(&d.Ok)
					return
				}

				d.IsError = true
				c.Fuzz(&d.Error)
			}),
		},
	)
)

func TestDispatchResultWithPostInfo_EncodeDecode(t *testing.T) {
	assertRoundTripFuzz[DispatchResultWithPostInfo](t, 1000, dispatchResultWithPostInfoFuzzOpts...)
	assertDecodeNilData[DispatchResultWithPostInfo](t)
	assertEncodeEmptyObj[DispatchResultWithPostInfo](t, 0)
}

func TestDispatchResultWithPostInfo_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{testDispatchResultWithPostInfo1, MustHexDecodeString("0x00017b0000000000000000")},
		{testDispatchResultWithPostInfo2, MustHexDecodeString("0x0101c8010000000000000100")},
	})
}

func TestDispatchResultWithPostInfo_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{MustHexDecodeString("0x00017b0000000000000000"), testDispatchResultWithPostInfo1},
		{MustHexDecodeString("0x0101c8010000000000000100"), testDispatchResultWithPostInfo2},
	})
}
