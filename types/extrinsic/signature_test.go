package extrinsic

import (
	"bytes"
	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSignature_Encode(t *testing.T) {
	signer, err := types.NewMultiAddressFromAccountID(signature.TestKeyringPairAlice.PublicKey)
	assert.NoError(t, err)

	multiSignature := types.MultiSignature{IsSr25519: true, AsSr25519: types.SignatureHash{}}

	signedField := &SignedField{
		Name:    "signed_field",
		Value:   uint64(1),
		Mutated: true,
	}

	signature := Signature{
		Signer:       signer,
		Signature:    multiSignature,
		SignedFields: []*SignedField{signedField},
	}

	encodedSigner, err := codec.Encode(signer)
	assert.NoError(t, err)
	encodedMultiSignature, err := codec.Encode(multiSignature)
	assert.NoError(t, err)
	encodedSignedFieldValue, err := codec.Encode(signedField.Value)
	assert.NoError(t, err)

	expectedResult := append(encodedSigner, append(encodedMultiSignature, encodedSignedFieldValue...)...)

	b := bytes.NewBuffer(nil)

	encoder := scale.NewEncoder(b)

	err = signature.Encode(*encoder)
	assert.NoError(t, err)

	assert.Equal(t, expectedResult, b.Bytes())
}
