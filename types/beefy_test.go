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
	"github.com/stretchr/testify/assert"
)

func TestBeefySignature(t *testing.T) {
	empty := NewOptionBeefySignatureEmpty()
	assert.True(t, empty.IsNone())
	assert.False(t, empty.IsSome())

	sig := NewOptionBeefySignature(BeefySignature{})
	sig.SetNone()
	assert.True(t, sig.IsNone())
	sig.SetSome(BeefySignature{})
	assert.True(t, sig.IsSome())
	ok, _ := sig.Unwrap()
	assert.True(t, ok)
	assertRoundtrip(t, sig)
}

func TestBeefySignature_EncodeDecode(t *testing.T) {
	assertRoundTripFuzz[BeefySignature](t, 100)
	assertDecodeNilData[BeefySignature](t)
	assertEncodeEmptyObj[BeefySignature](t, 65)
}

var (
	optionBeefySignatureFuzzOpts = []fuzzOpt{
		withFuzzFuncs(func(o *OptionBeefySignature, c fuzz.Continue) {
			if c.RandBool() {
				*o = NewOptionBeefySignatureEmpty()
				return
			}

			var b BeefySignature
			c.Fuzz(&b)

			*o = NewOptionBeefySignature(b)
		}),
	}
)

func TestOptionBeefySignature_EncodeDecode(t *testing.T) {
	assertRoundTripFuzz[OptionBeefySignature](t, 100, optionBeefySignatureFuzzOpts...)
	assertEncodeEmptyObj[OptionBeefySignature](t, 1)
}
