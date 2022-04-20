package types

import "github.com/centrifuge/go-substrate-rpc-client/v4/scale"

type JunctionV0 struct {
	IsParent bool

	IsParachain bool
	ParachainID U32

	IsAccountId32        bool
	AccountId32NetworkID NetworkID
	AccountID            []U8

	IsAccountIndex64        bool
	AccountIndex64NetworkID NetworkID
	AccountIndex            U64

	IsAccountKey20        bool
	AccountKey20NetworkID NetworkID
	AccountKey            []U8

	IsPalletInstance bool
	PalletIndex      U8

	IsGeneralIndex bool
	GeneralIndex   U128

	IsGeneralKey bool
	GeneralKey   []U8

	IsOnlyChild bool

	IsPlurality   bool
	PluralityID   BodyID
	PluralityPart BodyPart
}

func (j *JunctionV0) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch b {
	case 0:
		j.IsParent = true
	case 1:
		j.IsParachain = true

		return decoder.Decode(&j.ParachainID)
	case 2:
		j.IsAccountId32 = true

		if err := decoder.Decode(&j.AccountId32NetworkID); err != nil {
			return err
		}

		return decoder.Decode(&j.AccountID)
	case 3:
		j.IsAccountIndex64 = true

		if err := decoder.Decode(&j.AccountIndex64NetworkID); err != nil {
			return err
		}

		return decoder.Decode(&j.AccountIndex)
	case 4:
		j.IsAccountKey20 = true

		if err := decoder.Decode(&j.AccountKey20NetworkID); err != nil {
			return err
		}

		return decoder.Decode(&j.AccountKey)
	case 5:
		j.IsPalletInstance = true

		return decoder.Decode(&j.PalletIndex)
	case 6:
		j.IsGeneralIndex = true

		return decoder.Decode(&j.GeneralIndex)
	case 7:
		j.IsGeneralKey = true

		return decoder.Decode(&j.GeneralKey)
	case 8:
		j.IsOnlyChild = true
	case 9:
		j.IsPlurality = true

		if err := decoder.Decode(&j.PluralityID); err != nil {
			return err
		}

		return decoder.Decode(&j.PluralityPart)
	}

	return nil
}

func (j JunctionV0) Encode(encoder scale.Encoder) error {
	switch {
	case j.IsParent:
		return encoder.PushByte(0)
	case j.IsParachain:
		if err := encoder.PushByte(1); err != nil {
			return err
		}

		return encoder.Encode(j.ParachainID)
	case j.IsAccountId32:
		if err := encoder.PushByte(2); err != nil {
			return err
		}

		if err := encoder.Encode(j.AccountId32NetworkID); err != nil {
			return err
		}

		return encoder.Encode(j.AccountID)
	case j.IsAccountIndex64:
		if err := encoder.PushByte(3); err != nil {
			return err
		}

		if err := encoder.Encode(j.AccountIndex64NetworkID); err != nil {
			return err
		}

		return encoder.Encode(j.AccountIndex)
	case j.IsAccountKey20:
		if err := encoder.PushByte(4); err != nil {
			return err
		}

		if err := encoder.Encode(j.AccountKey20NetworkID); err != nil {
			return err
		}

		return encoder.Encode(j.AccountKey)
	case j.IsPalletInstance:
		if err := encoder.PushByte(5); err != nil {
			return err
		}

		return encoder.Encode(j.PalletIndex)
	case j.IsGeneralIndex:
		if err := encoder.PushByte(6); err != nil {
			return err
		}

		return encoder.Encode(j.GeneralIndex)
	case j.IsGeneralKey:
		if err := encoder.PushByte(7); err != nil {
			return err
		}

		return encoder.Encode(j.GeneralKey)
	case j.IsOnlyChild:
		return encoder.PushByte(8)
	case j.IsPlurality:
		if err := encoder.PushByte(9); err != nil {
			return err
		}

		if err := encoder.Encode(j.PluralityID); err != nil {
			return err
		}

		return encoder.Encode(j.PluralityPart)
	}

	return nil
}
