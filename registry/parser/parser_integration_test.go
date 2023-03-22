//go:build integration

package parser

import (
	"log"
	"sync"
	"testing"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/registry"
	"github.com/centrifuge/go-substrate-rpc-client/v4/registry/state"
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

func TestParser_GetEvents(t *testing.T) {
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

			parser, err := NewDefaultEventParser(state.NewProvider(api.RPC.State), registry.NewFactory())

			if err != nil {
				log.Printf("Couldn't create eventParser: %s", err)
				return
			}

			blockHash, err := api.RPC.Chain.GetBlockHashLatest()

			if err != nil {
				log.Printf("Couldn't get latest block hash for '%s': %s\n", testURL, err)
				return
			}

			var previousBlock *types.SignedBlock

			processedBlockCount := 0

			for {
				block, err := api.RPC.Chain.GetBlock(blockHash)

				if err != nil {
					log.Printf("Skipping block %d for '%s' due to a block retrieval error\n", previousBlock.Block.Header.Number-1, testURL)

					if blockHash, err = api.RPC.Chain.GetBlockHash(uint64(previousBlock.Block.Header.Number - 2)); err != nil {
						log.Printf("Couldn't get block hash for block %d: %s", previousBlock.Block.Header.Number-2, err)

						return
					}

					continue
				}

				previousBlock = block

				blockHash = block.Block.Header.ParentHash

				_, err = parser.GetEvents(blockHash)

				if err != nil {
					log.Printf("Couldn't parse events for '%s', block number %d: %s\n", testURL, block.Block.Header.Number, err)
					return
				}

				processedBlockCount++

				if processedBlockCount%500 == 0 {
					log.Printf("Parsed events for %d blocks for '%s' so far, last block number %d\n", processedBlockCount, testURL, block.Block.Header.Number)
				}
			}
		}()
	}

	wg.Wait()
}
