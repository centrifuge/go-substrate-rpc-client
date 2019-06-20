package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"os/exec"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/vimukthi-git/go-substrate"
	"github.com/vimukthi-git/go-substrate/scalecodec"
)

type Signature struct {
	Hash [64]byte
}

func NewSignature(b []byte) *Signature {
	s := &Signature{}
	copy(s.Hash[:], b)
	return s
}

func (s *Signature) ParityDecode(decoder scalecodec.Decoder) {
	decoder.Read(s.Hash[:])
}

func (s Signature) ParityEncode(encoder scalecodec.Encoder) {
	encoder.Write(s.Hash[:])
}

type MortalEra struct {
	Period uint64
	Phase uint64
}

type ExtrinsicSignature struct {

	SignatureOptional uint8
	// 19              03             81                  ff
	// 0001 1001                      1000 0001           1111 1111
	Signer substrate.Address // 0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d
	Signature Signature
	Nonce uint64
	Era uint8 // era enum
}


func NewExtrinsicSignature(signature Signature, Nonce uint64) ExtrinsicSignature {
	return ExtrinsicSignature{Signature: signature, Nonce: Nonce}
}

func (e *ExtrinsicSignature) ParityDecode(decoder scalecodec.Decoder) {
	// length of the encoded signature struct
	//l := decoder.DecodeUintCompact()
	//fmt.Println(l)
	decoder.Decode(&e.SignatureOptional) // implement decodeExtrinsicSignature logic to derive if the request is signed

	b := decoder.ReadOneByte()
	// need to add other address representations from Address.decodeAddress
	if b == 255 {
		e.Signer = substrate.Address{}
		decoder.Decode(&e.Signer)
	}

	e.Signature = Signature{}
	decoder.Decode(&e.Signature)
	e.Nonce = decoder.DecodeUintCompact()
	// assuming immortal for now TODO
	decoder.Decode(&e.Era)

}

func (e ExtrinsicSignature) ParityEncode(encoder scalecodec.Encoder) {
	// always signed
	e.SignatureOptional = 129
	// Alice
	s, _ := hexutil.Decode("0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d")
	e.Signer = *substrate.NewAddress(s)
	e.Era = 0

	encoder.Encode(e.SignatureOptional)
	encoder.Encode(&e.Signer)
	encoder.Encode(&e.Signature)
	encoder.EncodeUintCompact(e.Nonce)
	encoder.Encode(e.Era)
}

type SignaturePayload struct {
	Nonce uint64
	Method Method
	Era uint8 // era enum
	//ImmortalEra []byte
	PriorBlock [32]byte
}

func (e SignaturePayload) ParityEncode(encoder scalecodec.Encoder) {
	encoder.EncodeUintCompact(e.Nonce)
	encoder.Encode(e.Method)
	encoder.Encode(e.Era)
	// encoder.Encode(e.ImmortalEra) // always immortal
	encoder.Write(e.PriorBlock[:])
}

type Args interface {
	scalecodec.Encodeable
	//scalecodec.Decodeable
}

type Method struct {

	CallIndex substrate.MethodIDX
	//  dynamic struct with the list of arguments defined as fields
	Args Args
}

func (e *Method) ParityDecode(decoder scalecodec.Decoder) {
	decoder.Decode(&e.CallIndex)
	e.Args = &AnchorParams{}
	decoder.Decode(e.Args)
}

func (m Method) ParityEncode(encoder scalecodec.Encoder) {
	encoder.Encode(&m.CallIndex)
	encoder.Encode(m.Args)
}

type Extrinsic struct {
	Signature ExtrinsicSignature
	Method Method
}

func (e *Extrinsic) ParityDecode(decoder scalecodec.Decoder) {
	// length (not used)
	decoder.DecodeUintCompact()

	e.Signature = ExtrinsicSignature{}
	decoder.Decode(&e.Signature)
	decoder.Decode(&e.Method)
}

func (e Extrinsic) ParityEncode(encoder scalecodec.Encoder) {
	var nonce uint64 = 102
	b := make([]byte, 0, 1000)
	bb := bytes.NewBuffer(b)
	tempEnc := scalecodec.NewEncoder(bb)
	pb, _ := hexutil.Decode("0xe9e76cdefbb3843a61814a77795078ebf5b07e0400ecd00d5390f3f0e4a96178")
	sigPay := SignaturePayload{
		Nonce: nonce,
		Method: e.Method,
		// Immortal
		Era: 0,
	}
	copy(sigPay.PriorBlock[:], pb)
	tempEnc.Encode(sigPay)
	bbb := bb.Bytes()
	encoded := hex.EncodeToString(bbb)


	// use "subKey" command for signature
	out, err := exec.Command("/Users/vimukthi/.cargo/bin/subkey", "sign-blob", encoded, "//Alice").Output()
	fmt.Println("/Users/vimukthi/.cargo/bin/subkey", "sign-blob", encoded, "//Alice")
	if err != nil {
		log.Fatal(err.Error())
	}

	v := string(out)
	fmt.Println("out", v)

	vs, err := hex.DecodeString(v)

	e.Signature = NewExtrinsicSignature(*NewSignature(vs), nonce)

	b = make([]byte, 0, 1000)
	bb = bytes.NewBuffer(b)
	tempEnc = scalecodec.NewEncoder(bb)
	tempEnc.Encode(&e.Signature)
	tempEnc.Encode(&e.Method)

	// encode with length prefix
	eb := bb.Bytes()
	encoder.EncodeUintCompact(uint64(len(eb)))
	encoder.Write(eb)
}


type AnchorParams struct {
	AnchorIDPreimage [32]byte
	DocRoot [32]byte
	Proof [32]byte
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
	client, err := substrate.Connect("ws://10.99.1.86:9944")
	if err != nil {
		panic(err)
	}

	// Alice
	//al, _ := hexutil.Decode("0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d")
	//alice := *substrate.NewAddress(al)
	//sys := substrate.NewSystemRPC(client)
	//_, err = sys.AccountNonce(alice)
	//if err != nil {
	//	panic(err)
	//}

	// 0xd133045f0efad58582772cbdb6f5f0cd6af7bb4bf1f30d039a4b18b4bdaf4901
	hs, err := hexutil.Decode("0xe9e76cdefbb3843a61814a77795078ebf5b07e0400ecd00d5390f3f0e4a96178")
	if err != nil {
		panic(err)
	}

	s := substrate.NewStateRPC(client)

	n, err := s.MetaData(hs)
	if err != nil {
		panic(err)
	}

	//s.Keys(hs)

	//err = client.Call(&res, "state_queryStorage", "5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY")
	//if err != nil {
	//	panic(err)
	//}

	a := NewAnchorParamsFromHex("0x0000000000000000000000000000000000000000000000000000000000000901", "0x0000000000000000000000000000000000000000000000000000000000000000", "0x0000000000000000000000000000000000000000000000000000000000000000")
	e := Extrinsic{Method:Method{CallIndex:n.Metadata.MethodIndex("kerplunk.commit"), Args:&a}}
	bb := make([]byte, 0, 1000)
	bbb := bytes.NewBuffer(bb)
	tempEnc := scalecodec.NewEncoder(bbb)
	tempEnc.Encode(&e)
	eb := hexutil.Encode(bbb.Bytes())
	fmt.Println(eb)
	var res string
	err = client.Call(&res, "author_submitExtrinsic", eb)
	if err != nil {
		panic(err.Error())
	}

	//wg := sync.WaitGroup{}
	//wg.Add(1)
	//go func(ctx context.Context) {
	//	for {
	//		select {
	//		case res := <-res:
	//			fmt.Println(res)
	//			wg.Done()
	//		case <-ctx.Done():
	//			wg.Done()
	//			break
	//		}
	//	}
	//}(ctx)
	//
	//wg.Wait()
	//
	fmt.Println("tx hash - ", res)
	//
	//canc()
	//sub.Unsubscribe()

}