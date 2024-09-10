package extrinsic

import (
	libErr "github.com/centrifuge/go-substrate-rpc-client/v4/error"
	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

const (
	ErrSignatureFieldEncoding = libErr.Error("signature field encoding failed")
)

// Signature holds all the relevant fields for an extrinsic signature.
type Signature struct {
	Signer       types.MultiAddress
	Signature    types.MultiSignature
	SignedFields []*SignedField
}

// Encode is encoding the Signer, Signature, and SignedFields.
//
// Note - the ordering of the SignedFields is the order in which they are provided in
// the metadata.
func (s Signature) Encode(encoder scale.Encoder) error {
	if err := encoder.Encode(s.Signer); err != nil {
		return err
	}

	if err := encoder.Encode(s.Signature); err != nil {
		return err
	}

	for _, signedField := range s.SignedFields {
		if err := encoder.Encode(signedField.Value); err != nil {
			return ErrSignatureFieldEncoding.Wrap(err)
		}
	}

	return nil
}
