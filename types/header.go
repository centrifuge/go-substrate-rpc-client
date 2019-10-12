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

	"github.com/centrifuge/go-substrate-rpc-client/scale"
)

type Header struct {
	ParentHash     Hash        `json:"parentHash"`
	Number         BlockNumber `json:"number"`
	StateRoot      Hash        `json:"stateRoot"`
	ExtrinsicsRoot Hash        `json:"extrinsicsRoot"`
	Digest         Digest      `json:"digest"`
}

type BlockNumber U32

// UnmarshalJSON fills u with the JSON encoded byte array given by b
func (b *BlockNumber) UnmarshalJSON(bz []byte) error {
	var tmp string
	if err := json.Unmarshal(bz, &tmp); err != nil {
		return err
	}
	return DecodeFromHexString(tmp, b)
}

// MarshalJSON returns a JSON encoded byte array of u
func (b BlockNumber) MarshalJSON() ([]byte, error) {
	s, err := EncodeToHexString(b)
	if err != nil {
		return nil, err
	}
	return json.Marshal(s)
}

// Encode implements encoding for BlockNumber, which just unwraps the bytes of BlockNumber
func (b BlockNumber) Encode(encoder scale.Encoder) error {
	return encoder.EncodeUintCompact(uint64(b))
}

// Decode implements decoding for BlockNumber, which just wraps the bytes in BlockNumber
func (b *BlockNumber) Decode(decoder scale.Decoder) error {
	u, err := decoder.DecodeUintCompact()
	if err != nil {
		return err
	}
	*b = BlockNumber(u)
	return err
}

// Digest contains logs
type Digest []DigestItem

// UnmarshalJSON fills u with the JSON encoded byte array given by b
func (d *Digest) UnmarshalJSON(bz []byte) error {
	var tmp struct {
		Logs []string `json:"logs"`
	}
	if err := json.Unmarshal(bz, &tmp); err != nil {
		return err
	}
	*d = make([]DigestItem, len(tmp.Logs))
	for i, log := range tmp.Logs {
		err := DecodeFromHexString(log, &(*d)[i])
		if err != nil {
			return err
		}
	}
	return nil
}

// MarshalJSON returns a JSON encoded byte array of u
func (d Digest) MarshalJSON() ([]byte, error) {
	logs := make([]string, len(d))
	var err error
	for i, di := range d {
		logs[i], err = EncodeToHexString(di)
		if err != nil {
			return nil, err
		}
	}
	return json.Marshal(struct {
		Logs []string `json:"logs"`
	}{
		Logs: logs,
	})
}

// DigestItem speciefies the item in the logs of a digest
type DigestItem struct {
	IsOther             bool
	AsOther             Bytes // 0
	IsAuthoritiesChange bool
	AsAuthoritiesChange []AuthorityID // 1
	IsChangesTrieRoot   bool
	AsChangesTrieRoot   Hash // 2
	IsSealV0            bool
	AsSealV0            SealV0 // 3
	IsConsensus         bool
	AsConsensus         Consensus // 4
	IsSeal              bool
	AsSeal              Seal // 5
	IsPreRuntime        bool
	AsPreRuntime        PreRuntime // 6
}

func (m *DigestItem) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()

	if err != nil {
		return err
	}

	switch b {
	case 0:
		m.IsOther = true
		err = decoder.Decode(&m.AsOther)
	case 1:
		m.IsAuthoritiesChange = true
		err = decoder.Decode(&m.AsAuthoritiesChange)
	case 2:
		m.IsChangesTrieRoot = true
		err = decoder.Decode(&m.AsChangesTrieRoot)
	case 3:
		m.IsSealV0 = true
		err = decoder.Decode(&m.AsSealV0)
	case 4:
		m.IsConsensus = true
		err = decoder.Decode(&m.AsConsensus)
	case 5:
		m.IsSeal = true
		err = decoder.Decode(&m.AsSeal)
	case 6:
		m.IsPreRuntime = true
		err = decoder.Decode(&m.AsPreRuntime)
	}

	if err != nil {
		return err
	}

	return nil
}

func (m DigestItem) Encode(encoder scale.Encoder) error {
	var err1, err2 error
	switch {
	case m.IsOther:
		err1 = encoder.PushByte(0)
		err2 = encoder.Encode(m.AsOther)
	case m.IsAuthoritiesChange:
		err1 = encoder.PushByte(1)
		err2 = encoder.Encode(m.AsAuthoritiesChange)
	case m.IsChangesTrieRoot:
		err1 = encoder.PushByte(2)
		err2 = encoder.Encode(m.AsChangesTrieRoot)
	case m.IsSealV0:
		err1 = encoder.PushByte(3)
		err2 = encoder.Encode(m.AsSealV0)
	case m.IsConsensus:
		err1 = encoder.PushByte(4)
		err2 = encoder.Encode(m.AsConsensus)
	case m.IsSeal:
		err1 = encoder.PushByte(5)
		err2 = encoder.Encode(m.AsSeal)
	case m.IsPreRuntime:
		err1 = encoder.PushByte(6)
		err2 = encoder.Encode(m.AsPreRuntime)
	}

	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}

	return nil
}

// AuthorityID represents a public key (an 32 byte array)
type AuthorityID [32]byte

// NewAuthorityID creates a new AuthorityID type
func NewAuthorityID(b [32]byte) AuthorityID {
	return AuthorityID(b)
}

type SealV0 struct {
	Signer    U64
	Signature Signature
}

type Seal struct {
	ConsensusEngineID ConsensusEngineID
	Bytes             Bytes
}

// ConsensusEngineID is a 4-byte identifier (actually a [u8; 4]) identifying the engine, e.g. for Aura it would be
// [b'a', b'u', b'r', b'a']
type ConsensusEngineID U32

type Consensus struct {
	ConsensusEngineID ConsensusEngineID
	Bytes             Bytes
}

type PreRuntime struct {
	ConsensusEngineID ConsensusEngineID
	Bytes             Bytes
}
