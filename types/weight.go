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

import "github.com/centrifuge/go-substrate-rpc-client/v4/scale"

type OptionWeight struct {
	option
	value Weight
}

// Weight is a numeric range of a transaction weight
type Weight uint64

// NewWeight creates a new Weight type
func NewWeight(u uint64) Weight {
	return Weight(u)
}

func NewOptionWeight(value Weight) OptionWeight {
	return OptionWeight{
		option: option{
			hasValue: true,
		},
		value: value,
	}
}

func NewOptionWeightEmpty() OptionWeight {
	return OptionWeight{
		option: option{
			hasValue: false,
		},
	}
}

func (o OptionWeight) Encode(encoder scale.Encoder) error {
	return encoder.EncodeOption(o.hasValue, o.value)
}

func (o *OptionWeight) Decode(decoder scale.Decoder) error {
	return decoder.DecodeOption(&o.hasValue, &o.value)
}

// SetSome sets a value
func (o *OptionWeight) SetSome(value Weight) {
	o.hasValue = true
	o.value = value
}

// SetNone removes a value and marks it as missing
func (o *OptionWeight) SetNone() {
	o.hasValue = false
	o.value = Weight(0)
}

// Unwrap returns a flag that indicates whether a value is present and the stored value
func (o *OptionWeight) Unwrap() (ok bool, value Weight) {
	return o.hasValue, o.value
}
