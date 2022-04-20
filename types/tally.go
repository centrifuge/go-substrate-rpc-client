package types

import "github.com/centrifuge/go-substrate-rpc-client/v4/scale"

// Tally is a wrapper struct for the Referenda Tally:
//
// type Tally: VoteTally<Self::Votes> + Default + Clone + Codec + Eq + Debug + TypeInfo;
type Tally struct {
	Votes U128
	Total U128
}

func (t *Tally) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&t.Votes)
	if err != nil {
		return err
	}

	err = decoder.Decode(&t.Total)
	if err != nil {
		return err
	}

	return nil
}

func (t Tally) Encode(encoder scale.Encoder) error {
	err := encoder.Encode(t.Votes)
	if err != nil {
		return err
	}

	err = encoder.Encode(t.Total)
	if err != nil {
		return err
	}

	return nil
}
