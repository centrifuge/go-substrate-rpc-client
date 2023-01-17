// go:build ignore

package main

import (
	"log"
	"sync"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/events"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

var (
	testURLs = []string{
		"wss://fullnode.parachain.centrifuge.io",
		"wss://rpc.polkadot.io",
		"wss://statemint-rpc.polkadot.io",
		"wss://acala-rpc-0.aca-api.network",
		"wss://wss.api.moonbeam.network",
	}
)

func main() {
	var wg sync.WaitGroup

	for _, testURL := range testURLs {
		testURL := testURL

		wg.Add(1)

		go func() {
			defer wg.Done()

			api, err := gsrpc.NewSubstrateAPI(testURL)

			if err != nil {
				log.Printf("Couldn't connect to '%s': %s\n", testURL, err)
				return
			}

			meta, err := api.RPC.State.GetMetadataLatest()

			if err != nil {
				log.Printf("Couldn't retrieve metadata for '%s': %s\n", testURL, err)
				return
			}

			blockHash, err := api.RPC.Chain.GetBlockHashLatest()

			if err != nil {
				log.Printf("Couldn't get latest block hash for '%s': %s\n", testURL, err)
				return
			}

			for {
				block, err := api.RPC.Chain.GetBlock(blockHash)

				if err != nil {
					log.Printf("Couldn't get block for '%s': %s\n", testURL, err)
					return
				}

				if _, err = parseEvents(api, meta, &blockHash); err != nil {
					log.Printf("Couldn't parse events for '%s', block number %d: %s\n", testURL, block.Block.Header.Number, err)
				}

				blockHash = block.Block.Header.ParentHash
			}
		}()
	}

	wg.Wait()
}

func parseEvents(
	api *gsrpc.SubstrateAPI,
	meta *types.Metadata,
	blockHash *types.Hash,
) ([]*events.Event, error) {
	key, err := types.CreateStorageKey(meta, "System", "Events", nil)

	if err != nil {
		return nil, err
	}

	sd, err := api.RPC.State.GetStorageRaw(key, *blockHash)

	if err != nil {
		return nil, err
	}

	return events.ParseEvents(meta, sd)
}
