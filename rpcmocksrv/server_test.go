package rpcmocksrv

import (
	"testing"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
)

const (
	SubKeySign = "sign-blob"

	// SubKeyCmd subkey command to create signatures
	SubKeyCmd = "/Users/vimukthi/.cargo/bin/subkey"
)

// type AnchorParams struct {
// 	AnchorIDPreimage [32]byte
// 	DocRoot          [32]byte
// 	Proof            [32]byte
// }

// func NewRandomAnchor() AnchorParams {
// 	ap := AnchorParams{}
// 	copy(ap.AnchorIDPreimage[:], utils.RandomSlice(32))
// 	copy(ap.DocRoot[:], utils.RandomSlice(32))
// 	copy(ap.Proof[:], utils.RandomSlice(32))
// 	return ap
// }

// func (a *AnchorParams) Decode(decoder scale.Decoder) error {
// 	decoder.Read(a.AnchorIDPreimage[:])
// 	decoder.Read(a.DocRoot[:])
// 	decoder.Read(a.Proof[:])
// 	return nil
// }

// func (a AnchorParams) Encode(encoder scale.Encoder) error {
// 	encoder.Write(a.AnchorIDPreimage[:])
// 	encoder.Write(a.DocRoot[:])
// 	encoder.Write(a.Proof[:])
// 	return nil
// }

type TestService struct {
}

func (ts *TestService) Ping(s string) string {
	return s
}

func TestServer(t *testing.T) {
	s := New()

	ts := new(TestService)
	err := s.RegisterName("testserv3", ts)
	assert.NoError(t, err)

	c, err := rpc.Dial(s.URL)
	assert.NoError(t, err)

	var res string
	err = c.Call(&res, "testserv3_ping", "hello")
	assert.NoError(t, err)

	assert.Equal(t, "hello", res)
}
