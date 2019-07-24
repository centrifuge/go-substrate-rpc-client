package substrate

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/scale"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/minio/blake2b-simd"
	"github.com/stretchr/testify/assert"
	bbb "golang.org/x/crypto/blake2b"
)

func TestState_GetMetaData(t *testing.T) {
	s := State{nonetwork: true}
	res, err := s.MetaData([]byte{})
	assert.NoError(t, err)
	assert.Equal(t, "system", res.Metadata.Modules[0].Name)
	// fmt.Println(res)
}

type AnchorData struct {
	ID [32]byte
	DocRoot [32]byte
	AnchoredBlock uint64
}

func (a *AnchorData) Decode(decoder scale.Decoder) error {
	decoder.Read(a.ID[:])
	decoder.Read(a.DocRoot[:])
	decoder.Decode(&a.AnchoredBlock)
	return nil
}

func (a AnchorData) Encode(encoder scale.Encoder) error {
	//encoder.Write(a.AnchorIDPreimage[:])
	//encoder.Write(a.DocRoot[:])
	//encoder.Write(a.AnchoredBlock)
	return nil
}

func TestState_GetStorage_Anchors(t *testing.T) {
	c, _ := Connect("ws://127.0.0.1:9944")
	s := NewStateRPC(c)
	b, _ := hexutil.Decode("0x33e423980c9b37d048bd5fadbd4a2aeb95146922045405accc2f468d0ef96988")
	h, _ := hexutil.Decode("0x142d4b3d1946e4956b4bd5a5bfc906142e921b51415ceccb3c82b3bd3ff3daf1")

	m ,_ := s.MetaData(h)
	key, _ := NewStorageKey(*m,"System", "AccountNonce", b)
	res, _ := s.Storage(key,  nil)
	fmt.Println(res)

	buf := bytes.NewBuffer(res)
	tempDec := scale.NewDecoder(buf)
	a := AnchorData{}
	tempDec.Decode(&a)

}

func TestState_GetStorage_AccountNonce(t *testing.T) {
	c, _ := Connect("ws://127.0.0.1:9944")
	s := NewStateRPC(c)
	b, _ := hexutil.Decode(AlicePubKey)
	h, _ := hexutil.Decode("0x142d4b3d1946e4956b4bd5a5bfc906142e921b51415ceccb3c82b3bd3ff3daf1")

	m ,_ := s.MetaData(h)
	key, _ := NewStorageKey(*m,"System", "AccountNonce", b)
	res, _ := s.Storage(key,  nil)

	buf := bytes.NewBuffer(res)
	tempDec := scale.NewDecoder(buf)
	var nonce uint64
	tempDec.Decode(&nonce)
	fmt.Println(nonce)
}

func TestState_GetStorage_TimeNow(t *testing.T) {
	c, _ := Connect("ws://127.0.0.1:9944")
	fn := []byte("Timestamp Now")
	h, err := getStorageHasher("")
	if err != nil {
		panic(err)
	}
	h.Write(fn)
	key := create2xXxhash(fn, 2)
	// TODO ask why this is not needed for "Timestamp Now"?
	//tempEnc.EncodeUintCompact(uint64(len(key)))
	//tempEnc.Write(key)
	// key := buf.Bytes()
	s := NewStateRPC(c)
	res, _ := s.Storage(key,  nil)
	fmt.Println(res)
}

func TestBlake(t *testing.T) {
	bb, _ := hexutil.Decode("0x0000000000000000000000000000000000000000000000000000000000000001")
	b := blake2b.Sum256(bb)
	b2 := bbb.Sum256(bb)
	fmt.Println(hexutil.Encode(b[:]))
	fmt.Println(hexutil.Encode(b2[:]))
}

func blake128(b []byte) []byte {
	h, err := blake2b.New(&blake2b.Config{Size: 16})
	if err != nil {
		fmt.Println(err)
	}
	h.Write(b)
	return h.Sum(nil)
}







