// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
// Copyright (C) 2019  Philip Stanislaus, Philip Stehlik, Vimukthi Wickramasinghe
//
// This file is part of Go Substrate RPC Client (GSRPC).
//
// GSRPC is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// GSRPC is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package types_test

import (
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/scale"
	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/stretchr/testify/assert"
)

type PhaseEnum struct {
	IsApplyExtrinsic bool
	AsApplyExtrinsic uint32
	IsFinalization   bool
}

func (m *PhaseEnum) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()

	if err != nil {
		return err
	}

	if b == 0 {
		m.IsApplyExtrinsic = true
		err = decoder.Decode(&m.AsApplyExtrinsic)
	} else if b == 1 {
		m.IsFinalization = true
	}

	if err != nil {
		return err
	}

	return nil
}

func (m PhaseEnum) Encode(encoder scale.Encoder) error {
	var err1, err2 error
	if m.IsApplyExtrinsic {
		err1 = encoder.PushByte(0)
		err2 = encoder.Encode(m.AsApplyExtrinsic)
	} else if m.IsFinalization {
		err1 = encoder.PushByte(1)
	}

	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}

	return nil
}

func TestPhaseEnumApplyExtrinsic(t *testing.T) {
	applyExtrinsic := PhaseEnum{
		IsApplyExtrinsic: true,
		AsApplyExtrinsic: 1234,
	}

	enc, err := types.EncodeToHexString(applyExtrinsic)
	assert.NoError(t, err)

	var dec PhaseEnum
	err = types.DecodeFromHexString(enc, &dec)
	assert.NoError(t, err)

	assert.Equal(t, applyExtrinsic, dec)
}

func TestPhaseEnumFinalization(t *testing.T) {
	finalization := PhaseEnum{
		IsFinalization: true,
	}

	enc, err := types.EncodeToHexString(finalization)
	assert.NoError(t, err)

	var dec PhaseEnum
	err = types.DecodeFromHexString(enc, &dec)
	assert.NoError(t, err)

	assert.Equal(t, finalization, dec)
}
