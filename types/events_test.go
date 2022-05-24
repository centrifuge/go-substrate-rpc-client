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

var (
	dispatchInfoFuzzOpts = combineFuzzOpts(
		dispatchClassFuzzOpts,
		paysFuzzOpts,
	)
)

func TestDispatchInfo_EncodeDecode(t *testing.T) {
	assertRoundTripFuzz[DispatchInfo](t, 100, dispatchInfoFuzzOpts...)
	assertDecodeNilData[DispatchInfo](t)
	assertEncodeEmptyObj[DispatchInfo](t, 8)
}

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

var (
	dispatchResultFuzzOpts = combineFuzzOpts(
		dispatchErrorFuzzOpts,
		[]fuzzOpt{
			withFuzzFuncs(func(d *DispatchResult, c fuzz.Continue) {
				if c.RandBool() {
					d.Ok = true
					return
				}

				c.Fuzz(&d.Error)
			}),
		},
	)
)

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

	assertRoundTripFuzz[DispatchResult](t, 100, dispatchResultFuzzOpts...)
	assertDecodeNilData[DispatchResult](t)
	assertEncodeEmptyObj[DispatchResult](t, 1)
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
	assertDecodeNilData[Pays](t)
	assertEncodeEmptyObj[Pays](t, 0)
}

var (
	dispatchClassFuzzOpts = []fuzzOpt{
		withFuzzFuncs(func(d *DispatchClass, c fuzz.Continue) {
			switch c.Intn(3) {
			case 0:
				d.IsNormal = true
			case 1:
				d.IsOperational = true
			case 2:
				d.IsMandatory = true
			}
		}),
	}
)

func TestDispatchClassEncodeDecode(t *testing.T) {
	assertRoundTripFuzz[DispatchClass](t, 100, dispatchClassFuzzOpts...)
	assertDecodeNilData[DispatchClass](t)
	assertEncodeEmptyObj[DispatchClass](t, 0)
}

var (
	democracyConvictionFuzzOpts = []fuzzOpt{
		withFuzzFuncs(func(d *DemocracyConviction, c fuzz.Continue) {
			*d = DemocracyConviction(c.Intn(7))
		}),
	}
)

func TestDemocracyConviction_EncodeDecode(t *testing.T) {
	assertRoundTripFuzz[DemocracyConviction](t, 100, democracyConvictionFuzzOpts...)
	assertDecodeNilData[DemocracyConviction](t)
	assertEncodeEmptyObj[DemocracyConviction](t, 1)
}

var (
	voteAccountVoteFuzzOpts = combineFuzzOpts(
		democracyConvictionFuzzOpts,
		[]fuzzOpt{
			withFuzzFuncs(func(v *VoteAccountVote, c fuzz.Continue) {
				if c.RandBool() {
					v.IsStandard = true
					c.Fuzz(&v.AsStandard)
					return
				}

				v.IsSplit = true
				c.Fuzz(&v.AsSplit)
			}),
		},
	)
)

func TestVoteAccountVote_EncodeDecode(t *testing.T) {
	assertRoundTripFuzz[VoteAccountVote](t, 100, voteAccountVoteFuzzOpts...)
	assertDecodeNilData[VoteAccountVote](t)
	assertEncodeEmptyObj[VoteAccountVote](t, 0)
}

var (
	schedulerLookupErrorFuzzOpts = []fuzzOpt{
		withFuzzFuncs(func(e *SchedulerLookupError, c fuzz.Continue) {
			*e = SchedulerLookupError(c.Intn(2))
		}),
	}
)

func TestSchedulerLookupError_EncodeDecode(t *testing.T) {
	assertRoundTripFuzz[SchedulerLookupError](t, 100, schedulerLookupErrorFuzzOpts...)
	assertDecodeNilData[SchedulerLookupError](t)
	assertEncodeEmptyObj[SchedulerLookupError](t, 1)
}
