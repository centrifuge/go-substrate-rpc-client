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
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v2/scale"
)

// ExtrinsicEra indicates either a mortal or immortal extrinsic
type ExtrinsicEra struct {
	IsImmortalEra bool
	// AsImmortalEra ImmortalEra
	IsMortalEra bool
	AsMortalEra MortalEra
}

// MortalEra for an extrinsic, indicating period and phase
type MortalEra struct {
	First  U64
	Second U64
}

func (e *ExtrinsicEra) Decode(decoder scale.Decoder) error {
	tag, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch tag {
	case 0:
		e.IsImmortalEra = true
	case 1:
		e.IsMortalEra = true

		err = decoder.Decode(&e.AsMortalEra.First)
		if err != nil {
			return err
		}

		err = decoder.Decode(&e.AsMortalEra.Second)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("No such variant for ExtrinsicEra")
	}

	return nil
}

func (e ExtrinsicEra) Encode(encoder scale.Encoder) error {
	var err error
	switch {
	case e.IsImmortalEra:
		err = encoder.PushByte(0)
		if err != nil {
			return err
		}
	case e.IsMortalEra:
		err = encoder.PushByte(1)
		if err != nil {
			return err
		}

		err = encoder.Encode(e.AsMortalEra.First)
		if err != nil {
			return err
		}

		err = encoder.Encode(e.AsMortalEra.Second)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("No such variant for ExtrinsicEra")
	}

	return nil
}
