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
	"os"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/client"
	"github.com/centrifuge/go-substrate-rpc-client/rpcmocksrv"
)

var chain *Chain

func TestMain(m *testing.M) {
	s := rpcmocksrv.New()
	err := s.RegisterName("chain", &mockSrv)
	if err != nil {
		panic(err)
	}

	cl, err := client.Connect(s.URL)
	if err != nil {
		panic(err)
	}
	chain = NewChain(&cl)

	os.Exit(m.Run())
}

// MockSrv holds data and methods exposed by the RPC Mock Server used in integration tests
type MockSrv struct {
	blockHash       string
	blockHashLatest string
}

func (s *MockSrv) GetBlockHash(height *uint64) string {
	if height != nil {
		return mockSrv.blockHash
	}
	return mockSrv.blockHashLatest
}

// mockSrv sets default data used in tests. This data might become stale when substrate is updated – just run the tests
// against real servers and update the values stored here. To do that, replace s.URL with
// config.NewDefaultConfig().RPCURL
var mockSrv = MockSrv{
	blockHash:       "0xc407ff9f28da7e8cedda956195d3e911c8615a2ecf0dbd6c25cf2667fb09a72a",
	blockHashLatest: "0xc407ff9f28da7e8cedda956195d3e911c8615a2ecf0dbd6c25cf2667fb09a72b",
}
