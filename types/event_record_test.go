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
	"fmt"
	"testing"

	. "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	. "github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	. "github.com/centrifuge/go-substrate-rpc-client/v4/types/test_utils"
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
)

var examplePhaseApp = Phase{
	IsApplyExtrinsic: true,
	AsApplyExtrinsic: 42,
}

var examplePhaseFin = Phase{
	IsFinalization: true,
}

var exampleEventApp = EventSystemExtrinsicSuccess{
	Phase:        examplePhaseApp,
	DispatchInfo: DispatchInfo{Weight: testWeight, Class: DispatchClass{IsNormal: true}, PaysFee: Pays{IsYes: true}},
	Topics:       []Hash{{1, 2}},
}

var exampleEventFin = EventSystemExtrinsicSuccess{
	Phase:        examplePhaseFin,
	DispatchInfo: DispatchInfo{Weight: testWeight, Class: DispatchClass{IsNormal: true}, PaysFee: Pays{IsYes: true}},
	Topics:       []Hash{{1, 2}},
}

var exampleEventAppEnc = []byte{0x0, 0x2a, 0x0, 0x0, 0x0, 0x2c, 0xe9, 0x9, 0x0, 0x0, 0x4, 0x1, 0x2, 0x0, 0x0, 0x0, 0x0,
	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0} //nolint:lll

var exampleEventFinEnc = []byte{0x1, 0x2c, 0xe9, 0x9, 0x0, 0x0, 0x4, 0x1, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0} //nolint:lll

var (
	eventSystemExtrinsicSuccessFuzzOpts = CombineFuzzOpts(
		phaseFuzzOpts,
		dispatchInfoFuzzOpts,
	)
)

func TestEventSystemExtrinsicSuccess_EncodeDecode(t *testing.T) {
	AssertRoundTripFuzz[EventSystemExtrinsicSuccess](t, 100, eventSystemExtrinsicSuccessFuzzOpts...)
	AssertDecodeNilData[EventSystemExtrinsicSuccess](t)
	AssertEncodeEmptyObj[EventSystemExtrinsicSuccess](t, 3)
}

func TestEventSystemExtrinsicSuccess_Encode(t *testing.T) {
	encoded, err := Encode(exampleEventFin)
	assert.NoError(t, err)
	assert.Equal(t, exampleEventFinEnc, encoded)

	encoded, err = Encode(exampleEventApp)
	assert.NoError(t, err)
	assert.Equal(t, exampleEventAppEnc, encoded)
}

func TestEventSystemExtrinsicSuccess_Decode(t *testing.T) {
	decoded := EventSystemExtrinsicSuccess{}
	err := Decode(exampleEventFinEnc, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, exampleEventFin, decoded)

	decoded = EventSystemExtrinsicSuccess{}
	err = Decode(exampleEventAppEnc, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, exampleEventApp, decoded)
}

func TestEventRecordsRaw_Decode_FailsNumFields(t *testing.T) {
	e := EventRecordsRaw(MustHexDecodeString("0x0400020000000302d43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48266d00000000000000000000000000000010a5d4e8000000000000000000000000")) //nolint:lll

	events := struct {
		Balances_Transfer []struct{ Abc uint8 } //nolint:stylecheck,golint
	}{}
	err := e.DecodeEventRecords(ExamplaryMetadataV8, &events)
	assert.EqualError(t, err, "expected event #0 with EventID [3 2], field Balances_Transfer to have at least 2 fields (for Phase and Topics), but has 1 fields") //nolint:lll
}

func TestEventRecordsRaw_Decode_FailsFirstNotPhase(t *testing.T) {
	e := EventRecordsRaw(MustHexDecodeString("0x0400020000000302d43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48266d00000000000000000000000000000010a5d4e8000000000000000000000000")) //nolint:lll

	events := struct {
		Balances_Transfer []struct { //nolint:stylecheck,golint
			P     uint8
			Other uint32
			T     []Hash
		}
	}{}
	err := e.DecodeEventRecords(ExamplaryMetadataV8, &events)
	assert.EqualError(t, err, "expected the first field of event #0 with EventID [3 2], field Balances_Transfer to be of type types.Phase, but got uint8") //nolint:lll
}

func TestEventRecordsRaw_Decode_FailsLastNotHash(t *testing.T) {
	e := EventRecordsRaw(MustHexDecodeString("0x0400020000000302d43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48266d00000000000000000000000000000010a5d4e8000000000000000000000000")) //nolint:lll

	events := struct {
		Balances_Transfer []struct { //nolint:stylecheck,golint
			P     Phase
			Other uint32
			T     Phase
		}
	}{}
	err := e.DecodeEventRecords(ExamplaryMetadataV8, &events)
	assert.EqualError(t, err, "expected the last field of event #0 with EventID [3 2], field Balances_Transfer to be of type []types.Hash for Topics, but got types.Phase") //nolint:lll
}

func ExampleEventRecordsRaw_Decode() {
	e := EventRecordsRaw(MustHexDecodeString(
		"0x10" +
			"0000000000" +
			"0000" +
			"2ce909" + // Weight
			"01" + // Operational
			"01" + // PaysFee
			"00" +
			"0001000000" +
			"0000" +
			"2ce909" + // Weight
			"01" + // operational
			"01" + // PaysFee
			"00" +
			"0001000000" + // ApplyExtrinsic(1)
			"0302" + // Balances_Transfer
			"d43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d" + // From
			"8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48" + // To
			"391b0000000000000000000000000000" + // Value
			"00" + // Topics

			"0002000000" +
			"0000" +
			"2ce909" + // Weight
			"00" + // Normal
			"01" + // PaysFee
			"00",
	))

	events := EventRecords{}
	err := e.DecodeEventRecords(ExamplaryMetadataV8, &events)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Got %v System_ExtrinsicSuccess events\n", len(events.System_ExtrinsicSuccess))
	fmt.Printf("Got %v Balances_Transfer events\n", len(events.Balances_Transfer))
	t := events.Balances_Transfer[0]
	fmt.Printf("Transfer: %v tokens from %#x to\n%#x", t.Value, t.From, t.To)

	// Output: Got 3 System_ExtrinsicSuccess events
	// Got 1 Balances_Transfer events
	// Transfer: 6969 tokens from 0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d to
	// 0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48
}

func TestDispatchError(t *testing.T) {
	AssertRoundTripFuzz[DispatchError](t, 1000, dispatchErrorFuzzOpts...)
	AssertDecodeNilData[DispatchError](t)
	AssertEncodeEmptyObj[DispatchError](t, 0)
}

var (
	phaseFuzzOpts = []FuzzOpt{
		WithFuzzFuncs(func(p *Phase, c fuzz.Continue) {
			switch c.Intn(3) {
			case 0:
				p.IsApplyExtrinsic = true
				c.Fuzz(&p.AsApplyExtrinsic)
			case 1:
				p.IsFinalization = true
			case 2:
				p.IsInitialization = true
			}
		}),
	}
)

func TestPhase(t *testing.T) {
	AssertRoundTripFuzz[Phase](t, 100, phaseFuzzOpts...)
	AssertDecodeNilData[Phase](t)
	AssertEncodeEmptyObj[Phase](t, 0)
}
