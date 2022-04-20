package types

import "github.com/centrifuge/go-substrate-rpc-client/v4/scale"

type NetworkID struct {
	IsAny bool

	IsNamed      bool
	NamedNetwork []U8

	IsPolkadot bool

	IsKusama bool
}

func (n *NetworkID) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch b {
	case 0:
		n.IsAny = true
	case 1:
		n.IsNamed = true

		return decoder.Decode(&n.NamedNetwork)
	case 2:
		n.IsPolkadot = true
	case 3:
		n.IsKusama = true
	}

	return nil
}

func (n NetworkID) Encode(encoder scale.Encoder) error {
	switch {
	case n.IsAny:
		return encoder.PushByte(0)
	case n.IsNamed:
		if err := encoder.PushByte(1); err != nil {
			return err
		}

		return encoder.Encode(n.NamedNetwork)
	case n.IsPolkadot:
		return encoder.PushByte(2)
	case n.IsKusama:
		return encoder.PushByte(3)
	}

	return nil
}
