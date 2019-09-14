// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
// Copyright (C) 2019  Philip Stanislaus, Philip Stehlik, Vimukthi Wickramasinghe
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

package types

import "github.com/centrifuge/go-substrate-rpc-client/scale"

// OptionH160 is a structure that can store a H160 or a missing value
type OptionH160 struct {
	option
	value H160
}

// NewOptionH160 creates an OptionH160 with a value
func NewOptionH160(value H160) OptionH160 {
	return OptionH160{option{true}, value}
}

// NewOptionH160Empty creates an OptionH160 without a value
func NewOptionH160Empty() OptionH160 {
	return OptionH160{option: option{false}}
}

func (o OptionH160) Encode(encoder scale.Encoder) error {
	return encoder.EncodeOption(o.hasValue, o.value)
}

func (o *OptionH160) Decode(decoder scale.Decoder) error {
	return decoder.DecodeOption(&o.hasValue, &o.value)
}

// SetSome sets a value
func (o *OptionH160) SetSome(value H160) {
	o.hasValue = true
	o.value = value
}

// SetNone removes a value and marks it as missing
func (o *OptionH160) SetNone() {
	o.hasValue = false
	o.value = H160{}
}

// Unwrap returns a flag that indicates whether a value is present and the stored value
func (o OptionH160) Unwrap() (ok bool, value H160) {
	return o.hasValue, o.value
}

// OptionH256 is a structure that can store a H256 or a missing value
type OptionH256 struct {
	option
	value H256
}

// NewOptionH256 creates an OptionH256 with a value
func NewOptionH256(value H256) OptionH256 {
	return OptionH256{option{true}, value}
}

// NewOptionH256Empty creates an OptionH256 without a value
func NewOptionH256Empty() OptionH256 {
	return OptionH256{option: option{false}}
}

func (o OptionH256) Encode(encoder scale.Encoder) error {
	return encoder.EncodeOption(o.hasValue, o.value)
}

func (o *OptionH256) Decode(decoder scale.Decoder) error {
	return decoder.DecodeOption(&o.hasValue, &o.value)
}

// SetSome sets a value
func (o *OptionH256) SetSome(value H256) {
	o.hasValue = true
	o.value = value
}

// SetNone removes a value and marks it as missing
func (o *OptionH256) SetNone() {
	o.hasValue = false
	o.value = H256{}
}

// Unwrap returns a flag that indicates whether a value is present and the stored value
func (o OptionH256) Unwrap() (ok bool, value H256) {
	return o.hasValue, o.value
}

// OptionH512 is a structure that can store a H512 or a missing value
type OptionH512 struct {
	option
	value H512
}

// NewOptionH512 creates an OptionH512 with a value
func NewOptionH512(value H512) OptionH512 {
	return OptionH512{option{true}, value}
}

// NewOptionH512Empty creates an OptionH512 without a value
func NewOptionH512Empty() OptionH512 {
	return OptionH512{option: option{false}}
}

func (o OptionH512) Encode(encoder scale.Encoder) error {
	return encoder.EncodeOption(o.hasValue, o.value)
}

func (o *OptionH512) Decode(decoder scale.Decoder) error {
	return decoder.DecodeOption(&o.hasValue, &o.value)
}

// SetSome sets a value
func (o *OptionH512) SetSome(value H512) {
	o.hasValue = true
	o.value = value
}

// SetNone removes a value and marks it as missing
func (o *OptionH512) SetNone() {
	o.hasValue = false
	o.value = H512{}
}

// Unwrap returns a flag that indicates whether a value is present and the stored value
func (o OptionH512) Unwrap() (ok bool, value H512) {
	return o.hasValue, o.value
}

// OptionHash is a structure that can store a Hash or a missing value
type OptionHash struct {
	option
	value Hash
}

// NewOptionHash creates an OptionHash with a value
func NewOptionHash(value Hash) OptionHash {
	return OptionHash{option{true}, value}
}

// NewOptionHashEmpty creates an OptionHash without a value
func NewOptionHashEmpty() OptionHash {
	return OptionHash{option: option{false}}
}

func (o OptionHash) Encode(encoder scale.Encoder) error {
	return encoder.EncodeOption(o.hasValue, o.value)
}

func (o *OptionHash) Decode(decoder scale.Decoder) error {
	return decoder.DecodeOption(&o.hasValue, &o.value)
}

// SetSome sets a value
func (o *OptionHash) SetSome(value Hash) {
	o.hasValue = true
	o.value = value
}

// SetNone removes a value and marks it as missing
func (o *OptionHash) SetNone() {
	o.hasValue = false
	o.value = Hash{}
}

// Unwrap returns a flag that indicates whether a value is present and the stored value
func (o OptionHash) Unwrap() (ok bool, value Hash) {
	return o.hasValue, o.value
}
