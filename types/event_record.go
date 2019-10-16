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

package types

import (
	"bytes"

	"github.com/centrifuge/go-substrate-rpc-client/scale"
)

// EventRecord is a record for an Event (as specified by Metadata) with the specific Phase of application
// type EventRecord struct {
// 	Phase  Phase
// 	Event  Event
// 	Topics []Hash
// }

// EventRecordsRaw is a raw record for a set of events, represented as the raw bytes. It exists since
// decoding of events can only be done with metadata, so events can't follow the static way of decoding
// other types do. It exposes functions to decode events using metadata and targets.
type EventRecordsRaw []byte

type EventSystemExtrinsicSuccess struct{}
type EventSystemExtrinsicFailed struct {
	DispatchError DispatchError // TODO only for V8
}
type EventBalancesTransfer struct {
	From  AccountID
	To    AccountID
	Value U128
	Fees  U128
}
type EventIndicesNewAccountIndex struct {
	AccountID    AccountID
	AccountIndex AccountIndex
}
type EventBalancesNewAccount struct {
	AccountID AccountID
	Balance   U128
}
type EventBalancesReapedAccount struct {
	AccountID AccountID
}
type EventSessionNewSession struct {
	SessionIndex U32
}

type EventRecords struct {
	System_ExtrinsicSuccess []EventSystemExtrinsicSuccess // 00 in MetadataV8
	System_ExtrinsicFailed  []EventSystemExtrinsicFailed  // 01 in MetadataV8
	Indices_NewAccountIndex []EventIndicesNewAccountIndex // 20 in MetadataV8
	Balances_NewAccount     []EventBalancesNewAccount     // 30 in MetadataV8
	Balances_ReapedAccount  []EventBalancesReapedAccount  // 31 in MetadataV8
	Balances_Transfer       []EventBalancesTransfer       // 32 in MetadataV8
	Session_NewSession      []EventSessionNewSession
}

// DecodeEvents can be used to decode the events from an EventRecordRaw into a target t using the given Metadata m
func (e EventRecordsRaw) Decode(m *Metadata, t interface{}) error {
	type Target1 struct {
		// EventID EventID
		Phase  Phase
		Event  EventSystemExtrinsicSuccess
		Topics []Hash
	}
	t1 := Target1{}

	type Target2 struct {
		// EventID EventID
		Phase  Phase
		Event  EventSystemExtrinsicFailed
		Topics []Hash
	}
	t2 := Target2{}

	type Target3 struct {
		// EventID EventID
		Phase  Phase
		Event  EventBalancesTransfer
		Topics []Hash
	}
	t3 := Target3{}

	type Target4 struct {
		// EventID EventID
		Phase  Phase
		Event  EventIndicesNewAccountIndex
		Topics []Hash
	}
	t4 := Target4{}

	decoder := scale.NewDecoder(bytes.NewReader(e))

	// determine number of events
	n, err := decoder.DecodeUintCompact()
	if err != nil {
		return err
	}

	// tar := make([]Target, n)

	// iterate over events
	for i := uint64(0); i < n; i++ {
		// decode EventID
		id := EventID{}
		err := decoder.Decode(&id)
		if err != nil {
			return err
		}

		if i == 0 {
			err = decoder.Decode(&t1)
		} else if i == 1 {
			err = decoder.Decode(&t2)
		} else if i == 2 {
			err = decoder.Decode(&t3)
		} else if i == 3 {
			err = decoder.Decode(&t4)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// Phase is an enum describing the current phase of the event (applying the extrinsic or finalized)
type Phase struct {
	IsApplyExtrinsic bool
	AsApplyExtrinsic uint32
	IsFinalization   bool
}

func (p *Phase) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	if b == 0 {
		p.IsApplyExtrinsic = true
		err = decoder.Decode(&p.AsApplyExtrinsic)
	} else if b == 1 {
		p.IsFinalization = true
	}

	if err != nil {
		return err
	}

	return nil
}

func (p Phase) Encode(encoder scale.Encoder) error {
	var err1, err2 error
	if p.IsApplyExtrinsic {
		err1 = encoder.PushByte(0)
		err2 = encoder.Encode(p.AsApplyExtrinsic)
	} else if p.IsFinalization {
		err1 = encoder.PushByte(1)
	}

	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}

	return nil
}

// DispatchError is an error occuring during extrinsic dispatch
type DispatchError struct {
	HasModule bool
	Module    uint8
	Error     uint8
}

func (d *DispatchError) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	if b == 1 {
		d.HasModule = true
		err = decoder.Decode(&d.Module)
	}
	if err != nil {
		return err
	}

	return decoder.Decode(&d.Error)
}

func (d DispatchError) Encode(encoder scale.Encoder) error {
	var err error
	if d.HasModule {
		err = encoder.PushByte(1)
		if err != nil {
			return err
		}
		err = encoder.Encode(d.Module)
	} else {
		err = encoder.PushByte(0)
	}

	if err != nil {
		return err
	}

	return encoder.Encode(&d.Error)
}

// type EventAndTopicsRaw struct {
// 	EventID       EventID
// 	DataAndTopics Data
// }

// type Event struct {
// 	EventID       EventID
// 	DataAndTopics Data
// }

type EventID [2]byte

// // Decode implements decoding for EventAndTopicsRaw, which just reads all the remaining bytes into EventAndTopicsRaw
// func (e *EventAndTopicsRaw) Decode(decoder scale.Decoder) error {
// 	for {
// 		b, err := decoder.ReadOneByte()
// 		if err == io.EOF {
// 			// fmt.Println(err)
// 			break
// 		}
// 		if err != nil {
// 			return err
// 		}
// 		*e = append(*e, b)
// 	}
// 	return nil
// }

// // Encode implements encoding for Data, which just unwraps the bytes of Data
// func (e EventAndTopicsRaw) Encode(encoder scale.Encoder) error {
// 	return encoder.Write(e)
// }

// type Event struct {
// 	// Section string
// 	// Method  string
// 	// TypeDef []string
// 	Index EventID
// 	// Data  []byte
// 	// Data interface{}
// }

// type EventID [2]byte

// type EventData []byte
