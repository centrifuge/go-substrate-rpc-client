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

package types_test

import (
	"testing"

	. "github.com/centrifuge/go-substrate-rpc-client/v3/types"
	"github.com/stretchr/testify/assert"
)

func TestExtrinsicEra_Immortal(t *testing.T) {
	var e ExtrinsicEra
	err := DecodeFromHexString("0x00", &e)
	assert.NoError(t, err)
	assert.Equal(t, ExtrinsicEra{IsImmortalEra: true}, e)
}

func TestExtrinsicEra_Mortal(t *testing.T) {
	var e ExtrinsicEra
	err := DecodeFromHexString("0x4e9c", &e)
	assert.NoError(t, err)
	assert.Equal(t, ExtrinsicEra{
		IsMortalEra: true, AsMortalEra: MortalEra{78, 156},
	}, e)
}

func TestExtrinsicEra_EncodeDecode(t *testing.T) {
	var e ExtrinsicEra
	err := DecodeFromHexString("0x4e9c", &e)
	assert.NoError(t, err)
	assertRoundtrip(t, e)
}

func TestExtrinsicEra_MortalEraEncoding(t *testing.T) {
	testCases := []struct {
		Current        uint64
		ValidityPeriod uint64
		Era            MortalEra
	}{
		{
			Current:        2251516,
			ValidityPeriod: 64,
			Era:            MortalEra{First: 0xc5, Second: 0x03},
		},
	}

	for _, testCase := range testCases {
		period, phase := Mortal(testCase.ValidityPeriod, testCase.Current)

		era := NewMortalEra(period, phase)
		assert.Equal(t, testCase.Era.First, era.First)
		assert.Equal(t, testCase.Era.Second, era.Second)
	}
}

func TestMortal(t *testing.T) {
	period, phase := Mortal(200, 1400)
	assert.Equal(t, uint64(120), phase)
	assert.Equal(t, uint64(256), period)
}

func TestNewMortalEra(t *testing.T) {
	testCases := []struct {
		Period uint64
		Phase  uint64
		Era    MortalEra
	}{
		{
			Period: 128,
			Phase:  125,
			Era:    MortalEra{First: 0xd6, Second: 0x07},
		},
		{
			Period: 4096,
			Phase:  3641,
			Era:    MortalEra{First: 0x9b, Second: 0xe3},
		},
		{
			Period: 64,
			Phase:  38,
			Era:    MortalEra{First: 0x65, Second: 0x02},
		},
		{
			Period: 64,
			Phase:  40,
			Era:    MortalEra{First: 0x85, Second: 0x02},
		},
		{
			Period: 32768,
			Phase:  20000,
			Era:    MortalEra{First: 78, Second: 156},
		},
	}

	for _, testCase := range testCases {
		era := NewMortalEra(testCase.Period, testCase.Phase)
		assert.Equal(t, testCase.Era.First, era.First)
		assert.Equal(t, testCase.Era.Second, era.Second)
	}
}
