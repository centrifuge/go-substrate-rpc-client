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

package types

import (
	"encoding/json"
)

// U8 is an unsigned 8-bit integer
type U8 uint8

// NewU8 creates a new U8 type
func NewU8(u uint8) U8 {
	return U8(u)
}

// UnmarshalJSON fills u with the JSON encoded byte array given by b
func (u *U8) UnmarshalJSON(b []byte) error {
	var tmp uint8
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	*u = U8(tmp)
	return nil
}

// MarshalJSON returns a JSON encoded byte array of u
func (u *U8) MarshalJSON() ([]byte, error) {
	return json.Marshal(uint8(*u))
}

// U16 is an unsigned 16-bit integer
type U16 uint16

// NewU16 creates a new U16 type
func NewU16(u uint16) U16 {
	return U16(u)
}

// UnmarshalJSON fills u with the JSON encoded byte array given by b
func (u *U16) UnmarshalJSON(b []byte) error {
	var tmp uint16
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	*u = U16(tmp)
	return nil
}

// MarshalJSON returns a JSON encoded byte array of u
func (u *U16) MarshalJSON() ([]byte, error) {
	return json.Marshal(uint16(*u))
}

// U32 is an unsigned 32-bit integer
type U32 uint32

// NewU32 creates a new U32 type
func NewU32(u uint32) U32 {
	return U32(u)
}

// UnmarshalJSON fills u with the JSON encoded byte array given by b
func (u *U32) UnmarshalJSON(b []byte) error {
	var tmp uint32
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	*u = U32(tmp)
	return nil
}

// MarshalJSON returns a JSON encoded byte array of u
func (u *U32) MarshalJSON() ([]byte, error) {
	return json.Marshal(uint32(*u))
}

// U64 is an unsigned 64-bit integer
type U64 uint64

// NewU64 creates a new U64 type
func NewU64(u uint64) U64 {
	return U64(u)
}

// UnmarshalJSON fills u with the JSON encoded byte array given by b
func (u *U64) UnmarshalJSON(b []byte) error {
	var tmp uint64
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	*u = U64(tmp)
	return nil
}

// MarshalJSON returns a JSON encoded byte array of u
func (u *U64) MarshalJSON() ([]byte, error) {
	return json.Marshal(uint64(*u))
}
