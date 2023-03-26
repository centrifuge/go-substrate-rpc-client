package types

import (
	"errors"
	"fmt"
	"math"
	"math/bits"

	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
)

type BitOrder uint

const (
	BitOrderLsb0 BitOrder = iota
	BitOrderMsb0
)

var (
	BitOrderName = map[BitOrder]string{
		BitOrderLsb0: "Lsb0",
		BitOrderMsb0: "Msb0",
	}

	BitOrderValue = map[string]BitOrder{
		"Lsb0": BitOrderLsb0,
		"Msb0": BitOrderMsb0,
	}
)

func (b *BitOrder) String() string {
	return BitOrderName[*b]
}

func NewBitOrderFromString(s string) (BitOrder, error) {
	bitOrder, ok := BitOrderValue[s]

	if !ok {
		return 0, fmt.Errorf("bit order '%s' not supported", s)
	}

	return bitOrder, nil
}

type BitVec struct {
	BitOrder BitOrder

	NumBytes uint
	Bits     uint
}

func NewBitVec(bitOrder BitOrder) *BitVec {
	return &BitVec{
		BitOrder: bitOrder,
	}
}

func (b *BitVec) Decode(decoder scale.Decoder) error {
	if err := b.GetMinimumNumberOfBytes(decoder); err != nil {
		return err
	}

	var total uint

	for i := uint(0); i < b.NumBytes; i++ {
		total = total << 8

		cb, err := decoder.ReadOneByte()

		if err != nil {
			return err
		}

		if b.BitOrder == BitOrderLsb0 {
			cb = bits.Reverse8(cb)
		}

		total = total | uint(cb)
	}

	b.Bits = total

	return nil
}

func (b *BitVec) GetMinimumNumberOfBytes(decoder scale.Decoder) error {
	nb, err := decoder.DecodeUintCompact()

	if err != nil {
		return err
	}

	numberOfBits := nb.Uint64()

	if numberOfBits == 0 {
		return errors.New("invalid number of bits")
	}

	b.NumBytes = uint(math.Ceil(float64(numberOfBits) / 8))

	return nil
}

func (b *BitVec) String() string {
	fmtArgs := b.NumBytes * 8

	return fmt.Sprintf("%#0*b", fmtArgs, b.Bits)
}
