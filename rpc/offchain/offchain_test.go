package offchain

import (
	"os"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v2/client"
	"github.com/centrifuge/go-substrate-rpc-client/v2/rpcmocksrv"
)

var offchain *Offchain

func TestMain(m *testing.M) {
	s := rpcmocksrv.New()
	err := s.RegisterName("offchain", &mockSrv)
	if err != nil {
		panic(err)
	}

	cl, err := client.Connect(s.URL)
	// cl, err := client.Connect(config.Default().RPCURL)
	if err != nil {
		panic(err)
	}
	offchain = NewOffchain(cl)

	os.Exit(m.Run())
}

// MockSrv holds data and methods exposed by the RPC Mock Server used in integration tests
type MockSrv struct {
	storageKeyHex   string
	storageValueHex string
}

func (s *MockSrv) LocalStorageGet(kind string, key string) string {
	if key != s.storageKeyHex {
		return ""
	}
	return mockSrv.storageValueHex
}

var mockSrv = MockSrv{
	storageKeyHex:   "0x666f6f",
	storageValueHex: "0xdeadbeef", //nolint:lll
}
