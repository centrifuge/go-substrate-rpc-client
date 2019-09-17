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

package config

import "os"

type Config struct {
	RPCURL string
}

// NewDefaultConfig returns a new config with default values. Default values can be overwritten with env variables,
// most importantly RPC_URL for a custom endpoint.
func NewDefaultConfig() Config {
	c := Config{}
	c.ExtractDefaultRPCURL()
	return c
}

// ExtractDefaultRPCURL reads the env variable RPC_URL and sets it in the config. If that variable is unset or empty,
// it will fallback to "http://127.0.0.1:9933"
func (c *Config) ExtractDefaultRPCURL() {
	if url, ok := os.LookupEnv("RPC_URL"); ok {
		c.RPCURL = url
		return
	}

	// Fallback
	c.RPCURL = "ws://127.0.0.1:9944"
}
