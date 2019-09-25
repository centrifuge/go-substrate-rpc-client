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
