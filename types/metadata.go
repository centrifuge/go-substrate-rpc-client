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
	"hash"

	"github.com/centrifuge/go-substrate-rpc-client/scale"
)

const MagicNumber uint32 = 0x6174656d

// Modelled after https://github.com/paritytech/substrate/blob/v1.0.0rc2/srml/metadata/src/lib.rs

type Metadata struct {
	MagicNumber  uint32
	Version      uint8
	IsMetadataV4 bool
	AsMetadataV4 MetadataV4
	IsMetadataV7 bool
	AsMetadataV7 MetadataV7
	IsMetadataV8 bool
	AsMetadataV8 MetadataV8
}

func NewMetadataV4() *Metadata {
	return &Metadata{Version: 4, IsMetadataV4: true, AsMetadataV4: MetadataV4{make([]ModuleMetadataV4, 0)}}
}

func NewMetadataV7() *Metadata {
	return &Metadata{Version: 7, IsMetadataV7: true, AsMetadataV7: MetadataV7{make([]ModuleMetadataV7, 0)}}
}

func NewMetadataV8() *Metadata {
	return &Metadata{Version: 8, IsMetadataV8: true, AsMetadataV8: MetadataV8{make([]ModuleMetadataV8, 0)}}
}

func (m *Metadata) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&m.MagicNumber)
	if err != nil {
		return err
	}
	if m.MagicNumber != MagicNumber {
		return fmt.Errorf("magic number mismatch: expected %#x, found %#x", MagicNumber, m.MagicNumber)
	}

	err = decoder.Decode(&m.Version)
	if err != nil {
		return err
	}

	switch m.Version {
	case 4:
		m.IsMetadataV4 = true
		err = decoder.Decode(&m.AsMetadataV4)
	case 7:
		m.IsMetadataV7 = true
		err = decoder.Decode(&m.AsMetadataV7)
	case 8:
		m.IsMetadataV8 = true
		err = decoder.Decode(&m.AsMetadataV8)
	default:
		return fmt.Errorf("unsupported metadata version %v", m.Version)
	}

	return err
}

func (m Metadata) Encode(encoder scale.Encoder) error {
	err := encoder.Encode(m.MagicNumber)
	if err != nil {
		return err
	}

	err = encoder.Encode(m.Version)
	if err != nil {
		return err
	}

	switch m.Version {
	case 4:
		err = encoder.Encode(m.AsMetadataV4)
	case 7:
		err = encoder.Encode(m.AsMetadataV7)
	case 8:
		err = encoder.Encode(m.AsMetadataV8)
	default:
		return fmt.Errorf("unsupported metadata version %v", m.Version)
	}

	return err
}

func (m *Metadata) FindCallIndex(call string) (CallIndex, error) {
	switch {
	case m.IsMetadataV4:
		return m.AsMetadataV4.FindCallIndex(call)
	case m.IsMetadataV7:
		return m.AsMetadataV7.FindCallIndex(call)
	case m.IsMetadataV8:
		return m.AsMetadataV8.FindCallIndex(call)
	default:
		return CallIndex{}, fmt.Errorf("unsupported metadata version")
	}
}

func (m *Metadata) FindEventNamesForEventID(eventID EventID) (Text, Text, error) {
	switch {
	case m.IsMetadataV4:
		return m.AsMetadataV4.FindEventNamesForEventID(eventID)
	case m.IsMetadataV7:
		return m.AsMetadataV7.FindEventNamesForEventID(eventID)
	case m.IsMetadataV8:
		return m.AsMetadataV8.FindEventNamesForEventID(eventID)
	default:
		return "", "", fmt.Errorf("unsupported metadata version")
	}
}

func (m *Metadata) FindStorageKeyHasher(module string, fn string) (hash.Hash, error) {
	switch {
	case m.IsMetadataV4:
		return m.AsMetadataV4.FindStorageKeyHasher(module, fn)
	case m.IsMetadataV7:
		return m.AsMetadataV7.FindStorageKeyHasher(module, fn)
	case m.IsMetadataV8:
		return m.AsMetadataV8.FindStorageKeyHasher(module, fn)
	default:
		return nil, fmt.Errorf("unsupported metadata version")
	}
}
