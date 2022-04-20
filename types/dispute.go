package types

import "github.com/centrifuge/go-substrate-rpc-client/v4/scale"

// /// Whether the dispute is local or remote.
//#[derive(Encode, Decode, Clone, PartialEq, Eq, RuntimeDebug, TypeInfo)]
//pub enum DisputeLocation {
//	Local,
//	Remote,
//}
//
///// The result of a dispute, whether the candidate is deemed valid (for) or invalid (against).
//#[derive(Encode, Decode, Clone, PartialEq, Eq, RuntimeDebug, TypeInfo)]
//pub enum DisputeResult {
//	Valid,
//	Invalid,
//}

type DisputeLocation struct {
	IsLocal bool

	IsRemote bool
}

func (d *DisputeLocation) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()

	if err != nil {
		return err
	}

	switch b {
	case 0:
		d.IsLocal = true
	case 1:
		d.IsRemote = true
	}

	return nil
}

func (d DisputeLocation) Encode(encoder scale.Encoder) error {
	switch {
	case d.IsLocal:
		return encoder.PushByte(0)
	case d.IsRemote:
		return encoder.PushByte(1)
	}

	return nil
}

type DisputeResult struct {
	IsValid bool

	IsInvalid bool
}

func (d *DisputeResult) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()

	if err != nil {
		return err
	}

	switch b {
	case 0:
		d.IsValid = true
	case 1:
		d.IsInvalid = true
	}

	return nil
}

func (d DisputeResult) Encode(encoder scale.Encoder) error {
	switch {
	case d.IsValid:
		return encoder.PushByte(0)
	case d.IsInvalid:
		return encoder.PushByte(1)
	}

	return nil
}
