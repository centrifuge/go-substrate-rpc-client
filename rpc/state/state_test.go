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

package state

import (
	"os"
	"strings"
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
	// cl, err := client.Connect("ws://35.246.140.178:9944")
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
	metadata        *types.Metadata
	runtimeVersion  types.RuntimeVersion
	storageKeyHex   string
	storageDataHex  string
	storageSize     types.U64
	storageHashHex  string
}

func (s *MockSrv) GetMetadata(hash *string) string {
	return mockSrv.metadataString
}

func (s *MockSrv) GetRuntimeVersion(hash *string) types.RuntimeVersion {
	return mockSrv.runtimeVersion
}

func (s *MockSrv) GetKeys(key string, hash *string) []string {
	if !strings.HasPrefix(mockSrv.storageKeyHex, key) {
		panic("key not found")
	}
	return []string{mockSrv.storageKeyHex}
}

func (s *MockSrv) GetStorage(key string, hash *string) string {
	if key != s.storageKeyHex {
		panic("key not found")
	}
	return mockSrv.storageDataHex
}

func (s *MockSrv) GetStorageSize(key string, hash *string) types.U64 {
	return mockSrv.storageSize
}

func (s *MockSrv) GetStorageHash(key string, hash *string) string {
	return mockSrv.storageHashHex
}

// mockSrv sets default data used in tests. This data might become stale when substrate is updated – just run the tests
// against real servers and update the values stored here. To do that, replace s.URL with
// config.NewDefaultConfig().RPCURL
var mockSrv = MockSrv{
	blockHashLatest: types.Hash{1, 2, 3},
	metadata:        types.ExamplaryMetadataV4,
	metadataString:  types.ExamplaryMetadataV4String,
	runtimeVersion:  types.RuntimeVersion{APIs: []types.RuntimeVersionAPI{{APIID: "0xdf6acb689907609b", Version: 0x2}, {APIID: "0x37e397fc7c91f5e4", Version: 0x1}, {APIID: "0x40fe3ad401f8959a", Version: 0x3}, {APIID: "0xd2bc9897eed08f15", Version: 0x1}, {APIID: "0xf78b278be53f454c", Version: 0x1}, {APIID: "0xed99c5acb25eedf5", Version: 0x2}, {APIID: "0xdd718d5cc53262d4", Version: 0x1}, {APIID: "0x7801759919ee83e5", Version: 0x1}}, AuthoringVersion: 0xa, ImplName: "substrate-node", ImplVersion: 0x3e, SpecName: "node", SpecVersion: 0x3c}, //nolint:lll
	storageKeyHex:   "0x0e4944cfd98d6f4cc374d16f5a4e3f9c",
	storageDataHex:  "0xb82d895d00000000",
	storageSize:     926778,
	storageHashHex:  "0xdf0e877ee1cb973b9a566f53707d365b269d7131b55e65b9790994e4e63b95e1",
}
