package types

import "github.com/centrifuge/go-substrate-rpc-client/v4/scale"

type BodyID struct {
	IsUnit bool

	IsNamed bool
	Body    []U8

	IsIndex bool
	Index   U32

	IsExecutive bool

	IsTechnical bool

	IsLegislative bool

	IsJudicial bool
}

func (b *BodyID) Decode(decoder scale.Decoder) error {
	bb, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch bb {
	case 0:
		b.IsUnit = true
	case 1:
		b.IsNamed = true

		return decoder.Decode(&b.Body)
	case 2:
		b.IsIndex = true

		return decoder.Decode(&b.Index)
	case 3:
		b.IsExecutive = true
	case 4:
		b.IsTechnical = true
	case 5:
		b.IsLegislative = true
	case 6:
		b.IsJudicial = true
	}

	return nil
}

func (b BodyID) Encode(encoder scale.Encoder) error {
	switch {
	case b.IsUnit:
		return encoder.PushByte(0)
	case b.IsNamed:
		if err := encoder.PushByte(1); err != nil {
			return err
		}

		return encoder.Encode(b.Body)
	case b.IsIndex:
		if err := encoder.PushByte(2); err != nil {
			return err
		}

		return encoder.Encode(b.Index)
	case b.IsExecutive:
		return encoder.PushByte(3)
	case b.IsTechnical:
		return encoder.PushByte(4)
	case b.IsLegislative:
		return encoder.PushByte(5)
	case b.IsJudicial:
		return encoder.PushByte(6)
	}

	return nil
}

type BodyPart struct {
	IsVoice bool

	IsMembers    bool
	MembersCount U32

	IsFraction bool
	Nom        U32
	Denom      U32

	IsAtLeastProportion bool
	// Also contains Nom
	// Also contains Denom

	IsMoreThanProportion bool
	// Also contains Nom
	// Also contains Denom
}

func (b *BodyPart) Decode(decoder scale.Decoder) error {
	bb, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch bb {
	case 0:
		b.IsVoice = true
	case 1:
		b.IsMembers = true

		return decoder.Decode(&b.MembersCount)
	case 2:
		b.IsFraction = true

		if err := decoder.Decode(&b.Nom); err != nil {
			return err
		}

		return decoder.Decode(&b.Denom)
	case 3:
		b.IsAtLeastProportion = true

		if err := decoder.Decode(&b.Nom); err != nil {
			return err
		}

		return decoder.Decode(&b.Denom)
	case 4:
		b.IsMoreThanProportion = true

		if err := decoder.Decode(&b.Nom); err != nil {
			return err
		}

		return decoder.Decode(&b.Denom)
	}

	return nil
}

func (b BodyPart) Encode(encoder scale.Encoder) error {
	switch {
	case b.IsVoice:
		return encoder.PushByte(0)
	case b.IsMembers:
		if err := encoder.PushByte(1); err != nil {
			return err
		}

		return encoder.Encode(b.MembersCount)
	case b.IsFraction:
		if err := encoder.PushByte(2); err != nil {
			return err
		}

		if err := encoder.Encode(b.Nom); err != nil {
			return err
		}

		return encoder.Encode(b.Denom)
	case b.IsAtLeastProportion:
		if err := encoder.PushByte(3); err != nil {
			return err
		}

		if err := encoder.Encode(b.Nom); err != nil {
			return err
		}

		return encoder.Encode(b.Denom)
	case b.IsMoreThanProportion:
		if err := encoder.PushByte(4); err != nil {
			return err
		}

		if err := encoder.Encode(b.Nom); err != nil {
			return err
		}

		return encoder.Encode(b.Denom)
	}

	return nil
}
