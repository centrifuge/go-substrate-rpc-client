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

type NetworkIDV3 struct {
	IsByGenesis bool
	ByGenesis Hash

	IsByFork bool
	ByForkBlockNumber U64
	ByForkBlockHash Hash

	IsPolkadot bool
	IsKusama bool
	IsWestend bool
	IsRococo bool
	IsWococo bool

	IsEthereum bool
	Ethereum UCompact

	IsBitcoinCore bool
	IsBitcoinCash bool
}

func (n *NetworkIDV3) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch b {
	case 0:
		n.IsByGenesis = true
		return decoder.Decode(&n.ByGenesis)
	case 1:
		n.IsByFork = true

		if err := decoder.Decode(&n.ByForkBlockNumber); err != nil {
			return err
		}

		return decoder.Decode(&n.ByForkBlockHash)
	case 2:
		n.IsPolkadot = true
	case 3:
		n.IsKusama = true
	case 4:
		n.IsWestend = true
	case 5:
		n.IsRococo = true
	case 6:
		n.IsWococo = true
	case 7:
		n.IsEthereum = true
		return decoder.Decode(&n.Ethereum)
	case 8:
		n.IsBitcoinCore = true
	case 9:
		n.IsBitcoinCash = true
	}

	return nil
}

func (n NetworkIDV3) Encode(encoder scale.Encoder) error {
	switch {
	case n.IsByGenesis:
		if err := encoder.PushByte(0); err != nil {
			return err
		}
		return encoder.Encode(n.ByGenesis)
	case n.IsByFork:
		if err := encoder.PushByte(1); err != nil {
			return err
		}

		if err := encoder.Encode(n.ByForkBlockNumber); err != nil {
			return err
		}

		return encoder.Encode(n.ByForkBlockHash)
	case n.IsPolkadot:
		return encoder.PushByte(2)
	case n.IsKusama:
		return encoder.PushByte(3)
	case n.IsWestend:
		return encoder.PushByte(4)
	case n.IsRococo:
		return encoder.PushByte(5)
	case n.IsWococo:
		return encoder.PushByte(6)
	case n.IsEthereum:
		if err := encoder.PushByte(7); err != nil {
			return err
		}
		return encoder.Encode(n.Ethereum)
	case n.IsBitcoinCore:
		return encoder.PushByte(8)
	case n.IsBitcoinCash:
		return encoder.PushByte(9)
	}

	return nil
}
