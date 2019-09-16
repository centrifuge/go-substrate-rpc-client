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

package state

import (
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/rpc/chain"
	"github.com/stretchr/testify/assert"
)

func TestState_GetMetadataLatest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping end-to-end test in short mode.")
	}

	metadata, err := state.GetMetadataLatest()
	assert.NoError(t, err)
	assert.Equal(t, "system", metadata.Metadata.Modules[0].Name)
}

func TestState_GetMetadata(t *testing.T) {
	chain := chain.NewChain(state.client)

	hash, err := chain.GetBlockHashLatest()
	assert.NoError(t, err)

	metadata, err := state.GetMetadataLatest(hash)
	assert.NoError(t, err)
	assert.Equal(t, "system", metadata.Metadata.Modules[0].Name)
}
