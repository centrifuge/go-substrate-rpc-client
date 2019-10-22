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
	"errors"
	"fmt"
	"io"
	"reflect"

	"github.com/centrifuge/go-substrate-rpc-client/scale"
)

// EventRecordsRaw is a raw record for a set of events, represented as the raw bytes. It exists since
// decoding of events can only be done with metadata, so events can't follow the static way of decoding
// other types do. It exposes functions to decode events using metadata and targets.
// Be careful using this in your own structs – it only works as the last value in a struct since it will consume the
// remainder of the encoded data. The reason for this is that it does not contain any length encoding, so it would
// not know where to stop.
type EventRecordsRaw []byte

// Encode implements encoding for Data, which just unwraps the bytes of Data
func (e EventRecordsRaw) Encode(encoder scale.Encoder) error {
	return encoder.Write(e)
}

// Decode implements decoding for Data, which just reads all the remaining bytes into Data
func (e *EventRecordsRaw) Decode(decoder scale.Decoder) error {
	for i := 0; true; i++ {
		b, err := decoder.ReadOneByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		*e = append((*e)[:i], b)
	}
	return nil
}

type EventSystemExtrinsicSuccess struct {
	Phase  Phase
	Topics []Hash
}
type EventSystemExtrinsicFailed struct {
	Phase         Phase
	DispatchError DispatchError
	Topics        []Hash
}

type EventBalancesTransfer struct {
	Phase  Phase
	From   AccountID
	To     AccountID
	Value  U128
	Fees   U128
	Topics []Hash
}
type EventIndicesNewAccountIndex struct {
	Phase        Phase
	AccountID    AccountID
	AccountIndex AccountIndex
	Topics       []Hash
}
type EventBalancesNewAccount struct {
	Phase     Phase
	AccountID AccountID
	Balance   U128
	Topics    []Hash
}
type EventBalancesReapedAccount struct {
	Phase     Phase
	AccountID AccountID
	Topics    []Hash
}
type EventSessionNewSession struct {
	Phase        Phase
	SessionIndex U32
	Topics       []Hash
}

// EventRecords is a default set of possible event records that can be used as a target for
// `func (e EventRecordsRaw) Decode(...`
type EventRecords struct {
	System_ExtrinsicSuccess []EventSystemExtrinsicSuccess //nolint:stylecheck,golint
	System_ExtrinsicFailed  []EventSystemExtrinsicFailed  //nolint:stylecheck,golint
	Indices_NewAccountIndex []EventIndicesNewAccountIndex //nolint:stylecheck,golint
	Balances_NewAccount     []EventBalancesNewAccount     //nolint:stylecheck,golint
	Balances_ReapedAccount  []EventBalancesReapedAccount  //nolint:stylecheck,golint
	Balances_Transfer       []EventBalancesTransfer       //nolint:stylecheck,golint
	Session_NewSession      []EventSessionNewSession      //nolint:stylecheck,golint
}

// DecodeEventRecords decodes the events records from an EventRecordRaw into a target t using the given Metadata m
func (e EventRecordsRaw) DecodeEventRecords(m *Metadata, t interface{}) error {
	// ensure t is a pointer
	ttyp := reflect.TypeOf(t)
	if ttyp.Kind() != reflect.Ptr {
		return errors.New("target must be a pointer, but is " + fmt.Sprint(ttyp))
	}
	// ensure t is not a nil pointer
	tval := reflect.ValueOf(t)
	if tval.IsNil() {
		return errors.New("target is a nil pointer")
	}
	val := tval.Elem()
	typ := val.Type()
	// ensure val can be set
	if !val.CanSet() {
		return fmt.Errorf("unsettable value %v", typ)
	}
	// ensure val points to a struct
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("target must point to a struct, but is " + fmt.Sprint(typ))
	}

	decoder := scale.NewDecoder(bytes.NewReader(e))

	// determine number of events
	n, err := decoder.DecodeUintCompact()
	if err != nil {
		return err
	}

	// iterate over events
	for i := uint64(0); i < n; i++ {
		// decode EventID
		id := EventID{}
		err := decoder.Decode(&id)
		if err != nil {
			return fmt.Errorf("unable to decode EventID for event #%v: %v", i, err)
		}

		// ask metadata for method & event name for event
		moduleName, eventName, err := m.FindEventNamesForEventID(id)
		// moduleName, eventName, err := "System", "ExtrinsicSuccess", nil
		if err != nil {
			return fmt.Errorf("unable to find event with EventID %v in metadata for event #%v", id, i)
		}

		// check whether name for eventID exists in t
		field := val.FieldByName(fmt.Sprintf("%v_%v", moduleName, eventName))
		if !field.IsValid() {
			return fmt.Errorf("unable to find field %v_%v for event #%v with EventID %v", moduleName, eventName, i, id)
		}

		// create a pointer to with the correct type that will hold the decoded event
		holder := reflect.New(field.Type().Elem())
		err = decoder.Decode(holder.Interface())
		if err != nil {
			return fmt.Errorf("unable to decode event #%v with EventID %v, field %v_%v: %v", i, id, moduleName,
				eventName, err)
		}

		// add the decoded event to the slice
		field.Set(reflect.Append(field, holder.Elem()))
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

// DispatchError is an error occurring during extrinsic dispatch
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

type EventID [2]byte
