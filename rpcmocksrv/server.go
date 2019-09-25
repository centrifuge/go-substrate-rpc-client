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
