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

	"github.com/stretchr/testify/assert"
)

func TestState_GetRuntimeVersionLatest(t *testing.T) {
	runtimeVersion, err := state.GetRuntimeVersionLatest()
	assert.NoError(t, err)
	assert.Equal(t, "substrate-node", runtimeVersion.ImplName)
}

// TODO make test dynamic
//func TestState_GetRuntimeVersion(t *testing.T) {
//	bz, _ := hex.DecodeString("cc9ea640d4d4f4dd260b1cbb65cb275df995b056710265f9becdd7e6e1a7b9e0")
//
//	var bz32 [32]byte
//	copy(bz32[:], bz)
//
//	runtimeVersion, err := state.GetRuntimeVersion(types.NewHash(bz32))
//	assert.NoError(t, err)
//	assert.Equal(t, "system", runtimeVersion.RuntimeVersion.Modules[0].Name)
//}
