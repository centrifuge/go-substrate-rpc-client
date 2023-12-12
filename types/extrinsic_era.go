// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
//
// Copyright 2019 Centrifuge GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	"errors"
	"math"
	"strconv"

	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
)

// ExtrinsicEra indicates either a mortal or immortal extrinsic
type ExtrinsicEra struct {
	IsImmortalEra bool
	// AsImmortalEra ImmortalEra
	IsMortalEra bool
	AsMortalEra MortalEra
}

func (e *ExtrinsicEra) Decode(decoder scale.Decoder) error {
	first, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	if first == 0 {
		e.IsImmortalEra = true
		return nil
	}

	second, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	encoded := uint64(first) + (uint64(second) << 8)
	period := 2 << (encoded % (1 << 4))
	quantizeFactor := period >> 12

	if quantizeFactor <= 1 {
		quantizeFactor = 1
	}

	phase := (encoded >> 4) * uint64(quantizeFactor)

	if period >= 4 && phase < uint64(period) {
		e.IsMortalEra = true
		e.AsMortalEra = MortalEra{
			First:  U64(period),
			Second: U64(phase),
		}
		return nil
	}

	return errors.New("invalid era")
}

func (e ExtrinsicEra) Encode(encoder scale.Encoder) error {
	if e.IsImmortalEra {
		return encoder.PushByte(0)
	}

	// let quantize_factor = (*period as u64 >> 12).max(1);
	quantizeFactor := e.AsMortalEra.First >> 12

	if quantizeFactor <= 1 {
		quantizeFactor = 1
	}

	// let encoded = (period.trailing_zeros() - 1).max(1).min(15) as u16 |
	// ((phase / quantize_factor) << 4) as u16;
	trailingZeroes := getTrailingZeroes(e.AsMortalEra.First) - 1

	if trailingZeroes <= 1 {
		trailingZeroes = 1
	}

	if trailingZeroes >= 15 {
		trailingZeroes = 15
	}

	r := U16((e.AsMortalEra.Second / quantizeFactor) << 4)

	encoded := trailingZeroes | r

	// encoded.encode_to(output);
	return encoder.Encode(encoded)
}

// MortalEra for an extrinsic, indicating period and phase
type MortalEra struct {
	First  U64
	Second U64
}

func NewMortalEra(currentBlock BlockNumber, blockHashCount U64) MortalEra {
	// BlockHashCount::get().checked_next_power_of_two().map(|c| c / 2).unwrap_or(2) as u64;
	np := getNextPowerOfTwo(blockHashCount, 2)

	var npb U64

	if np > 2 {
		npb = np / 2
	}

	// let period = period.checked_next_power_of_two().unwrap_or(1 << 16).max(4).min(1 << 16);
	period := getNextPowerOfTwo(npb, 1<<16)

	if period <= 4 {
		period = 4
	}

	if period >= 1<<16 {
		period = 1 << 16
	}

	// let phase = current % period;
	phase := U64(currentBlock) % period

	// let quantize_factor = (period >> 12).max(1);
	quantizeFactor := period >> 12

	if quantizeFactor <= 1 {
		quantizeFactor = 1
	}

	// let quantized_phase = phase / quantize_factor * quantize_factor;
	quantizedPhase := phase / quantizeFactor * quantizeFactor

	return MortalEra{period, quantizedPhase}
}

func getNextPowerOfTwo(n U64, def U64) U64 {
	bn := strconv.FormatInt(int64(n), 2)
	numBits := len(bn)

	if (1 << (numBits - 1)) == n {
		return n
	}

	res := uint(1 << numBits)

	if res > math.MaxUint64 {
		return def
	}

	return U64(res)
}

func getTrailingZeroes(n U64) U16 {
	var count U16

	for n > 0 {
		if n%2 == 1 {
			break
		}

		n = n / 2
		count++
	}

	return count
}
