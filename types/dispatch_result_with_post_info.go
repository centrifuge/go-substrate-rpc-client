package types

import "github.com/centrifuge/go-substrate-rpc-client/v4/scale"

// PostDispatchInfo is used in DispatchResultWithPostInfo.
// Weight information that is only available post dispatch.
type PostDispatchInfo struct {
	ActualWeight OptionWeight
	PaysFee      Pays
}

func (p *PostDispatchInfo) Decode(decoder scale.Decoder) error {
	if err := decoder.Decode(&p.ActualWeight); err != nil {
		return err
	}

	return decoder.Decode(&p.PaysFee)
}

func (p PostDispatchInfo) Encode(encoder scale.Encoder) error {
	if err := encoder.Encode(p.ActualWeight); err != nil {
		return err
	}

	return encoder.Encode(p.PaysFee)
}

// DispatchErrorWithPostInfo is used in DispatchResultWithPostInfo.
type DispatchErrorWithPostInfo struct {
	PostInfo PostDispatchInfo
	Error    DispatchError
}

func (d *DispatchErrorWithPostInfo) Decode(decoder scale.Decoder) error {
	if err := decoder.Decode(&d.PostInfo); err != nil {
		return err
	}

	return decoder.Decode(&d.Error)
}

func (d DispatchErrorWithPostInfo) Encode(encoder scale.Encoder) error {
	if err := encoder.Encode(d.PostInfo); err != nil {
		return err
	}

	return encoder.Encode(d.Error)
}

// DispatchResultWithPostInfo can be returned from dispatchable functions.
type DispatchResultWithPostInfo struct {
	IsOk bool
	Ok   PostDispatchInfo

	IsError bool
	Error   DispatchErrorWithPostInfo
}

func (d *DispatchResultWithPostInfo) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch b {
	case 0:
		d.IsOk = true

		return decoder.Decode(&d.Ok)
	case 1:
		d.IsError = true

		return decoder.Decode(&d.Error)
	}

	return nil
}

func (d DispatchResultWithPostInfo) Encode(encoder scale.Encoder) error {
	switch {
	case d.IsOk:
		if err := encoder.PushByte(0); err != nil {
			return err
		}

		return encoder.Encode(d.Ok)
	case d.IsError:
		if err := encoder.PushByte(1); err != nil {
			return err
		}

		return encoder.Encode(d.Error)
	}

	return nil
}
