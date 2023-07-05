//go:build live

package retriever

import (
	"errors"
	"log"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/registry"
	"github.com/centrifuge/go-substrate-rpc-client/v4/registry/exec"
	"github.com/centrifuge/go-substrate-rpc-client/v4/registry/parser"
	"github.com/centrifuge/go-substrate-rpc-client/v4/rpc/chain/generic"
	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

var (
	extTestChains = []testFnProvider{
		&testChain[
			types.MultiAddress,
			types.MultiSignature,
			generic.DefaultPaymentFields,
		]{
			url: "wss://fullnode.parachain.centrifuge.io",
		},
		&testChain[
			types.MultiAddress,
			types.MultiSignature,
			generic.DefaultPaymentFields,
		]{
			url: "wss://rpc.polkadot.io",
		},
		&testChain[
			types.MultiAddress,
			types.MultiSignature,
			generic.PaymentFieldsWithAssetID,
		]{
			url: "wss://statemint-rpc.polkadot.io",
		},
		&testChain[
			types.MultiAddress,
			AcalaMultiSignature,
			generic.DefaultPaymentFields,
		]{
			url: "wss://acala-rpc-0.aca-api.network",
		},
		&testChain[
			[20]byte,
			[65]byte,
			generic.DefaultPaymentFields,
		]{
			url: "wss://wss.api.moonbeam.network",
		},
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

	for _, c := range extTestChains {
		c := c

		wg.Add(1)

		go c.GetTestFn()(&wg, extrinsicsThreshold)
	}

	wg.Wait()
}

type testFn func(wg *sync.WaitGroup, extrinsicsThreshold int)

type testFnProvider interface {
	GetTestFn() testFn
}

type testChain[
	A any,
	S any,
	P any,
] struct {
	url string
}

func (t *testChain[A, S, P]) GetTestFn() testFn {
	return func(wg *sync.WaitGroup, extrinsicsThreshold int) {
		defer wg.Done()

		testURL := t.url

		api, err := gsrpc.NewSubstrateAPI(testURL)

		if err != nil {
			log.Printf("Couldn't connect to '%s': %s\n", testURL, err)
			return
		}

		chain := generic.NewChain[A, S, P, *generic.SignedBlock[A, S, P]](api.Client)

		extrinsicParser := parser.NewExtrinsicParser[A, S, P]()
		registryFactory := registry.NewFactory()

		chainExecutor := exec.NewRetryableExecutor[*generic.SignedBlock[A, S, P]](exec.WithRetryTimeout(1 * time.Second))
		extrinsicParsingExecutor := exec.NewRetryableExecutor[[]*parser.Extrinsic[A, S, P]](exec.WithMaxRetryCount(1))

		extrinsicRetriever, err := NewExtrinsicRetriever[A, S, P](
			extrinsicParser,
			chain,
			api.RPC.State,
			registryFactory,
			chainExecutor,
			extrinsicParsingExecutor,
		)

		if err != nil {
			log.Printf("Couldn't create extrinsic retriever: %s", err)
			return
		}

		header, err := api.RPC.Chain.GetHeaderLatest()

		if err != nil {
			log.Printf("Couldn't get latest header for '%s': %s\n", testURL, err)
			return
		}

		extrinsicsCount := 0

		for {
			blockHash, err := api.RPC.Chain.GetBlockHash(uint64(header.Number))

			if err != nil {
				log.Printf("Couldn't retrieve blockHash for '%s', block number %d: %s\n", testURL, header.Number, err)
				return
			}

			extrinsics, err := extrinsicRetriever.GetExtrinsics(blockHash)

			if err != nil {
				log.Printf("Couldn't retrieve extrinsics for '%s', block number %d: %s\n", testURL, header.Number, err)
			}

			log.Printf("Found %d extrinsics for '%s', at block number %d.\n", len(extrinsics), testURL, header.Number)

			extrinsicsCount += len(extrinsics)

			if extrinsicsCount >= extrinsicsThreshold {
				log.Printf("Retrieved a total of %d extrinsics for '%s', last block number %d. Stopping now.\n", extrinsicsCount, testURL, header.Number)

				return
			}

			header, err = api.RPC.Chain.GetHeader(header.ParentHash)

			if err != nil {
				log.Printf("Couldn't retrieve header for block number '%d' for '%s': %s\n", header.Number, testURL, err)

				return
			}
		}
	}
}

type AcalaMultiSignature struct {
	IsEd25519     bool
	AsEd25519     types.Signature
	IsSr25519     bool
	AsSr25519     types.Signature
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
