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

	. "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	. "github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	. "github.com/centrifuge/go-substrate-rpc-client/v4/types/test_utils"
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
)

var (
	optionWeightFuzzOpts = []FuzzOpt{
		WithFuzzFuncs(func(o *OptionWeight, c fuzz.Continue) {
			if c.RandBool() {
				*o = NewOptionWeightEmpty()
				return
			}

			var weight Weight

			c.Fuzz(&weight)

			*o = NewOptionWeight(weight)
		}),
	}
)

func TestOptionWeight_EncodeDecode(t *testing.T) {
	AssertRoundTripFuzz[OptionWeight](t, 100, optionWeightFuzzOpts...)
	AssertEncodeEmptyObj[OptionWeight](t, 1)
}

func TestOptionWeight_Encode(t *testing.T) {
	AssertEncode(t, []EncodingAssert{
		{NewOptionWeight(NewWeight(0)), MustHexDecodeString("0x010000000000000000")},
		{NewOptionWeight(NewWeight(1)), MustHexDecodeString("0x010100000000000000")},
		{NewOptionWeight(NewWeight(2)), MustHexDecodeString("0x010200000000000000")},
		{NewOptionWeightEmpty(), MustHexDecodeString("0x00")},
	})
}

func TestOptionWeight_Decode(t *testing.T) {
	AssertDecode(t, []DecodingAssert{
		{MustHexDecodeString("0x010000000000000000"), NewOptionWeight(NewWeight(0))},
		{MustHexDecodeString("0x010100000000000000"), NewOptionWeight(NewWeight(1))},
		{MustHexDecodeString("0x010200000000000000"), NewOptionWeight(NewWeight(2))},
		{MustHexDecodeString("0x00"), NewOptionWeightEmpty()},
	})
}

func TestOptionWeight_OptionMethods(t *testing.T) {
	o := NewOptionWeightEmpty()
	o.SetSome(Weight(11))

	ok, v := o.Unwrap()
	assert.True(t, ok)
	assert.NotNil(t, v)

	o.SetNone()

	ok, v = o.Unwrap()
	assert.False(t, ok)
	assert.Equal(t, Weight(0), v)
}

func TestWeight_EncodeDecode(t *testing.T) {
	AssertRoundTripFuzz[Weight](t, 100)
	AssertDecodeNilData[Weight](t)
	AssertEncodeEmptyObj[Weight](t, 8)
}

func TestWeight_EncodedLength(t *testing.T) {
	AssertEncodedLength(t, []EncodedLengthAssert{{NewWeight(13), 8}})
}

func TestWeight_Encode(t *testing.T) {
	AssertEncode(t, []EncodingAssert{
		{NewWeight(29), MustHexDecodeString("0x1d00000000000000")},
	})
}

func TestWeight_Hash(t *testing.T) {
	AssertHash(t, []HashAssert{
		{NewWeight(29), MustHexDecodeString("0x83e168a13a013e6d47b0778f046aaa05d6c01d6857d044d9e9b658a6d85eb865")},
	})
}

func TestWeight_Hex(t *testing.T) {
	AssertEncodeToHex(t, []EncodeToHexAssert{
		{NewWeight(29), "0x1d00000000000000"},
	})
}

func TestWeight_String(t *testing.T) {
	AssertString(t, []StringAssert{
		{NewWeight(29), "29"},
	})
}

func TestWeight_Eq(t *testing.T) {
	AssertEq(t, []EqAssert{
		{NewWeight(23), NewWeight(23), true},
		{NewWeight(23), NewBool(false), false},
	})
}
