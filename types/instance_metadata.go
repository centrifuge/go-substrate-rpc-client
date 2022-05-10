package types

import "github.com/centrifuge/go-substrate-rpc-client/v4/scale"

type InstanceMetadata struct {
	Deposit  U128
	Data     Bytes
	IsFrozen bool
}

func (i *InstanceMetadata) Decode(decoder scale.Decoder) error {
	if err := decoder.Decode(&i.Deposit); err != nil {
		return err
	}

	if err := decoder.Decode(&i.Data); err != nil {
		return err
	}

	return decoder.Decode(&i.IsFrozen)
}

func (i InstanceMetadata) Encode(encoder scale.Encoder) error {
	if err := encoder.Encode(i.Deposit); err != nil {
		return err
	}

	if err := encoder.Encode(i.Data); err != nil {
		return err
	}

	return encoder.Encode(i.IsFrozen)
}
