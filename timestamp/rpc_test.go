package timestamp

import (
	"strconv"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client"
	"github.com/centrifuge/go-substrate-rpc-client/testrpc"
	"github.com/stretchr/testify/assert"
)

var testServer *testrpc.Server
var testClient substrate.Client
var rpcPort int

func TestMain(m *testing.M) {
	testServer = new(testrpc.Server)
	var err error
	rpcPort, err = testServer.Init()
	if err != nil {
		panic(err)
	}

	testClient, err = substrate.Connect("ws://localhost:" + strconv.Itoa(rpcPort))
	if err != nil {
		panic(err)
	}
	m.Run()
}

func TestNow(t *testing.T) {
	testServer.AddStorageKey("0x0e4944cfd98d6f4cc374d16f5a4e3f9c", "0x7ab3425d00000000")
	ts, err := Now(testClient)
	assert.NoError(t, err)
	assert.Equal(t, "2019-08-01 11:40:10 +0200 CEST", ts.String())
}
