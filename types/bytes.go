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

// Bytes represents byte slices
type Bytes []byte

// NewBytes creates a new Bytes type
func NewBytes(b []byte) Bytes {
	return Bytes(b)
}

// Bytes8 represents an 8 byte array
type Bytes8 [8]byte

// NewBytes8 creates a new Bytes8 type
func NewBytes8(b [8]byte) Bytes8 {
	return Bytes8(b)
}

// Bytes16 represents an 16 byte array
type Bytes16 [16]byte

// NewBytes16 creates a new Bytes16 type
func NewBytes16(b [16]byte) Bytes16 {
	return Bytes16(b)
}

// Bytes32 represents an 32 byte array
type Bytes32 [32]byte

// NewBytes32 creates a new Bytes32 type
func NewBytes32(b [32]byte) Bytes32 {
	return Bytes32(b)
}

// Bytes64 represents an 64 byte array
type Bytes64 [64]byte

// NewBytes64 creates a new Bytes64 type
func NewBytes64(b [64]byte) Bytes64 {
	return Bytes64(b)
}

// Bytes128 represents an 128 byte array
type Bytes128 [128]byte

// NewBytes128 creates a new Bytes128 type
func NewBytes128(b [128]byte) Bytes128 {
	return Bytes128(b)
}

// Bytes256 represents an 256 byte array
type Bytes256 [256]byte

// NewBytes256 creates a new Bytes256 type
func NewBytes256(b [256]byte) Bytes256 {
	return Bytes256(b)
}

// Bytes512 represents an 512 byte array
type Bytes512 [512]byte

// NewBytes512 creates a new Bytes512 type
func NewBytes512(b [512]byte) Bytes512 {
	return Bytes512(b)
}

// Bytes1024 represents an 1024 byte array
type Bytes1024 [1024]byte

// NewBytes1024 creates a new Bytes1024 type
func NewBytes1024(b [1024]byte) Bytes1024 {
	return Bytes1024(b)
}

// Bytes2048 represents an 2048 byte array
type Bytes2048 [2048]byte

// NewBytes2048 creates a new Bytes2048 type
func NewBytes2048(b [2048]byte) Bytes2048 {
	return Bytes2048(b)
}
