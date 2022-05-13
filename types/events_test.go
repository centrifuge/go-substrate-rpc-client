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
	"bytes"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	. "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
)

func TestVoteThreshold_Decoder(t *testing.T) {
	// SuperMajorityAgainst
	decoder := scale.NewDecoder(bytes.NewReader([]byte{1}))
	vt := VoteThreshold(0)
	err := decoder.Decode(&vt)
	assert.NoError(t, err)
	assert.Equal(t, vt, SuperMajorityAgainst)

	// Error
	decoder = scale.NewDecoder(bytes.NewReader([]byte{3}))
	err = decoder.Decode(&vt)
	assert.Error(t, err)
}

func TestVoteThreshold_Encode(t *testing.T) {
	vt := SuperMajorityAgainst
	var buf bytes.Buffer
	encoder := scale.NewEncoder(&buf)
	assert.NoError(t, encoder.Encode(vt))
	assert.Equal(t, buf.Len(), 1)
	assert.Equal(t, buf.Bytes(), []byte{1})
}

func TestDispatchResult_Decode(t *testing.T) {
	// ok
	decoder := scale.NewDecoder(bytes.NewReader([]byte{0}))
	var res DispatchResult
	err := decoder.Decode(&res)
	assert.NoError(t, err)
	assert.True(t, res.Ok)

	// Dispatch Error
	decoder = scale.NewDecoder(bytes.NewReader([]byte{1, 3, 1, 1}))
	res = DispatchResult{}
	assert.NoError(t, decoder.Decode(&res))

	assert.False(t, res.Ok)
	assert.True(t, res.Error.IsModule)
	assert.Equal(t, res.Error.ModuleError.Index, U8(1))
	assert.Equal(t, res.Error.ModuleError.Error, U8(1))

	// decoder error
	decoder = scale.NewDecoder(bytes.NewReader([]byte{1, 3, 1}))
	res = DispatchResult{}
	assert.Error(t, decoder.Decode(&res))
}

func TestProxyTypeEncodeDecode(t *testing.T) {
	// encode
	pt := Governance
	var buf bytes.Buffer
	encoder := scale.NewEncoder(&buf)
	assert.NoError(t, encoder.Encode(pt))
	assert.Equal(t, buf.Len(), 1)
	assert.Equal(t, buf.Bytes(), []byte{2})

	//decode
	decoder := scale.NewDecoder(bytes.NewReader(buf.Bytes()))
	pt0 := ProxyType(0)
	err := decoder.Decode(&pt0)
	assert.NoError(t, err)
	assert.Equal(t, pt0, Governance)

	//decode error
	decoder = scale.NewDecoder(bytes.NewReader([]byte{5}))
	pt0 = ProxyType(0)
	err = decoder.Decode(&pt0)
	assert.Error(t, err)
}

var (
	paysFuzzOpts = []fuzzOpt{
		withFuzzFuncs(func(p *Pays, c fuzz.Continue) {
			r := c.RandBool()
			p.IsYes = r
			p.IsNo = !r
		}),
	}
)

func TestPaysEncodeDecode(t *testing.T) {
	assertRoundTripFuzz[Pays](t, 1000, paysFuzzOpts...)
}

func TestDispatchClassEncodeDecode(t *testing.T) {
	// encode
	dc := DispatchClass{IsMandatory: true}
	var buf bytes.Buffer
	encoder := scale.NewEncoder(&buf)
	assert.NoError(t, encoder.Encode(dc))
	assert.Equal(t, buf.Len(), 1)
	assert.Equal(t, buf.Bytes(), []byte{2})

	// decode supported
	var dcc DispatchClass
	decoder := scale.NewDecoder(bytes.NewReader(buf.Bytes()))
	err := decoder.Decode(&dcc)
	assert.NoError(t, err)
	assert.True(t, dcc.IsMandatory)
}
