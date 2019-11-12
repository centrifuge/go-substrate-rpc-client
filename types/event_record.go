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

// EventBalancesNewAccount is emitted when a new account was created
type EventBalancesNewAccount struct {
	Phase     Phase
	AccountID AccountID
	Balance   U128
	Topics    []Hash
}

// EventBalancesReapedAccount is emitted when an account was reaped
type EventBalancesReapedAccount struct {
	Phase     Phase
	AccountID AccountID
	Topics    []Hash
}

// EventBalancesTransfer is emitted when a transfer succeeded (from, to, value, fees)
type EventBalancesTransfer struct {
	Phase  Phase
	From   AccountID
	To     AccountID
	Value  U128
	Fees   U128
	Topics []Hash
}

// EventGrandpaNewAuthorities is emitted when a new authority set has been applied
type EventGrandpaNewAuthorities struct {
	Phase          Phase
	NewAuthorities []struct {
		AuthorityID     AuthorityID
		AuthorityWeight U64
	}
	Topics []Hash
}

// EventGrandpaPaused is emitted when the current authority set has been paused
type EventGrandpaPaused struct {
	Phase  Phase
	Topics []Hash
}

// EventGrandpaResumed is emitted when the current authority set has been resumed
type EventGrandpaResumed struct {
	Phase  Phase
	Topics []Hash
}

// EventImOnlineHeartbeatReceived is emitted when a new heartbeat was received from AuthorityId
type EventImOnlineHeartbeatReceived struct {
	Phase       Phase
	AuthorityID AuthorityID
	Topics      []Hash
}

// EventIndicesNewAccountIndex is emitted when a new account index was assigned. This event is not triggered
// when an existing index is reassigned to another AccountId
type EventIndicesNewAccountIndex struct {
	Phase        Phase
	AccountID    AccountID
	AccountIndex AccountIndex
	Topics       []Hash
}

// EventOffencesOffence is emitted when there is an offence reported of the given kind happened at the session_index
// and (kind-specific) time slot. This event is not deposited for duplicate slashes
type EventOffencesOffence struct {
	Phase          Phase
	Kind           Bytes16
	OpaqueTimeSlot Bytes
	Topics         []Hash
}

// EventSessionNewSession is emitted when a new session has happened. Note that the argument is the session index,
// not the block number as the type might suggest
type EventSessionNewSession struct {
	Phase        Phase
	SessionIndex U32
	Topics       []Hash
}

// EventStakingOldSlashingReportDiscarded is emitted when an old slashing report from a prior era was discarded because
// it could not be processed
type EventStakingOldSlashingReportDiscarded struct {
	Phase        Phase
	SessionIndex U32
	Topics       []Hash
}

// EventStakingReward is emitted when all validators have been rewarded by the given balance
type EventStakingReward struct {
	Phase   Phase
	Balance U128
	Topics  []Hash
}

// EventStakingSlash is emitted when one validator (and its nominators) has been slashed by the given amount
type EventStakingSlash struct {
	Phase     Phase
	AccountID AccountID
	Balance   U128
	Topics    []Hash
}

// EventSystemExtrinsicSuccess is emitted when an extrinsic completed successfully
type EventSystemExtrinsicSuccess struct {
	Phase  Phase
	Topics []Hash
}

// EventSystemExtrinsicFailed is emitted when an extrinsic failed
type EventSystemExtrinsicFailed struct {
	Phase         Phase
	DispatchError DispatchError
	Topics        []Hash
}

// EventRecords is a default set of possible event records that can be used as a target for
// `func (e EventRecordsRaw) Decode(...`
type EventRecords struct {
	Balances_NewAccount                []EventBalancesNewAccount                //nolint:stylecheck,golint
	Balances_ReapedAccount             []EventBalancesReapedAccount             //nolint:stylecheck,golint
	Balances_Transfer                  []EventBalancesTransfer                  //nolint:stylecheck,golint
	Grandpa_NewAuthorities             []EventGrandpaNewAuthorities             //nolint:stylecheck,golint
	Grandpa_Paused                     []EventGrandpaPaused                     //nolint:stylecheck,golint
	Grandpa_Resumed                    []EventGrandpaResumed                    //nolint:stylecheck,golint
	ImOnline_HeartbeatReceived         []EventImOnlineHeartbeatReceived         //nolint:stylecheck,golint
	Indices_NewAccountIndex            []EventIndicesNewAccountIndex            //nolint:stylecheck,golint
	Offences_Offence                   []EventOffencesOffence                   //nolint:stylecheck,golint
	Session_NewSession                 []EventSessionNewSession                 //nolint:stylecheck,golint
	Staking_OldSlashingReportDiscarded []EventStakingOldSlashingReportDiscarded //nolint:stylecheck,golint
	Staking_Reward                     []EventStakingReward                     //nolint:stylecheck,golint
	Staking_Slash                      []EventStakingSlash                      //nolint:stylecheck,golint
	System_ExtrinsicSuccess            []EventSystemExtrinsicSuccess            //nolint:stylecheck,golint
	System_ExtrinsicFailed             []EventSystemExtrinsicFailed             //nolint:stylecheck,golint
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
		// decode Phase
		phase := Phase{}
		err := decoder.Decode(&phase)
		if err != nil {
			return fmt.Errorf("unable to decode Phase for event #%v: %v", i, err)
		}

		// decode EventID
		id := EventID{}
		err = decoder.Decode(&id)
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

		// ensure first field is for Phase, last field is for Topics
		numFields := holder.Elem().NumField()
		if numFields < 2 {
			return fmt.Errorf("expected event #%v with EventID %v, field %v_%v to have at least 2 fields "+
				"(for Phase and Topics), but has %v fields", i, id, moduleName, eventName, numFields)
		}
		phaseField := holder.Elem().FieldByIndex([]int{0})
		if phaseField.Type() != reflect.TypeOf(phase) {
			return fmt.Errorf("expected the first field of event #%v with EventID %v, field %v_%v to be of type "+
				"types.Phase, but got %v", i, id, moduleName, eventName, phaseField.Type())
		}
		topicsField := holder.Elem().FieldByIndex([]int{numFields - 1})
		if topicsField.Type() != reflect.TypeOf([]Hash{}) {
			return fmt.Errorf("expected the last field of event #%v with EventID %v, field %v_%v to be of type "+
				"[]types.Hash for Topics, but got %v", i, id, moduleName, eventName, topicsField.Type())
		}

		// set the phase we decoded earlier
		phaseField.Set(reflect.ValueOf(phase))

		// set the remaining fields
		for j := 1; j < numFields; j++ {
			err = decoder.Decode(holder.Elem().FieldByIndex([]int{j}).Addr().Interface())
			if err != nil {
				return fmt.Errorf("unable to decode field %v event #%v with EventID %v, field %v_%v: %v", j, i, id, moduleName,
					eventName, err)
			}
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
