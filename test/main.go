package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/centrifuge/go-centrifuge/utils"
	"github.com/centrifuge/go-substrate-rpc-client"
	"github.com/centrifuge/go-substrate-rpc-client/scale"
	"github.com/centrifuge/go-substrate-rpc-client/system"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"golang.org/x/crypto/blake2b"
)

const (
	AnchorCommit = "anchor.commit"
	SubKeySign   = "sign-blob"

	// Adjust below params according to your env + chain state + requirement
	RPCEndPoint = "ws://127.0.0.1:9944"

	// SubKeyCmd subkey command to create signatures
	SubKeyCmd = "/Users/vimukthi/.cargo/bin/subkey"

	NumAnchorsPerThread = 2
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

type AnchorData struct {
	ID            [32]byte
	DocRoot       [32]byte
	AnchoredBlock uint64
}

func (a *AnchorData) Decode(decoder scale.Decoder) error {
	decoder.Read(a.ID[:])
	decoder.Read(a.DocRoot[:])
	decoder.Decode(&a.AnchoredBlock)
	return nil
}

func Anchors(client substrate.Client, anchorIDPreImage []byte) (*AnchorData, error) {
	h := blake2b.Sum256(anchorIDPreImage)
	m, err := client.MetaData(true)
	if err != nil {
		return nil, err
	}

	key, err := substrate.NewStorageKey(*m,"Anchor", "Anchors", h[:])
	if err != nil {
		return nil, err
	}

	s := substrate.NewStateRPC(client)
	res, err := s.Storage(key,  nil)
	if err != nil {
		return nil, err
	}

	tempDec := res.Decoder()
	a := AnchorData{}
	err = tempDec.Decode(&a)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func main() {
	// Connect the client.
	client, err := substrate.Connect(RPCEndPoint)
	if err != nil {
		panic(err)
	}
	alice, _ := hexutil.Decode(substrate.AlicePubKey)
	nonce, err := system.AccountNonce(client, alice)
	if err != nil {
		panic(err)
	}

	gs, err := system.BlockHash(client, 0)
	if err != nil {
		panic(err)
	}

	authRPC := substrate.NewAuthorRPC(client, gs, SubKeyCmd, SubKeySign)
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
				res, err := authRPC.SubmitExtrinsic(nonce, AnchorCommit, a)
				if err != nil {
					fmt.Printf("FAIL!!! anchor ID %s failed with %s\n", aID, err.Error())
					break
				} else {
					// verify anchor
					for i := 0; i < 10; i++ {
						time.Sleep(10 * time.Second)
						anc, err := Anchors(client, a.AnchorIDPreimage[:])
						fmt.Println(err)
						if anc != nil {
							fmt.Printf("SUCCESS!!! anchor %v\n", anc)
							break
						}
					}
					fmt.Printf("SUCCESS!!! anchor ID %s , tx hash %s\n", aID, res)
					atomic.AddUint64(&counter, 1)
					atomic.AddUint64(&nonce, 1)
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
