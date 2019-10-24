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

package state

import (
	"context"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/ethereum/go-ethereum/rpc"
)

// SubscribeStorageRaw subscribes the storage for the given keys, returning a subscription and a channel that will
// receive server notifications containing the storage change sets.
//
// Slow subscribers will be dropped eventually. Client buffers up to 20000 notifications before considering the
// subscriber dead. The subscription Err channel will receive ErrSubscriptionQueueOverflow. Use a sufficiently
// large buffer on the channel or ensure that the channel usually has at least one reader to prevent this issue.
func (s *State) SubscribeStorageRaw(keys []types.StorageKey) (
	*rpc.ClientSubscription, <-chan types.StorageChangeSet, error) {
	ctx, _ := context.WithTimeout(context.Background(), 1000*time.Second)

	c := make(chan types.StorageChangeSet)

	keyss := make([]string, len(keys))
	for i := range keys {
		keyss[i] = keys[i].Hex()
	}

	sub, err := (*s.client).Subscribe(ctx, "state", c, keyss)
	if err != nil {
		return nil, nil, err
	}

	return sub, c, nil
}
