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

import "github.com/centrifuge/go-substrate-rpc-client/scale"

type ExtrinsicSignatureV3 struct {
	Signer    Address
	Signature Signature
	Era       ExtrinsicEra // extra via system::CheckEra
	Nonce     UCompact     // extra via system::CheckNonce (Compact<Index> where Index is u32))
	Tip       UCompact     // extra via balances::TakeFees (Compact<Balance> where Balance is u128))
}

type ExtrinsicSignatureV4 struct {
	Signer    Address
	Signature MultiSignature
	Era       ExtrinsicEra // extra via system::CheckEra
	Nonce     UCompact     // extra via system::CheckNonce (Compact<Index> where Index is u32))
	Tip       UCompact     // extra via balances::TakeFees (Compact<Balance> where Balance is u128))
}

type SignatureOptions struct {
	Era         ExtrinsicEra // extra via system::CheckEra
	Nonce       UCompact     // extra via system::CheckNonce (Compact<Index> where Index is u32)
	Tip         UCompact     // extra via balances::TakeFees (Compact<Balance> where Balance is u128)
	SpecVersion U32          // additional via system::CheckVersion
	GenesisHash Hash         // additional via system::CheckGenesis
	BlockHash   Hash         // additional via system::CheckEra
}

func (es *ExtrinsicSignatureV4) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}
	if b == 0xff {
		err = decoder.Decode(&es.Signer.AsAccountID)
		if err != nil {
			return err
		}
		es.Signer.IsAccountID = true
	} else {
		err := decoder.Decode(&es.Signer)
		if err != nil {
			return err
		}
	}

	err = decoder.Decode(&es.Signature)
	if err != nil {
		return err
	}

	err = decoder.Decode(&es.Era)
	if err != nil {
		return err
	}

	err = decoder.Decode(&es.Nonce)
	if err != nil {
		return err
	}

	err = decoder.Decode(&es.Tip)
	if err != nil {
		return err
	}

	return nil
}
