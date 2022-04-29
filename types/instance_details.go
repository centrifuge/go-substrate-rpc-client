package types

import "github.com/centrifuge/go-substrate-rpc-client/v4/scale"

type InstanceDetails struct {
	Owner    AccountID
	Approved OptionAccountID
	IsFrozen bool
	Deposit  U128
}

func (i *InstanceDetails) Decode(decoder scale.Decoder) error {
	if err := decoder.Decode(&i.Owner); err != nil {
		return err
	}
	if err := decoder.Decode(&i.Approved); err != nil {
		return err
	}
	if err := decoder.Decode(&i.IsFrozen); err != nil {
		return err
	}

	return decoder.Decode(&i.Deposit)
}

func (i InstanceDetails) Encode(encoder scale.Encoder) error {
	if err := encoder.Encode(i.Owner); err != nil {
		return err
	}
	if err := encoder.Encode(i.Approved); err != nil {
		return err
	}
	if err := encoder.Encode(i.IsFrozen); err != nil {
		return err
	}

	return encoder.Encode(i.Deposit)
}
