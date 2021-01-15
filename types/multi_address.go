package types

import (
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v2/scale"
)

type MultiAddress struct {
	IsID        bool
	AsID        AccountID
	IsIndex     bool
	AsIndex     AccountIndex
	IsRaw       bool
	AsRaw       []byte
	IsAddress32 bool
	AsAddress32 [32]byte
	IsAddress20 bool
	AsAddress20 [20]byte
}

// NewMultiAddressFromAccountID creates an Address from the given AccountID (public key)
func NewMultiAddressFromAccountID(b []byte) MultiAddress {
	return MultiAddress{
		IsID: true,
		AsID: NewAccountID(b),
	}
}

// NewMultiAddressFromHexAccountID creates an Address from the given hex string that contains an AccountID (public key)
func NewMultiAddressFromHexAccountID(str string) (MultiAddress, error) {
	b, err := HexDecodeString(str)
	if err != nil {
		return MultiAddress{}, err
	}
	return NewMultiAddressFromAccountID(b), nil
}

func (m MultiAddress) Encode(encoder scale.Encoder) error {
	var err error
	switch {
	case m.IsID:
		err = encoder.PushByte(0)
		if err != nil {
			return err
		}
		err = encoder.Encode(m.AsID)
		if err != nil {
			return err
		}
	case m.IsIndex:
		err = encoder.PushByte(1)
		if err != nil {
			return err
		}
		err = encoder.Encode(m.AsIndex)
		if err != nil {
			return err
		}
	case m.IsRaw:
		err = encoder.PushByte(2)
		if err != nil {
			return err
		}
		err = encoder.Encode(m.AsRaw)
		if err != nil {
			return err
		}
	case m.IsAddress32:
		err = encoder.PushByte(3)
		if err != nil {
			return err
		}
		err = encoder.Encode(m.AsAddress32)
		if err != nil {
			return err
		}
	case m.IsAddress20:
		err = encoder.PushByte(4)
		if err != nil {
			return err
		}
		err = encoder.Encode(m.AsAddress20)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Invalid variant for MultiAddress")
	}

	return nil
}

func (m MultiAddress) Decode(decoder scale.Decoder) error {
	tag, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch tag {
	case 0:
		m.IsID = true
		err = decoder.Decode(&m.AsID)
	case 1:
		m.IsIndex = true
		err = decoder.Decode(&m.AsIndex)
	case 2:
		m.IsRaw = true
		err = decoder.Decode(&m.AsRaw)
	case 3:
		m.IsAddress32 = true
		err = decoder.Decode(&m.AsAddress32)
	case 4:
		m.IsAddress20 = true
		err = decoder.Decode(&m.AsAddress20)
	default:
		return fmt.Errorf("Invalid variant for MultiAddress")
	}

	if err != nil {
		return err
	}

	return nil
}
