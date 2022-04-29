package types

import "github.com/centrifuge/go-substrate-rpc-client/v4/scale"

type ClassDetails struct {
	Owner             AccountID
	Issuer            AccountID
	Admin             AccountID
	Freezer           AccountID
	TotalDeposit      U128
	FreeHolding       bool
	Instances         U32
	InstanceMetadatas U32
	Attributes        U32
	IsFrozen          bool
}

func (c *ClassDetails) Decode(decoder scale.Decoder) error {
	if err := decoder.Decode(&c.Owner); err != nil {
		return err
	}

	if err := decoder.Decode(&c.Issuer); err != nil {
		return err
	}
	if err := decoder.Decode(&c.Admin); err != nil {
		return err
	}
	if err := decoder.Decode(&c.Freezer); err != nil {
		return err
	}
	if err := decoder.Decode(&c.TotalDeposit); err != nil {
		return err
	}
	if err := decoder.Decode(&c.FreeHolding); err != nil {
		return err
	}
	if err := decoder.Decode(&c.Instances); err != nil {
		return err
	}
	if err := decoder.Decode(&c.InstanceMetadatas); err != nil {
		return err
	}
	if err := decoder.Decode(&c.Attributes); err != nil {
		return err
	}

	return decoder.Decode(&c.IsFrozen)
}

func (c ClassDetails) Encode(encoder scale.Encoder) error {
	if err := encoder.Encode(c.Owner); err != nil {
		return err
	}

	if err := encoder.Encode(c.Issuer); err != nil {
		return err
	}

	if err := encoder.Encode(c.Admin); err != nil {
		return err
	}

	if err := encoder.Encode(c.Freezer); err != nil {
		return err
	}

	if err := encoder.Encode(c.TotalDeposit); err != nil {
		return err
	}

	if err := encoder.Encode(c.FreeHolding); err != nil {
		return err
	}

	if err := encoder.Encode(c.Instances); err != nil {
		return err
	}

	if err := encoder.Encode(c.InstanceMetadatas); err != nil {
		return err
	}

	if err := encoder.Encode(c.Attributes); err != nil {
		return err
	}

	return encoder.Encode(c.IsFrozen)
}
