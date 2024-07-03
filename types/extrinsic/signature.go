package extrinsic

import (
	"fmt"
	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

type Signature struct {
	Signer       types.MultiAddress
	Signature    types.MultiSignature
	SignedFields []*SignedField
}

func (s Signature) Encode(encoder scale.Encoder) error {
	if err := encoder.Encode(s.Signer); err != nil {
		return err
	}

	if err := encoder.Encode(s.Signature); err != nil {
		return err
	}

	for _, signedField := range s.SignedFields {
		if err := encoder.Encode(signedField.Value); err != nil {
			return fmt.Errorf("unable to encode signed field: %w", err)
		}
	}

	return nil
}

func createSignature(signer types.MultiAddress, sig types.MultiSignature, signedFields []*SignedField) *Signature {
	return &Signature{
		Signer:       signer,
		Signature:    sig,
		SignedFields: signedFields,
	}
}
