package rpcmocksrv

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
)

type Server struct {
	*rpc.Server
	// Host consists of hostname and port
	Host string
	// URL consists of protocol, hostname and port
	URL string
}

// New creates a new RPC mock server with a random port that allows registration of services
func New() *Server {
	port := randomPort()
	host := "localhost:" + strconv.Itoa(port)

	_, rpcServ, err := rpc.StartWSEndpoint(host, []rpc.API{}, []string{}, []string{"*"}, true)
	if err != nil {
		panic(err)
	}
	s := Server{
		Server: rpcServ,
		Host:   host,
		URL:    "ws://" + host,
	}
	return &s
}

func randomPort() int {
	rand.Seed(time.Now().UnixNano())
	min := 10000
	max := 30000
	return rand.Intn(max-min+1) + min
}
