//go:build live

package retriever

import (
	"errors"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"log"
	"os"
	"strconv"
	"sync"
	"testing"
)

var (
	extrinsicTestURLs = []string{
		"wss://fullnode.parachain.centrifuge.io",
		"wss://rpc.polkadot.io",
		"wss://acala-rpc-0.aca-api.network",
		"wss://wss.api.moonbeam.network",
	}
)

func TestLive_ExtrinsicRetriever_GetExtrinsics(t *testing.T) {
	t.Parallel()

	var wg sync.WaitGroup

	extrinsicsThreshold, err := strconv.Atoi(os.Getenv("GSRPC_LIVE_TEST_EXTRINSICS_THRESHOLD"))
	if err != nil {
		extrinsicsThreshold = 50
		log.Printf("Env Var GSRPC_LIVE_TEST_EXTRINSICS_THRESHOLD not set, defaulting to %d", extrinsicsThreshold)
	}

	for _, testURL := range extrinsicTestURLs {
		testURL := testURL

		wg.Add(1)

		go func() {
			defer wg.Done()

			api, err := gsrpc.NewSubstrateAPI(testURL)

			if err != nil {
				log.Printf("Couldn't connect to '%s': %s\n", testURL, err)
				return
			}

			retriever, err := NewDefaultExtrinsicRetriever(api.RPC.Chain, api.RPC.State)

			if err != nil {
				log.Printf("Couldn't create extrinsic retriever: %s", err)
				return
			}

			header, err := api.RPC.Chain.GetHeaderLatest()

			if err != nil {
				log.Printf("Couldn't get latest header for '%s': %s\n", testURL, err)
				return
			}

			extrinsicCount := 0

			for {
				blockHash, err := api.RPC.Chain.GetBlockHash(uint64(header.Number))

				if err != nil {
					log.Printf("Couldn't retrieve blockHash for '%s', block number %d: %s\n", testURL, header.Number, err)
					return
				}

				extrinsics, err := retriever.GetExtrinsics(blockHash)

				if err != nil {
					log.Printf("Couldn't retrieve extrinsics for '%s', block number %d: %s\n", testURL, header.Number, err)
					return
				}

				log.Printf("Found %d extrinsics for '%s', at block number %d.\n", len(extrinsics), testURL, header.Number)

				extrinsicCount += len(extrinsics)

				if extrinsicCount > extrinsicsThreshold {
					log.Printf("Retrieved a total of %d extrinsics for '%s', last block number %d. Stopping now.\n", extrinsicCount, testURL, header.Number)

					return
				}

				header, err = api.RPC.Chain.GetHeader(header.ParentHash)

				if err != nil {
					log.Printf("Couldn't retrieve header for block number '%d' for '%s': %s\n", header.Number, testURL, err)

					return
				}
			}
		}()
	}

	wg.Wait()
}

type AcalaMultiSignature struct {
	IsEd25519     bool
	AsEd25519     types.SignatureHash
	IsSr25519     bool
	AsSr25519     types.SignatureHash
	IsEcdsa       bool
	AsEcdsa       types.EcdsaSignature
	IsEthereum    bool
	AsEthereum    [65]byte
	IsEip1559     bool
	AsEip1559     [65]byte
	IsAcalaEip712 bool
	AsAcalaEip712 [65]byte
}

func (m *AcalaMultiSignature) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch b {
	case 0:
		m.IsEd25519 = true
		err = decoder.Decode(&m.AsEd25519)
	case 1:
		m.IsSr25519 = true
		err = decoder.Decode(&m.AsSr25519)
	case 2:
		m.IsEcdsa = true
		err = decoder.Decode(&m.AsEcdsa)
	case 3:
		m.IsEthereum = true
		err = decoder.Decode(&m.AsEthereum)
	case 4:
		m.IsEip1559 = true
		err = decoder.Decode(&m.AsEip1559)
	case 5:
		m.IsAcalaEip712 = true
		err = decoder.Decode(&m.AsAcalaEip712)
	default:
		return errors.New("signature not supported")
	}

	if err != nil {
		return err
	}

	return nil
}
