package main

import (
	"fmt"
	"github.com/centrifuge/go-centrifuge/utils"
	"golang.org/x/crypto/blake2b"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/vimukthi-git/go-substrate"
	"github.com/vimukthi-git/go-substrate/scalecodec"
)

const (
	RPCEndPoint = "ws://127.0.0.1:9944"

	// TODO query these from the chain
	GenesisBlock  = "0x703bffa9bc816a5cd06ab8a95c2ff74e7a60cdaf269ffd64d325eb193266a656"
	BestBlock  = "0xcd515661ff266920416ba9f1d48d8c532c586ab0101cbc987fbd74b4503abe87"
	StartNonce = 83

	AnchorCommit = "kerplunk.commit"
)

type AnchorParams struct {
	AnchorIDPreimage [32]byte
	DocRoot [32]byte
	Proof [32]byte
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

func (a *AnchorParams) ParityDecode(decoder scalecodec.Decoder)  {
	decoder.Read(a.AnchorIDPreimage[:])
	decoder.Read(a.DocRoot[:])
	decoder.Read(a.Proof[:])
}

func (a AnchorParams) ParityEncode(encoder scalecodec.Encoder) {
	encoder.Write(a.AnchorIDPreimage[:])
	encoder.Write(a.DocRoot[:])
	encoder.Write(a.Proof[:])
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
	authRPC := substrate.NewAuthorRPC(StartNonce, gs, *n, client)

	for i := 0; i < 10; i++ {
		// a := NewAnchorParamsFromHex("0x0000000000000000000000000000000000000000000000000000000000000901", "0x0000000000000000000000000000000000000000000000000000000000000000", "0x0000000000000000000000000000000000000000000000000000000000000000")
		a := NewRandomAnchor()
		aID := a.AnchorIDHex()
		fmt.Println("submitting new anchor with anchor ID %s", a.AnchorIDHex())
		res, err := authRPC.SubmitExtrinsic(AnchorCommit, a)
		if err != nil {
			fmt.Println("anchor ID %s failed with %s", aID, err.Error())
		} else {
			fmt.Println("tx hash - ", res)
		}
	}
}