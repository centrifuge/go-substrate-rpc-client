package types

import "github.com/centrifuge/go-substrate-rpc-client/v4/scale"

const (
	// CentrifugeChainStringLimit is defined in the centrifuge-chain.
	CentrifugeChainStringLimit = 256
)

type ClassMetadata struct {
	Deposit  U128
	Data     [CentrifugeChainStringLimit]U8
	IsFrozen bool
}

func (c *ClassMetadata) Decode(decoder scale.Decoder) error {
	if err := decoder.Decode(&c.Deposit); err != nil {
		return err
	}

	if err := decoder.Decode(&c.Data); err != nil {
		return err
	}

	return decoder.Decode(&c.IsFrozen)
}

func (c ClassMetadata) Encode(encoder scale.Encoder) error {
	if err := encoder.Encode(c.Deposit); err != nil {
		return err
	}

	if err := encoder.Encode(c.Data); err != nil {
		return err
	}

	return encoder.Encode(c.IsFrozen)
}
