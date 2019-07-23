package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/centrifuge/go-centrifuge/utils"
	"github.com/centrifuge/go-substrate-rpc-client"
	"github.com/centrifuge/go-substrate-rpc-client/scale"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"golang.org/x/crypto/blake2b"
)

const (
	AnchorCommit = "anchor.commit"
	SubKeySign   = "sign-blob"

	// Adjust below params according to your env + chain state + requirement

	RPCEndPoint = "ws://127.0.0.1:9944"

	// TODO query these from the chain
	// 0x703bffa9bc816a5cd06ab8a95c2ff74e7a60cdaf269ffd64d325eb193266a656
	GenesisBlock = "0xcedd9988967d60f287700453839ce435584407b4a654040cd3c7f46b1bf61134"
	// BestBlock is the earliest block thats not already pruned
	BestBlock = "0xec5c55bdb4733cd3e6278284625379cbfd8ebc706e8a74529358d5673c434234"
	// StartNonce is the current account nonce for Alice (can't use other accounts for now)
	StartNonce = 2
	// SubKeyCmd subkey command to create signatures
	SubKeyCmd = "/Users/vimukthi/.cargo/bin/subkey"

	NumAnchorsPerThread = 1
	Concurrency         = 1
)

type AnchorParams struct {
	AnchorIDPreimage [32]byte
	DocRoot          [32]byte
	Proof            [32]byte
}

func NewRandomAnchor() AnchorParams {
	ap := AnchorParams{}
	copy(ap.AnchorIDPreimage[:], utils.RandomSlice(32))
	copy(ap.DocRoot[:], utils.RandomSlice(32))
	copy(ap.Proof[:], utils.RandomSlice(32))
	return ap
}

func NewAnchorParamsFromHex(apre, docr, proof string) AnchorParams {
	a, _ := hexutil.Decode(apre)
	d, _ := hexutil.Decode(docr)
	p, _ := hexutil.Decode(proof)
	ap := AnchorParams{}
	copy(ap.AnchorIDPreimage[:], a)
	copy(ap.DocRoot[:], d)
	copy(ap.Proof[:], p)
	return ap
}

func (a *AnchorParams) AnchorIDHex() string {
	b := blake2b.Sum256(a.AnchorIDPreimage[:])
	return hexutil.Encode(b[:])
}

func (a *AnchorParams) Decode(decoder scale.Decoder) error {
	decoder.Read(a.AnchorIDPreimage[:])
	decoder.Read(a.DocRoot[:])
	decoder.Read(a.Proof[:])
	return nil
}

func (a AnchorParams) Encode(encoder scale.Encoder) error {
	encoder.Write(a.AnchorIDPreimage[:])
	encoder.Write(a.DocRoot[:])
	encoder.Write(a.Proof[:])
	return nil
}

func main() {
	// Connect the client.
	client, err := substrate.Connect(RPCEndPoint)
	if err != nil {
		panic(err)
	}

	hs, err := hexutil.Decode(BestBlock)
	if err != nil {
		panic(err)
	}

	s := substrate.NewStateRPC(client)
	n, err := s.MetaData(hs)
	if err != nil {
		panic(err)
	}

	gs, err := hexutil.Decode(GenesisBlock)
	if err != nil {
		panic(err)
	}
	authRPC := substrate.NewAuthorRPC(StartNonce, gs, SubKeyCmd, SubKeySign, *n, client)
	wg := sync.WaitGroup{}
	start := time.Now()
	wg.Add(Concurrency)
	var counter uint64
	for i := 0; i < Concurrency; i++ {
		go func() {
			for i := 0; i < NumAnchorsPerThread; i++ {
				// a := NewAnchorParamsFromHex("0x0000000000000000000000000000000000000000000000000000000000000901", "0x0000000000000000000000000000000000000000000000000000000000000000", "0x0000000000000000000000000000000000000000000000000000000000000000")
				a := NewRandomAnchor()
				aID := a.AnchorIDHex()
				// fmt.Println("submitting new anchor with anchor ID", a.AnchorIDHex())
				res, err := authRPC.SubmitExtrinsic(AnchorCommit, a)
				if err != nil {
					fmt.Printf("FAIL!!! anchor ID %s failed with %s\n", aID, err.Error())
					break
				} else {
					fmt.Printf("SUCCESS!!! anchor ID %s , tx hash %s\n", aID, res)
					atomic.AddUint64(&counter, 1)
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()
	elapsed := time.Since(start)
	tps := float64(counter) / elapsed.Seconds()
	fmt.Printf("Successful execution of %d transactions took %s, amounting to %f TPS", counter, elapsed, tps)
}
