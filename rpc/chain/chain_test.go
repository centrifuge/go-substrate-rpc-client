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
	"github.com/centrifuge/go-substrate-rpc-client/types"
)

var chain *Chain

func TestMain(m *testing.M) {
	s := rpcmocksrv.New()
	err := s.RegisterName("chain", &mockSrv)
	if err != nil {
		panic(err)
	}

	cl, err := client.Connect(s.URL)
	// cl, err := client.Connect(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}
	chain = NewChain(cl)

	os.Exit(m.Run())
}

// MockSrv holds data and methods exposed by the RPC Mock Server used in integration tests
type MockSrv struct {
	blockHash       types.Hash
	blockHashLatest types.Hash
	header          types.Header
	signedBlock     types.SignedBlock
}

func (s *MockSrv) GetBlockHash(height *uint64) string {
	if height != nil {
		return mockSrv.blockHash.Hex()
	}
	return mockSrv.blockHashLatest.Hex()
}

func (s *MockSrv) GetBlock(hash *string) types.SignedBlock {
	return mockSrv.signedBlock
}

func (s *MockSrv) GetHeader(hash *string) types.Header {
	return mockSrv.header
}

func (s *MockSrv) GetFinalizedHead() string {
	return mockSrv.blockHashLatest.Hex()
}

// mockSrv sets default data used in tests. This data might become stale when substrate is updated – just run the tests
// against real servers and update the values stored here. To do that, replace s.URL with
// config.Default().RPCURL
var mockSrv = MockSrv{
	blockHash:       types.Hash{0xc4, 0x07, 0xff, 0x9f, 0x28, 0xda, 0x7e, 0x8c, 0xed, 0xda, 0x95, 0x61, 0x95, 0xd3, 0xe9, 0x11, 0xc8, 0x61, 0x5a, 0x2e, 0xcf, 0x0d, 0xbd, 0x6c, 0x25, 0xcf, 0x26, 0x67, 0xfb, 0x09, 0xa7, 0x2a}, //nolint:lll
	blockHashLatest: types.Hash{0xc4, 0x07, 0xff, 0x9f, 0x28, 0xda, 0x7e, 0x8c, 0xed, 0xda, 0x95, 0x61, 0x95, 0xd3, 0xe9, 0x11, 0xc8, 0x61, 0x5a, 0x2e, 0xcf, 0x0d, 0xbd, 0x6c, 0x25, 0xcf, 0x26, 0x67, 0xfb, 0x09, 0xa7, 0x2b}, //nolint:lll
	header:          types.ExamplaryHeader,
	signedBlock:     types.ExamplarySignedBlock,
}
