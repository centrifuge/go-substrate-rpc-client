//go:build live

package retriever

import (
	"log"
	"sync"
	"testing"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/registry/state"
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

func TestLive_EventRetriever_GetEvents(t *testing.T) {
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

			retriever, err := NewDefaultEventRetriever(state.NewProvider(api.RPC.State))

			if err != nil {
				log.Printf("Couldn't create event retriever: %s", err)
				return
			}

			header, err := api.RPC.Chain.GetHeaderLatest()

			if err != nil {
				log.Printf("Couldn't get latest header for '%s': %s\n", testURL, err)
				return
			}

			processedBlockCount := 0

			for {
				blockHash, err := api.RPC.Chain.GetBlockHash(uint64(header.Number))

				if err != nil {
					log.Printf("Couldn't retrieve blockHash for '%s', block number %d: %s\n", testURL, header.Number, err)
					return
				}

				_, err = retriever.GetEvents(blockHash)

				if err != nil {
					log.Printf("Couldn't retrieve events for '%s', block number %d: %s\n", testURL, header.Number, err)
					return
				}

				header, err = api.RPC.Chain.GetHeader(header.ParentHash)

				if err != nil {
					log.Printf("Couldn't retrieve header for block number '%d' for '%s': %s\n", header.Number, testURL, err)

					return
				}

				processedBlockCount++

				if processedBlockCount%1000 == 0 {
					log.Printf("Retrieved events for %d blocks for '%s' so far, last block number %d\n", processedBlockCount, testURL, header.Number)
				}
			}
		}()
	}

	wg.Wait()
}
