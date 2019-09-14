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

package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var exampleRuntimeVersion = RuntimeVersion{
	APIs:             []RuntimeVersionAPI{exampleRuntimeVersionAPI},
	AuthoringVersion: 13,
	ImplName:         "My impl",
	ImplVersion:      21,
	SpecName:         "My spec",
	SpecVersion:      39,
}

var exampleRuntimeVersionAPI = RuntimeVersionAPI{
	APIID:   "0x37e397fc7c91f5e4",
	Version: 23,
}

func TestRuntimeVersion_Encode_Decode(t *testing.T) {
	enc, err := EncodeToBytes(exampleRuntimeVersion)
	assert.NoError(t, err)

	var output RuntimeVersion
	err = DecodeFromBytes(enc, &output)
	assert.NoError(t, err)

	assert.Equal(t, exampleRuntimeVersion, output)
}

func TestRuntimeVersionAPI_Encode_Decode(t *testing.T) {
	enc, err := EncodeToBytes(exampleRuntimeVersionAPI)
	assert.NoError(t, err)

	var output RuntimeVersionAPI
	err = DecodeFromBytes(enc, &output)
	assert.NoError(t, err)

	assert.Equal(t, exampleRuntimeVersionAPI, output)
}
