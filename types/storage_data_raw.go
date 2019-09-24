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
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/scale"
)

// StorageDataRaw contains raw bytes that are not decoded/encoded
type StorageDataRaw []byte

// NewStorageDataRaw creates a new StorageDataRaw type
func NewStorageDataRaw(b []byte) StorageDataRaw {
	return StorageDataRaw(b)
}

// Encode implements encoding for StorageDataRaw, which just unwraps the bytes of StorageDataRaw
func (s StorageDataRaw) Encode(encoder scale.Encoder) error {
	return encoder.Write(s)
}

// Decode implements decoding for StorageDataRaw, which just wraps the bytes in StorageDataRaw
func (s *StorageDataRaw) Decode(decoder scale.Decoder) error {
	return decoder.Read(*s)
}

// Hex returns a hex string representation of the value
func (s StorageDataRaw) Hex() string {
	return fmt.Sprintf("%#x", s)
}
