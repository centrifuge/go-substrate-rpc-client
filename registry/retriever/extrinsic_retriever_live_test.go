//go:build live

package retriever

import (
	"log"
	"sync"
	"testing"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
)

func TestLive_ExtrinsicRetriever_GetExtrinsics(t *testing.T) {
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

			extrinsicRetriever, err := NewDefaultExtrinsicRetriever(api.RPC.Chain, api.RPC.State)

			if err != nil {
				log.Printf("Couldn't create extrinsic retriever: %s", err)
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

				_, err = extrinsicRetriever.GetExtrinsics(blockHash)

				if err != nil {
					log.Printf("Couldn't retrieve extrinsics for '%s', block number %d: %s\n", testURL, header.Number, err)
				}

				header, err = api.RPC.Chain.GetHeader(header.ParentHash)

				if err != nil {
					log.Printf("Couldn't retrieve header for block number '%d' for '%s': %s\n", header.Number, testURL, err)

					return
				}

				processedBlockCount++

				if processedBlockCount%500 == 0 {
					log.Printf("Retrieved calls for %d blocks for '%s' so far, last block number %d\n", processedBlockCount, testURL, header.Number)
				}
			}
		}()
	}

	wg.Wait()
}
