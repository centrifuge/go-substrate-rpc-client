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

package config

import "os"

type Config struct {
	RPCURL string
}

// NewDefaultConfig returns a new config with default values. Default values can be overwritten with env variables,
// most importantly RPC_URL for a custom endpoint.
// TODO: rewrite as init function
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
