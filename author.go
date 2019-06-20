package substrate

import (
	"bytes"
	"encoding/hex"
	"log"
	"os/exec"
	"sync"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/vimukthi-git/go-substrate/scalecodec"
)

const (
	Alice = "//Alice"
	AlicePubKey = "0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d"
)

type MortalEra struct {
	Period uint64
	Phase uint64
}

type ExtrinsicSignature struct {
	SignatureOptional uint8
	Signer Address
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
		e.Signer = Address{}
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
	s, _ := hexutil.Decode(AlicePubKey)
	e.Signer = *NewAddress(s)
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
}

type Method struct {

	CallIndex MethodIDX
	//  dynamic struct with the list of arguments defined as fields
	Args Args
}

func NewMethod(name string, a Args, metadata MetadataVersioned) Method {
	// "kerplunk.commit"
	return Method{CallIndex: metadata.Metadata.MethodIndex(name), Args:a}
}

func (e *Method) ParityDecode(decoder scalecodec.Decoder) {
	decoder.Decode(&e.CallIndex)
	//e.Args = &AnchorParams{}
	decoder.Decode(e.Args)
}

func (m Method) ParityEncode(encoder scalecodec.Encoder) {
	encoder.Encode(&m.CallIndex)
	encoder.Encode(m.Args)
}

type Extrinsic struct {
	subKeyCMD string
	subKeySign string
	Nonce uint64
	// BestKnownBlock genesis block
	BestKnownBlock []byte
	Signature      ExtrinsicSignature
	Method         Method
}

func NewExtrinsic(subKeyCMD string, subKeySign string, accountNonce uint64, bestKnownBlock []byte, method Method) *Extrinsic {
	return &Extrinsic{subKeyCMD: subKeyCMD, subKeySign: subKeySign, Nonce:accountNonce, BestKnownBlock: bestKnownBlock, Method:method}
}

func (e *Extrinsic) ParityDecode(decoder scalecodec.Decoder) {
	// length (not used)
	decoder.DecodeUintCompact()

	e.Signature = ExtrinsicSignature{}
	decoder.Decode(&e.Signature)
	decoder.Decode(&e.Method)
}

func (e Extrinsic) ParityEncode(encoder scalecodec.Encoder) {
	b := make([]byte, 0, 1000)
	bb := bytes.NewBuffer(b)
	tempEnc := scalecodec.NewEncoder(bb)

	sigPay := SignaturePayload{
		Nonce: e.Nonce,
		Method: e.Method,
		// Immortal
		Era: 0,
	}
	copy(sigPay.PriorBlock[:], e.BestKnownBlock)
	tempEnc.Encode(sigPay)
	bbb := bb.Bytes()
	encoded := hex.EncodeToString(bbb)


	// use "subKey" command for signature
	out, err := exec.Command(e.subKeyCMD, e.subKeySign, encoded, Alice).Output()
	// fmt.Println(SubKeyCmd, SubKeySign, encoded, Alice)
	if err != nil {
		log.Fatal(err.Error())
	}

	v := string(out)
	vs, err := hex.DecodeString(v)

	e.Signature = NewExtrinsicSignature(*NewSignature(vs), e.Nonce)

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

type Author struct {
	client Client
	meta MetadataVersioned

	// mu is an exclusive lock to manage the nonce
	mu sync.RWMutex

	subKeyCMD string
	subKeySign string

	// TODO obtain these using RPCs
	accountNonce uint64
	bestKnownBlock []byte

}

func NewAuthorRPC(startNonce uint64, bestKnownBlock []byte, subKeyCMD , SubKeySign string, meta MetadataVersioned, client Client) *Author {
	return &Author{ client, meta, sync.RWMutex{}, subKeyCMD, SubKeySign, startNonce, bestKnownBlock}
}

func (a *Author) SubmitExtrinsic(method string, args Args) (string, error) {
	a.mu.Lock()
	e :=  NewExtrinsic(a.subKeyCMD, a.subKeySign, a.accountNonce, a.bestKnownBlock, NewMethod(method, args, a.meta))
	a.accountNonce++
	a.mu.Unlock()
	bb := make([]byte, 0, 1000)
	bbb := bytes.NewBuffer(bb)
	tempEnc := scalecodec.NewEncoder(bbb)
	tempEnc.Encode(&e)
	eb := hexutil.Encode(bbb.Bytes())
	// fmt.Println(eb)

	var res string
	err := a.client.Call(&res, "author_submitExtrinsic", eb)
	if err != nil {
		return "", err
	}

	return res, nil
}