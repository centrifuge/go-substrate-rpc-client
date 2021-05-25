package types

import (
	"fmt"

	"github.com/Phala-Network/go-substrate-rpc-client/v3/scale"
)

type ElectionCompute byte

const (
	// Result was forcefully computed on chain at the end of the session.
	OnChain ElectionCompute = 0
	// Result was submitted and accepted to the chain via a signed transaction.
	Signed ElectionCompute = 1
	// Result was submitted and accepted to the chain via an unsigned transaction (by an authority).
	Unsigned ElectionCompute = 2
)

func (ec *ElectionCompute) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	vb := ElectionCompute(b)
	switch vb {
	case OnChain, Signed, Unsigned:
		*ec = vb
	default:
		return fmt.Errorf("unknown ElectionCompute enum: %v", vb)
	}
	return err
}

func (ec ElectionCompute) Encode(encoder scale.Encoder) error {
	return encoder.PushByte(byte(ec))
}

type OptionElectionCompute struct {
	option
	value ElectionCompute
}

func NewOptionElectionCompute(value ElectionCompute) OptionElectionCompute {
	return OptionElectionCompute{option{true}, value}
}

func NewOptionElectionComputeEmpty() OptionElectionCompute {
	return OptionElectionCompute{option: option{false}}
}

func (o OptionElectionCompute) Encode(encoder scale.Encoder) error {
	return encoder.EncodeOption(o.hasValue, o.value)
}

func (o *OptionElectionCompute) Decode(decoder scale.Decoder) error {
	return decoder.DecodeOption(&o.hasValue, &o.value)
}

func (o *OptionElectionCompute) SetSome(value ElectionCompute) {
	o.hasValue = true
	o.value = value
}

func (o *OptionElectionCompute) SetNone() {
	o.hasValue = false
	o.value = 0
}

func (o OptionElectionCompute) Unwrap() (ok bool, value ElectionCompute) {
	return o.hasValue, o.value
}
