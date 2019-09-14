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
	"github.com/centrifuge/go-substrate-rpc-client/types"
)

func (s *State) GetRuntimeVersion(blockHash types.Hash) (*types.RuntimeVersion, error) {
	return s.getRuntimeVersion(&blockHash)
}

func (s *State) GetRuntimeVersionLatest() (*types.RuntimeVersion, error) {
	return s.getRuntimeVersion(nil)
}

func (s *State) getRuntimeVersion(blockHash *types.Hash) (*types.RuntimeVersion, error) {
	var runtimeVersion types.RuntimeVersion
	var err error
	if blockHash == nil {
		err = (*s.client).Call(&runtimeVersion, "state_getRuntimeVersion")
	} else {
		err = (*s.client).Call(&runtimeVersion, "state_getRuntimeVersion", *blockHash)
	}
	if err != nil {
		return &runtimeVersion, err
	}
	return &runtimeVersion, err
}

//func (s *State) getRuntimeVersion(blockHash *types.Hash) (*types.RuntimeVersion, error) {
//	runtimeVersion := types.NewRuntimeVersion()
//
//	var res string
//	var err error
//	if blockHash == nil {
//		err = (*s.client).Call(&res, "state_getRuntimeVersion")
//	} else {
//		err = (*s.client).Call(&res, "state_getRuntimeVersion", *blockHash)
//	}
//	if err != nil {
//		return runtimeVersion, err
//	}
//
//	err = types.DecodeFromHexString(res, runtimeVersion)
//	return runtimeVersion, err
//}
