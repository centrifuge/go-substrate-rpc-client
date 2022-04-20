package types

import "github.com/centrifuge/go-substrate-rpc-client/v4/scale"

type HRMPChannelID struct {
	Sender    U32
	Recipient U32
}

func (h *HRMPChannelID) Decode(decoder scale.Decoder) error {
	if err := decoder.Decode(&h.Sender); err != nil {
		return err
	}

	return decoder.Decode(&h.Recipient)
}

func (h HRMPChannelID) Encode(encoder scale.Encoder) error {
	if err := encoder.Encode(h.Sender); err != nil {
		return err
	}

	return encoder.Encode(h.Recipient)
}
