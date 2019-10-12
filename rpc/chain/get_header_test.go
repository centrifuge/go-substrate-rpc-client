// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
// Copyright (C) 2019  Centrifuge GmbH
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

package chain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChain_GetHeaderLatest(t *testing.T) {
	header, err := chain.GetHeaderLatest()
	assert.NoError(t, err)
	assert.Equal(t, &mockSrv.header, header)
}

func TestChain_GetHeader(t *testing.T) {
	header, err := chain.GetHeader(mockSrv.blockHashLatest)
	assert.NoError(t, err)
	assert.Equal(t, &mockSrv.header, header)
}
