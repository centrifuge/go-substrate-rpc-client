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
