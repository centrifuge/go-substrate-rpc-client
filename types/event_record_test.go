// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
// Copyright (C) 2019  Centrifuge GmbH
//
// This file is part of Go Substrate RPC Client (GSRPC).
//
// GSRPC is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// GSRPC is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package types_test

import (
	"fmt"

	. "github.com/centrifuge/go-substrate-rpc-client/types"
)

// var examplePhaseApp = Phase{
// 	IsApplyExtrinsic: true,
// 	AsApplyExtrinsic: 42,
// }

// var examplePhaseFin = Phase{
// 	IsFinalization: true,
// }

// var exampleEvent = Event{
// 	Index: EventID{1, 2},
// 	// Data:  Data{0x03, 0x04, 0x05},
// }

// var exampleEventRecordApp = EventRecord{
// 	Phase:  examplePhaseApp,
// 	Event:  exampleEvent,
// 	Topics: []Hash{{1, 2}},
// }

// var exampleEventRecordFin = EventRecord{
// 	Phase:  examplePhaseFin,
// 	Event:  exampleEvent,
// 	Topics: []Hash{{1, 2}},
// }

// var exampleEventRecordFinEnc = []byte{0x1, 0x8, 0x1, 0x2, 0x3, 0x4, 0x5, 0x4, 0x80, 0x1, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}

// var exampleEventRecordRawApp = EventRecordRaw{
// 	Phase:             examplePhaseApp,
// 	EventAndTopicsRaw: []byte{1, 2, 3},
// }

// var exampleEventRecordRawFin = EventRecordRaw{
// 	Phase:             examplePhaseFin,
// 	EventAndTopicsRaw: []byte{1, 2, 3},
// }

// var exampleEventRecordRawFinEnc = []byte{0x1, 0x1, 0x2, 0x3}

// func TestEventRecord_EncodedLength(t *testing.T) {
// 	assertEncodedLength(t, []encodedLengthAssert{
// 		{exampleEventRecordApp, 45},
// 		{exampleEventRecordFin, 41},
// 	})
// }

// func TestEventRecord_Encode(t *testing.T) {
// 	encoded, err := EncodeToBytes(exampleEventRecordFin)
// 	assert.NoError(t, err)
// 	assert.Equal(t, exampleEventRecordFinEnc, encoded)
// }

// func TestEventRecord_Decode(t *testing.T) {
// 	decoded := EventRecord{}
// 	err := DecodeFromBytes(exampleEventRecordFinEnc, &decoded)
// 	assert.NoError(t, err)
// 	assert.Equal(t, exampleEventRecordFin, decoded)
// }

// func TestEventRecord_Hash(t *testing.T) {
// 	assertHash(t, []hashAssert{
// 		{exampleEventRecordFin, mustDecodeHexString(
// 			"0xc5b3444e6d277f1cc07246a16fe1ff5aa54d4aee174ca0c963b0aad28e1cb765")},
// 	})
// }

// func TestEventRecord_Hex(t *testing.T) {
// 	assertEncodeToHex(t, []encodeToHexAssert{
// 		{exampleEventRecordFin, "0x0108010203040504800102000000000000000000000000000000000000000000000000000000000000"},
// 	})
// }

// func TestEventRecord_String(t *testing.T) {
// 	assertString(t, []stringAssert{
// 		{exampleEventRecordFin, "{{false 0 true} {[1 2] [3 4 5]} [[1 2 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]]}"},
// 	})
// }

// func TestEventRecord_Eq(t *testing.T) {
// 	assertEq(t, []eqAssert{
// 		{exampleEventRecordFin, exampleEventRecordFin, true},
// 		{exampleEventRecordApp, exampleEventRecordFin, false},
// 		{exampleEventRecordApp, NewBool(true), false},
// 	})
// }

// func TestEventRecord_EncodedLength(t *testing.T) {
// 	assertEncodedLength(t, []encodedLengthAssert{
// 		{exampleEventRecordRawApp, 45},
// 		{exampleEventRecordRawFin, 41},
// 	})
// }

// func TestEventRecordRaw_Encode(t *testing.T) {
// 	encoded, err := EncodeToBytes(exampleEventRecordRawFin)
// 	assert.NoError(t, err)
// 	assert.Equal(t, exampleEventRecordRawFinEnc, encoded)
// }

// func TestEventRecordRaw_Decode(t *testing.T) {
// 	decoded := EventRecordRaw{}
// 	err := DecodeFromBytes(exampleEventRecordRawFinEnc, &decoded)
// 	assert.NoError(t, err)
// 	assert.Equal(t, exampleEventRecordRawFin, decoded)
// }

// func TestEventRecordRaw_Hash(t *testing.T) {
// 	assertHash(t, []hashAssert{
// 		{exampleEventRecordRawFin, mustDecodeHexString(
// 			"0xc5b3444e6d277f1cc07246a16fe1ff5aa54d4aee174ca0c963b0aad28e1cb765")},
// 	})
// }

// func TestEventRecordRaw_Hex(t *testing.T) {
// 	assertEncodeToHex(t, []encodeToHexAssert{
// 		{exampleEventRecordRawFin, "0x0108010203040504800102000000000000000000000000000000000000000000000000000000000000"},
// 	})
// }

// func TestEventRecordRaw_String(t *testing.T) {
// 	assertString(t, []stringAssert{
// 		{exampleEventRecordRawFin, "{{false 0 true} {[1 2] [3 4 5]} [[1 2 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]]}"},
// 	})
// }

// func TestEventRecordRaw_Eq(t *testing.T) {
// 	assertEq(t, []eqAssert{
// 		{exampleEventRecordRawFin, exampleEventRecordRawFin, true},
// 		{exampleEventRecordRawApp, exampleEventRecordRawFin, false},
// 		{exampleEventRecordRawApp, NewBool(true), false},
// 	})
// }

func ExampleEventRecordsRaw_Decode() {
	e := EventRecordsRaw(MustHexDecodeString("0x100000000000000000000100000000000000020000000302d43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48266d00000000000000000000000000000010a5d4e80000000000000000000000000002000000000000")) //nolint:lll

	events := EventRecords{}

	err := e.Decode(ExamplaryMetadataV8, &events)
	fmt.Println(err)

	fmt.Printf("%#v\n", events)

	// Output: "abc"
}
