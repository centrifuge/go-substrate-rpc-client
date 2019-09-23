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

package state

import (
	"os"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/client"
	"github.com/centrifuge/go-substrate-rpc-client/rpcmocksrv"
	"github.com/centrifuge/go-substrate-rpc-client/types"
)

var state *State

func TestMain(m *testing.M) {
	s := rpcmocksrv.New()
	err := s.RegisterName("state", &mockSrv)
	if err != nil {
		panic(err)
	}

	cl, err := client.Connect(s.URL)
	// cl, err := client.Connect(config.NewDefaultConfig().RPCURL)
	if err != nil {
		panic(err)
	}
	state = NewState(&cl)

	os.Exit(m.Run())
}

// MockSrv holds data and methods exposed by the RPC Mock Server used in integration tests
type MockSrv struct {
	blockHashLatest types.Hash
	metadataString  string
	metadata        types.Metadata
	runtimeVersion  types.RuntimeVersion
}

func (s *MockSrv) GetMetadata(hash *string) string {
	return mockSrv.metadataString
}

func (s *MockSrv) GetRuntimeVersion(hash *string) types.RuntimeVersion {
	return mockSrv.runtimeVersion
}

// mockSrv sets default data used in tests. This data might become stale when substrate is updated – just run the tests
// against real servers and update the values stored here. To do that, replace s.URL with
// config.NewDefaultConfig().RPCURL
var mockSrv = MockSrv{
	blockHashLatest: types.Hash{1, 2, 3},
	metadata:        types.ExamplaryMetadataV4,
	metadataString:  types.ExamplaryMetadataV4String,
	runtimeVersion:  types.RuntimeVersion{APIs: []types.RuntimeVersionAPI{{APIID: "0xdf6acb689907609b", Version: 0x2}, {APIID: "0x37e397fc7c91f5e4", Version: 0x1}, {APIID: "0x40fe3ad401f8959a", Version: 0x3}, {APIID: "0xd2bc9897eed08f15", Version: 0x1}, {APIID: "0xf78b278be53f454c", Version: 0x1}, {APIID: "0xed99c5acb25eedf5", Version: 0x2}, {APIID: "0xdd718d5cc53262d4", Version: 0x1}, {APIID: "0x7801759919ee83e5", Version: 0x1}}, AuthoringVersion: 0xa, ImplName: "substrate-node", ImplVersion: 0x3e, SpecName: "node", SpecVersion: 0x3c}, //nolint:lll
}
