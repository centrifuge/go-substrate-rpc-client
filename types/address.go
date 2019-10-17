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

// Address is a wrapper around an AccountId or an AccountIndex. It is encoded with a prefix in case of an AccountID.
// Basically the Address is encoded as `[ <prefix-byte>, ...publicKey/...bytes ]` as per spec
type Address struct {
	IsAccountID    bool
	AsAccountID    AccountID
	IsAccountIndex bool
	AsAccountIndex AccountIndex
}

// NewAddressFromAccountID creates an Address from the given AccountID (public key)
func NewAddressFromAccountID(b []byte) Address {
	return Address{
		IsAccountID: true,
		AsAccountID: NewAccountID(b),
	}
}

// NewAddressFromHexAccountID creates an Address from the given hex string that contains an AccountID (public key)
func NewAddressFromHexAccountID(str string) (Address, error) {
	b, err := HexDecodeString(str)
	if err != nil {
		return Address{}, err
	}
	return NewAddressFromAccountID(b), nil
}

// NewAddressFromAccountIndex creates an Address from the given AccountIndex
func NewAddressFromAccountIndex(u uint32) Address {
	return Address{
		IsAccountIndex: true,
		AsAccountIndex: AccountIndex(u),
	}
}

func (a *Address) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	if b == 0xff {
		err = decoder.Decode(&a.AsAccountID)
		a.IsAccountID = true
		return err
	}

	if b == 0xfe {
		return fmt.Errorf("decoding of Address with 0xfe prefix not supported")
	}

	if b == 0xfd {
		err = decoder.Decode(&a.AsAccountIndex)
		a.IsAccountIndex = true
		return err
	}

	if b == 0xfc {
		var aIndex uint16
		err = decoder.Decode(&aIndex)
		a.IsAccountIndex = true
		a.AsAccountIndex = AccountIndex(aIndex)
		return err
	}

	a.IsAccountIndex = true
	a.AsAccountIndex = AccountIndex(b)
	return nil
}

func (a Address) Encode(encoder scale.Encoder) error {
	// type of address - public key
	if a.IsAccountID {
		err := encoder.PushByte(255)
		if err != nil {
			return err
		}

		err = encoder.Write(a.AsAccountID[:])
		if err != nil {
			return err
		}

		return nil
	}

	if a.AsAccountIndex > 0xffff {
		err := encoder.PushByte(253)
		if err != nil {
			return err
		}

		return encoder.Encode(a.AsAccountIndex)
	}

	if a.AsAccountIndex >= 0xf0 {
		err := encoder.PushByte(252)
		if err != nil {
			return err
		}

		return encoder.Encode(uint16(a.AsAccountIndex))
	}

	return encoder.Encode(uint8(a.AsAccountIndex))
}
