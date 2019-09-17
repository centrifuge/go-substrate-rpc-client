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
	"reflect"

	"github.com/centrifuge/go-substrate-rpc-client/scale"
	. "github.com/centrifuge/go-substrate-rpc-client/types"
)

// PhaseEnum is an enum example. Since Go has no enums, it is implemented as a struct with flags for each
// potential value and a corresponding value if needed. This enum represents phases, the values are `ApplyExtrinsic`
// which can be a uint32 and `Finalization`. By implementing Encode and Decode methods on this struct that satisfy the
// scale.Encodeable and scale.Decodeable interfaces, we encode our enum struct to correspond to the scale codec
// (see https://substrate.dev/docs/en/overview/low-level-data-format for a description).
type PhaseEnum struct {
	IsApplyExtrinsic bool
	AsApplyExtrinsic uint32
	IsFinalization   bool
}

func (m *PhaseEnum) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()

	if err != nil {
		return err
	}

	if b == 0 {
		m.IsApplyExtrinsic = true
		err = decoder.Decode(&m.AsApplyExtrinsic)
	} else if b == 1 {
		m.IsFinalization = true
	}

	if err != nil {
		return err
	}

	return nil
}

func (m PhaseEnum) Encode(encoder scale.Encoder) error {
	var err1, err2 error
	if m.IsApplyExtrinsic {
		err1 = encoder.PushByte(0)
		err2 = encoder.Encode(m.AsApplyExtrinsic)
	} else if m.IsFinalization {
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

func ExampleExampleEnum_applyExtrinsic() {
	applyExtrinsic := PhaseEnum{
		IsApplyExtrinsic: true,
		AsApplyExtrinsic: 1234,
	}

	enc, err := EncodeToHexString(applyExtrinsic)
	if err != nil {
		panic(err)
	}

	var dec PhaseEnum
	err = DecodeFromHexString(enc, &dec)
	if err != nil {
		panic(err)
	}

	fmt.Println(reflect.DeepEqual(applyExtrinsic, dec))
}

func ExampleExampleEnum_finalization() {
	finalization := PhaseEnum{
		IsFinalization: true,
	}

	enc, err := EncodeToHexString(finalization)
	if err != nil {
		panic(err)
	}

	var dec PhaseEnum
	err = DecodeFromHexString(enc, &dec)
	if err != nil {
		panic(err)
	}

	fmt.Println(reflect.DeepEqual(finalization, dec))
}
