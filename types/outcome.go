package types

import (
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
)

type Outcome struct {
	IsComplete   bool
	IsIncomplete bool
	IsError      bool

	Weight Weight
	Error  XCMError
}

func (o *Outcome) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch b {
	case 0:
		o.IsComplete = true

		if err = decoder.Decode(&o.Weight); err != nil {
			return err
		}
	case 1:
		o.IsIncomplete = true

		if err = decoder.Decode(&o.Weight); err != nil {
			return err
		}

		if err = decoder.Decode(&o.Error); err != nil {
			return err
		}
	case 2:
		o.IsError = true

		if err = decoder.Decode(&o.Error); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown enum value for outcome - %d", b)
	}

	return nil
}

func (o Outcome) Encode(encoder scale.Encoder) error {
	switch {
	case o.IsComplete:
		if err := encoder.PushByte(0); err != nil {
			return err
		}

		return encoder.Encode(o.Weight)
	case o.IsIncomplete:
		if err := encoder.PushByte(1); err != nil {
			return err
		}

		if err := encoder.Encode(o.Weight); err != nil {
			return err
		}

		return encoder.Encode(o.Error)
	case o.IsError:
		if err := encoder.PushByte(2); err != nil {
			return err
		}

		return encoder.Encode(o.Error)
	}

	return nil
}
