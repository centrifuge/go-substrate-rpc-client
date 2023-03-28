//go:build live

package retriever

import (
	"log"
	"sync"
	"testing"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
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

			//blockHash, err := api.RPC.Chain.GetBlockHashLatest()
			//
			//if err != nil {
			//	log.Printf("Couldn't get latest block hash for '%s': %s\n", testURL, err)
			//	return
			//}

			//statemint = 3505698
			//acala = 3238492

			blockHash, err := api.RPC.Chain.GetBlockHash(3505698)

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

				_, err = extrinsicRetriever.GetExtrinsics(blockHash)

				if err != nil {
					log.Printf("Couldn't retrieve extrinsics for '%s', block number %d: %s\n", testURL, block.Block.Header.Number, err)
					return
				}

				previousBlock = block

				blockHash = block.Block.Header.ParentHash

				processedBlockCount++

				if processedBlockCount%500 == 0 {
					log.Printf("Retrieved calls for %d blocks for '%s' so far, last block number %d\n", processedBlockCount, testURL, block.Block.Header.Number)
				}
			}
		}()
	}

	wg.Wait()
}
