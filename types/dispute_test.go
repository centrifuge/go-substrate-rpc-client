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
	testDisputeLocation1 = DisputeLocation{
		IsLocal: true,
	}
	testDisputeLocation2 = DisputeLocation{
		IsRemote: true,
	}

	testDisputeResult1 = DisputeResult{
		IsValid: true,
	}

	testDisputeResult2 = DisputeResult{
		IsInvalid: true,
	}
)

var (
	disputeLocationFuzzOpts = []fuzzOpt{
		withFuzzFuncs(func(d *DisputeLocation, c fuzz.Continue) {
			if c.RandBool() {
				d.IsLocal = true
				return
			}

			d.IsRemote = true
		}),
	}
)

func TestDisputeLocation_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, testDisputeLocation1)
	assertRoundtrip(t, testDisputeLocation2)
	assertRoundTripFuzz[DisputeLocation](t, 100, disputeLocationFuzzOpts...)
	assertDecodeNilData[DisputeLocation](t)
	assertEncodeEmptyObj[DisputeLocation](t, 0)
}

func TestDisputeLocation_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{testDisputeLocation1, MustHexDecodeString("0x00")},
		{testDisputeLocation2, MustHexDecodeString("0x01")},
	})
}

func TestDisputeLocation_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{MustHexDecodeString("0x00"), testDisputeLocation1},
		{MustHexDecodeString("0x01"), testDisputeLocation2},
	})
}

var (
	disputeResultFuzzOpts = []fuzzOpt{
		withFuzzFuncs(func(d *DisputeResult, c fuzz.Continue) {
			if c.RandBool() {
				d.IsValid = true
				return
			}

			d.IsInvalid = true
			return
		}),
	}
)

func TestDisputeResult_EncodeDecode(t *testing.T) {
	assertRoundtrip(t, testDisputeResult1)
	assertRoundtrip(t, testDisputeResult2)
	assertRoundTripFuzz[DisputeResult](t, 100, disputeResultFuzzOpts...)
	assertDecodeNilData[DisputeResult](t)
	assertEncodeEmptyObj[DisputeResult](t, 0)
}

func TestDisputeResult_Encode(t *testing.T) {
	assertEncode(t, []encodingAssert{
		{testDisputeResult1, MustHexDecodeString("0x00")},
		{testDisputeResult2, MustHexDecodeString("0x01")},
	})
}

func TestDisputeResult_Decode(t *testing.T) {
	assertDecode(t, []decodingAssert{
		{MustHexDecodeString("0x00"), testDisputeResult1},
		{MustHexDecodeString("0x01"), testDisputeResult2},
	})
}
