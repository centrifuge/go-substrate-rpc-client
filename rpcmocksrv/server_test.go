package rpcmocksrv

import (
	"testing"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
)

type TestService struct {
}

func (ts *TestService) Ping(s string) string {
	return s
}

func TestServer(t *testing.T) {
	s := New()

	ts := new(TestService)
	err := s.RegisterName("testserv3", ts)
	assert.NoError(t, err)

	c, err := rpc.Dial(s.URL)
	assert.NoError(t, err)

	var res string
	err = c.Call(&res, "testserv3_ping", "hello")
	assert.NoError(t, err)

	assert.Equal(t, "hello", res)
}
