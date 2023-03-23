package types

import "github.com/centrifuge/go-substrate-rpc-client/v4/scale"

type VersionedMultiLocation struct {
	IsV0            bool
	MultiLocationV0 MultiLocationV0

	IsV1            bool
	MultiLocationV1 MultiLocationV1

	IsV2            bool
	MultiLocationV2 MultiLocationV2

	IsV3            bool
	MultiLocationV3 MultiLocationV3
}

func (m *VersionedMultiLocation) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch b {
	case 0:
		m.IsV0 = true

		return decoder.Decode(&m.MultiLocationV0)
	case 1:
		// V1 == V2, so default to decoding as V2:
		// https://github.com/paritytech/polkadot/blob/32b91380421d73e979d39f2325e1e0b5e23b83c7/xcm/src/lib.rs#L333
		m.IsV2 = true

		return decoder.Decode(&m.MultiLocationV2)
	case 3:
		m.IsV3 = true

		return decoder.Decode(&m.MultiLocationV3)
	}

	return nil
}

func (m VersionedMultiLocation) Encode(encoder scale.Encoder) error {
	switch {
	case m.IsV0:
		if err := encoder.PushByte(0); err != nil {
			return err
		}

		return encoder.Encode(m.MultiLocationV0)
	case m.IsV1:
		if err := encoder.PushByte(1); err != nil {
			return err
		}

		return encoder.Encode(m.MultiLocationV1)
	case m.IsV2:
		// V2 shares same index as V1:
		// https://github.com/paritytech/polkadot/blob/32b91380421d73e979d39f2325e1e0b5e23b83c7/xcm/src/lib.rs#L333
		if err := encoder.PushByte(1); err != nil {
			return err
		}

		return encoder.Encode(m.MultiLocationV2)
	case m.IsV3:
		if err := encoder.PushByte(3); err != nil {
			return err
		}

		return encoder.Encode(m.MultiLocationV3)
	}

	return nil
}
