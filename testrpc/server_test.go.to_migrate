// +build tests

package testrpc

import (
	"testing"

	"github.com/centrifuge/go-centrifuge/utils"
	"github.com/centrifuge/go-substrate-rpc-client"
	"github.com/centrifuge/go-substrate-rpc-client/scale"
	"github.com/stretchr/testify/assert"
)

const (
	SubKeySign = "sign-blob"

	// SubKeyCmd subkey command to create signatures
	SubKeyCmd = "/Users/vimukthi/.cargo/bin/subkey"
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

func TestServer(t *testing.T) {
	// TODO local only for now until subkey is included in the build
	t.SkipNow()
	testServer := new(Server)
	host, _ := testServer.Init(GetTestMetaData(), nil)
	c, err := substrate.Connect(host)
	assert.NoError(t, err)

	a := substrate.NewAuthorRPC(c, utils.RandomSlice(32), SubKeyCmd, SubKeySign)
	_, err = a.SubmitExtrinsic(1, "anchor.commit", NewRandomAnchor())
	assert.NoError(t, err)

}
