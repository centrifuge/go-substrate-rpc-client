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

type OptionMultiLocationV3 struct {
	option
	value MultiLocationV3
}

func NewOptionMultiLocationV3(value MultiLocationV3) OptionMultiLocationV3 {
	return OptionMultiLocationV3{option{hasValue: true}, value}
}

func NewOptionMultiLocationV3Empty() OptionMultiLocationV3 {
	return OptionMultiLocationV3{option: option{hasValue: false}}
}

func (o *OptionMultiLocationV3) Decode(decoder scale.Decoder) error {
	return decoder.DecodeOption(&o.hasValue, &o.value)
}

func (o OptionMultiLocationV3) Encode(encoder scale.Encoder) error {
	return encoder.EncodeOption(o.hasValue, o.value)
}

// SetSome sets a value
func (o *OptionMultiLocationV3) SetSome(value MultiLocationV3) {
	o.hasValue = true
	o.value = value
}

// SetNone removes a value and marks it as missing
func (o *OptionMultiLocationV3) SetNone() {
	o.hasValue = false
	o.value = MultiLocationV3{}
}

// Unwrap returns a flag that indicates whether a value is present and the stored value
func (o *OptionMultiLocationV3) Unwrap() (ok bool, value MultiLocationV3) {
	return o.hasValue, o.value
}

type MultiLocationV3 struct {
	Parents  U8
	Interior JunctionsV3
}

func (m *MultiLocationV3) Decode(decoder scale.Decoder) error {
	if err := decoder.Decode(&m.Parents); err != nil {
		return err
	}

	return decoder.Decode(&m.Interior)
}

func (m *MultiLocationV3) Encode(encoder scale.Encoder) error {
	if err := encoder.Encode(&m.Parents); err != nil {
		return err
	}

	return encoder.Encode(&m.Interior)
}
