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

package client

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/rpc"
)

type Client interface {
	Call(result interface{}, method string, args ...interface{}) error

	Subscribe(ctx context.Context, namespace string, channel interface{}, args ...interface{}) (
		*rpc.ClientSubscription, error)

	URL() string
	//MetaData(cache bool) (*MetadataVersioned, error)
}

type client struct {
	*rpc.Client

	url string

	// metadataVersioned is the metadata cache to prevent unnecessary requests
	//metadataVersioned *MetadataVersioned

	//metadataLock sync.RWMutex
}

// Returns the URL the client connects to
func (c client) URL() string {
	return c.url
}

// TODO move to State struct
//func (c *client) MetaData(cache bool) (m *MetadataVersioned, err error) {
//	if cache && c.metadataVersioned != nil {
//		c.metadataLock.RLock()
//		defer c.metadataLock.RUnlock()
//		m = c.metadataVersioned
//	} else {
//		s := NewStateRPC(c)
//		m, err = s.MetaData(nil)
//		if err != nil {
//			return nil, err
//		}
//		// set cache
//		c.metadataLock.Lock()
//		defer c.metadataLock.Unlock()
//		c.metadataVersioned = m
//	}
//	return
//}

// Connect connects to the provided url
func Connect(url string) (Client, error) {
	log.Printf("Connecting to %v...", url)
	c, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	cc := client{c, url}
	return &cc, nil
}
