package types

import "github.com/centrifuge/go-substrate-rpc-client/v4/scale"

// LotteryCallIndex is a 16 bit wrapper around the following Lottery CallIndex:
//
//   Any runtime call can be encoded into two bytes which represent the pallet and call index.
//   We use this to uniquely match someone's incoming call with the calls configured for the lottery.
//   type CallIndex = (u8, u8);
type LotteryCallIndex struct {
	PalletIndex uint8
	CallIndex   uint8
}

func (m *LotteryCallIndex) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&m.PalletIndex)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.CallIndex)
	if err != nil {
		return err
	}

	return nil
}

func (m LotteryCallIndex) Encode(encoder scale.Encoder) error {
	err := encoder.Encode(m.PalletIndex)
	if err != nil {
		return err
	}

	err = encoder.Encode(m.CallIndex)
	if err != nil {
		return err
	}

	return nil
}
