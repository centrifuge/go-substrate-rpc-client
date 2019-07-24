package substrate

import (
	"bytes"
	"encoding/hex"
	"log"
	"os/exec"

	"github.com/centrifuge/go-substrate-rpc-client/scale"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	Alice       = "//Alice"
	AlicePubKey = "0xd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d"
)

type MortalEra struct {
	Period uint64
	Phase  uint64
}

type ExtrinsicSignature struct {
	SignatureOptional uint8
	Signer            Address
	Signature         Signature
	Nonce             uint64
	Era               uint8 // era enum
}

func NewExtrinsicSignature(signature Signature, Nonce uint64) ExtrinsicSignature {
	return ExtrinsicSignature{Signature: signature, Nonce: Nonce}
}

func (e *ExtrinsicSignature) Decode(decoder scale.Decoder) error {
	// length of the encoded signature struct
	//l := decoder.DecodeUintCompact()
	//fmt.Println(l)
	err := decoder.Decode(&e.SignatureOptional) // implement decodeExtrinsicSignature logic to derive if the request is signed
	if err != nil {
		return err
	}

	b, _ := decoder.ReadOneByte()
	// need to add other address representations from Address.decodeAddress
	if b == 255 {
		e.Signer = Address{}
		err = decoder.Decode(&e.Signer)
		if err != nil {
			return err
		}
	}

	e.Signature = Signature{}
	err = decoder.Decode(&e.Signature)
	if err != nil {
		return err
	}
	e.Nonce, _ = decoder.DecodeUintCompact()
	// assuming immortal for now TODO
	err = decoder.Decode(&e.Era)
	if err != nil {
		return err
	}

	return nil
}

func (e ExtrinsicSignature) Encode(encoder scale.Encoder) error {
	// always signed
	e.SignatureOptional = 129
	// Alice
	s, _ := hexutil.Decode(AlicePubKey)
	e.Signer = *NewAddress(s)
	e.Era = 0

	err := encoder.Encode(e.SignatureOptional)
	if err != nil {
		return err
	}

	err = encoder.Encode(&e.Signer)
	if err != nil {
		return err
	}

	err = encoder.Encode(&e.Signature)
	if err != nil {
		return err
	}

	err = encoder.EncodeUintCompact(e.Nonce)
	if err != nil {
		return err
	}

	err = encoder.Encode(e.Era)
	if err != nil {
		return err
	}

	return nil
}

type SignaturePayload struct {
	Nonce  uint64
	Method Method
	Era    uint8 // era enum
	//ImmortalEra []byte
	PriorBlock [32]byte
}

func (e SignaturePayload) Encode(encoder scale.Encoder) error {
	err := encoder.EncodeUintCompact(e.Nonce)
	if err != nil {
		return err
	}
	err = encoder.Encode(e.Method)
	if err != nil {
		return err
	}
	err = encoder.Encode(e.Era)
	if err != nil {
		return err
	}
	// encoder.Encode(e.ImmortalEra) // always immortal
	err = encoder.Write(e.PriorBlock[:])
	if err != nil {
		return err
	}
	return nil
}

type Args interface {
	scale.Encodeable
}

type Method struct {
	CallIndex MethodIDX
	//  dynamic struct with the list of arguments defined as fields
	Args Args
}

func NewMethod(name string, a Args, metadata MetadataVersioned) Method {
	// "kerplunk.commit"
	return Method{CallIndex: metadata.Metadata.MethodIndex(name), Args: a}
}

func (e *Method) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&e.CallIndex)
	if err != nil {
		return err
	}
	//e.Args = &AnchorParams{}
	err = decoder.Decode(e.Args)
	if err != nil {
		return err
	}
	return nil
}

func (m Method) Encode(encoder scale.Encoder) error {
	err := encoder.Encode(&m.CallIndex)
	if err != nil {
		return err
	}

	err = encoder.Encode(m.Args)
	if err != nil {
		return err
	}
	return nil
}

type Extrinsic struct {
	subKeyCMD  string
	subKeySign string
	Nonce      uint64

	GenesisBlock []byte
	Signature    ExtrinsicSignature
	Method       Method
}

func NewExtrinsic(subKeyCMD string, subKeySign string, accountNonce uint64, genesisBlock []byte, method Method) *Extrinsic {
	return &Extrinsic{subKeyCMD: subKeyCMD, subKeySign: subKeySign, Nonce: accountNonce, GenesisBlock: genesisBlock, Method: method}
}

func (e *Extrinsic) Decode(decoder scale.Decoder) error {
	// length (not used)
	_, err := decoder.DecodeUintCompact()
	if err != nil {
		return err
	}

	e.Signature = ExtrinsicSignature{}
	err = decoder.Decode(&e.Signature)
	if err != nil {
		return err
	}

	err = decoder.Decode(&e.Method)
	if err != nil {
		return err
	}

	return nil
}

func (e Extrinsic) Encode(encoder scale.Encoder) error {
	b := make([]byte, 0, 1000)
	bb := bytes.NewBuffer(b)
	tempEnc := scale.NewEncoder(bb)

	sigPay := SignaturePayload{
		Nonce:  e.Nonce,
		Method: e.Method,
		// Immortal
		Era: 0,
	}
	copy(sigPay.PriorBlock[:], e.GenesisBlock)
	err := tempEnc.Encode(sigPay)
	if err != nil {
		return err
	}
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
	tempEnc = scale.NewEncoder(bb)
	err = tempEnc.Encode(&e.Signature)
	if err != nil {
		return err
	}
	err = tempEnc.Encode(&e.Method)
	if err != nil {
		return err
	}

	// encode with length prefix
	eb := bb.Bytes()
	err = encoder.EncodeUintCompact(uint64(len(eb)))
	if err != nil {
		return err
	}
	err = encoder.Write(eb)
	if err != nil {
		return err
	}
	return nil
}

type Author struct {
	client Client
	genesisBlock []byte

	subKeyCMD  string
	subKeySign string
}

func NewAuthorRPC(client Client, genesisBlock []byte, subKeyCMD, SubKeySign string) *Author {
	return &Author{client, genesisBlock, subKeyCMD, SubKeySign}
}

func (a *Author) SubmitExtrinsic(accountNonce uint64, method string, args Args) (string, error) {
	m, err := a.client.MetaData(true)
	if err != nil {
		return "", err
	}
	e := NewExtrinsic(a.subKeyCMD, a.subKeySign, accountNonce, a.genesisBlock, NewMethod(method, args, *m))
	bb := make([]byte, 0, 1000)
	bbb := bytes.NewBuffer(bb)
	tempEnc := scale.NewEncoder(bbb)
	err = tempEnc.Encode(&e)
	if err != nil {
		return "", err
	}
	eb := hexutil.Encode(bbb.Bytes())
	// fmt.Println(eb)

	var res string
	err = a.client.Call(&res, "author_submitExtrinsic", eb)
	if err != nil {
		return "", err
	}

	return res, nil
}
