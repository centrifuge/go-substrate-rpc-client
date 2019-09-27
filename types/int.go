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
	"encoding/json"
)

// I8 is a signed 8-bit integer
type I8 int8

// NewI8 creates a new I8 type
func NewI8(i int8) I8 {
	return I8(i)
}

// UnmarshalJSON fills i with the JSON encoded byte array given by b
func (i *I8) UnmarshalJSON(b []byte) error {
	var tmp int8
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	*i = I8(tmp)
	return nil
}

// MarshalJSON returns a JSON encoded byte array of i
func (i *I8) MarshalJSON() ([]byte, error) {
	return json.Marshal(int8(*i))
}

// I16 is a signed 16-bit integer
type I16 int16

// NewI16 creates a new 16 type
func NewI16(i int16) I16 {
	return I16(i)
}

// UnmarshalJSON fills i with the JSON encoded byte array given by b
func (i *I16) UnmarshalJSON(b []byte) error {
	var tmp int16
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	*i = I16(tmp)
	return nil
}

// MarshalJSON returns a JSON encoded byte array of i
func (i *I16) MarshalJSON() ([]byte, error) {
	return json.Marshal(int16(*i))
}

// I32 is a signed 32-bit integer
type I32 int32

// NewI32 creates a new 32 type
func NewI32(i int32) I32 {
	return I32(i)
}

// UnmarshalJSON fills i with the JSON encoded byte array given by b
func (i *I32) UnmarshalJSON(b []byte) error {
	var tmp int32
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	*i = I32(tmp)
	return nil
}

// MarshalJSON returns a JSON encoded byte array of i
func (i *I32) MarshalJSON() ([]byte, error) {
	return json.Marshal(int32(*i))
}

// I64 is a signed 64-bit integer
type I64 int64

// NewI64 creates a new 64 type
func NewI64(i int64) I64 {
	return I64(i)
}

// UnmarshalJSON fills i with the JSON encoded byte array given by b
func (i *I64) UnmarshalJSON(b []byte) error {
	var tmp int64
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	*i = I64(tmp)
	return nil
}

// MarshalJSON returns a JSON encoded byte array of i
func (i *I64) MarshalJSON() ([]byte, error) {
	return json.Marshal(int64(*i))
}
