package config

import "os"

type Config struct {
	RPCURL string
}

// NewDefaultConfig returns a new config with default values. Default values can be overwritten with env variables,
// most importantly RPC_URL for a custom endpoint.
func NewDefaultConfig() Config {
	c := Config{}
	c.SetDefaultRPCURL()
	return c
}

// SetDefaultRPCURL reads the env variable RPC_URL and sets it in the config. If that variable is unset or empty,
// it will fallback to "http://127.0.0.1:9933"
func (c *Config) SetDefaultRPCURL() {
	if url, ok := os.LookupEnv("RPC_URL"); ok {
		c.RPCURL = url
		return
	}
	// FIXME: due to a size limit, websocket connections don't work with getMetadata right now.
	// 	Related issue: https://github.com/ethereum/go-ethereum/issues/16846
	//  Should get fixed with https://github.com/ethereum/go-ethereum/pull/19866 , released in 1.9.1
	//c.RPCURL = "ws://127.0.0.1:9944"
	c.RPCURL = "http://127.0.0.1:9933"
}
