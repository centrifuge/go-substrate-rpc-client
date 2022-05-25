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

import "github.com/ComposableFi/go-substrate-rpc-client/v4/scale"

type Payload struct {
	ID    [2]byte
	Value []byte
}

// Commitment is a beefy commitment
type Commitment struct {
	Payload        []Payload
	BlockNumber    U32
	ValidatorSetID U64
}

// SignedCommitment is a beefy commitment with optional signatures from the set of validators
type SignedCommitment struct {
	Commitment Commitment
	Signatures []OptionBeefySignature
}

type CompactSignedCommitment struct {
	Commitment        Commitment
	SignaturesFrom    []byte
	ValidatorSetLen   uint32
	SignaturesCompact []BeefySignature
}

const ContainerBitSize = 8

func (cs *CompactSignedCommitment) Unpack() SignedCommitment {
	var bits []byte

	for _, block := range cs.SignaturesFrom {
		for i := 0; i < ContainerBitSize; i++ {
			bits = append(bits, (block>>(ContainerBitSize-i-1))&1)
		}
	}

	bits = bits[:cs.ValidatorSetLen]

	var sigs []OptionBeefySignature
	count := 0
	for _, v := range bits {
		if v == 1 {
			sigs = append(sigs, OptionBeefySignature{option{true}, cs.SignaturesCompact[count]})
			count++
		} else {
			sigs = append(sigs, OptionBeefySignature{option: option{false}})

		}
	}

	return SignedCommitment{Commitment: cs.Commitment, Signatures: sigs}
}

// BeefySignature is a beefy signature
type BeefySignature [65]byte

// OptionBeefySignature is a structure that can store a BeefySignature or a missing value
type OptionBeefySignature struct {
	option
	value BeefySignature
}

// NewOptionBeefySignature creates an OptionBeefySignature with a value
func NewOptionBeefySignature(value BeefySignature) OptionBeefySignature {
	return OptionBeefySignature{option{true}, value}
}

// NewOptionBeefySignatureEmpty creates an OptionBeefySignature without a value
func NewOptionBeefySignatureEmpty() OptionBeefySignature {
	return OptionBeefySignature{option: option{false}}
}

func (o OptionBeefySignature) Encode(encoder scale.Encoder) error {
	return encoder.EncodeOption(o.HasValue, o.value)
}

func (o *OptionBeefySignature) Decode(decoder scale.Decoder) error {
	return decoder.DecodeOption(&o.HasValue, &o.value)
}

// SetSome sets a value
func (o *OptionBeefySignature) SetSome(value BeefySignature) {
	o.HasValue = true
	o.value = value
}

// SetNone removes a value and marks it as missing
func (o *OptionBeefySignature) SetNone() {
	o.HasValue = false
	o.value = BeefySignature{}
}

// Unwrap returns a flag that indicates whether a value is present and the stored value
func (o OptionBeefySignature) Unwrap() (ok bool, value BeefySignature) {
	return o.HasValue, o.value
}
