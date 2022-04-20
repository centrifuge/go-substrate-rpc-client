package types

import "github.com/centrifuge/go-substrate-rpc-client/v4/scale"

// PostDispatchInfo is used in DispatchResultWithPostInfo.
// Weight information that is only available post dispatch.
type PostDispatchInfo struct {
	ActualWeight OptionWeight
	PaysFee      Pays
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
	Ok    *PostDispatchInfo
	Error *DispatchErrorWithPostInfo
}

func (d *DispatchResultWithPostInfo) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch b {
	case 0:
		err = decoder.Decode(&d.Ok)
		return nil
	default:
		derr := DispatchErrorWithPostInfo{}
		err = decoder.Decode(&derr)
		if err != nil {
			return err
		}
		d.Error = &derr
		return nil
	}
}

func (d DispatchResultWithPostInfo) Encode(encoder scale.Encoder) error {
	if d.Ok != nil {
		return encoder.Encode(d.Ok)
	}

	return d.Error.Encode(encoder)
}
