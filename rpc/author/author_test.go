// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based author RPC calls
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

package author

import (
	"os"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/client"
	"github.com/centrifuge/go-substrate-rpc-client/rpcmocksrv"
)

var author *Author

func TestMain(m *testing.M) {
	s := rpcmocksrv.New()
	err := s.RegisterName("author", &mockSrv)
	if err != nil {
		panic(err)
	}

	cl, err := client.Connect(s.URL)
	// cl, err := client.Connect(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}
	author = NewAuthor(cl)

	os.Exit(m.Run())
}

// MockSrv holds data and methods exposed by the RPC Mock Server used in integration tests
type MockSrv struct {
	submitExtrinsicHash string
	pendingExtrinsics   []string
}

func (s *MockSrv) SubmitExtrinsic(extrinsic string) string {
	return mockSrv.submitExtrinsicHash
}

func (s *MockSrv) PendingExtrinsics() []string {
	return mockSrv.pendingExtrinsics
}

// mockSrv sets default data used in tests. This data might become stale when substrate is updated – just run the tests
// against real servers and update the values stored here. To do that, replace s.URL with
// config.Default().RPCURL
var mockSrv = MockSrv{
	submitExtrinsicHash: "0x9a8ef9794ded03b4d1ae45034351210e87f970f1f9500994bca82f9cd5a1166e",
	pendingExtrinsics:   []string{"0x290284ffd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d00a023bbe883405b5fac2aa114093fcf3d0802d2f3d3715e09129b00a4bf741048caf53d8c7d97e872caa703e7d04f174a4e2ed4acadee4173a8b6bab7e45c0a06000c000600ff8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48e56c"}, //nolint:lll
}
