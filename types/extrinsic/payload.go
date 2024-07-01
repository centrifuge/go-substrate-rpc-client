package extrinsic

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/extrinsic/extensions"
)

// DynamicExtrinsicPayload is the payload type used by the DynamicExtrinsic.
type DynamicExtrinsicPayload struct {
	Method             types.BytesBare
	Era                *types.ExtrinsicEra
	Nonce              *types.UCompact
	CheckMetadataMode  *extensions.CheckMetadataMode
	Tip                *types.UCompact
	SpecVersion        *types.U32
	TransactionVersion *types.U32
	GenesisHash        *types.Hash
	BlockHash          *types.Hash
	CheckMetadataHash  *extensions.CheckMetadataHash
}

// Sign encodes the payload and then signs the encoded bytes using the provided signer.
func (e *DynamicExtrinsicPayload) Sign(signer signature.KeyringPair) (types.Signature, error) {
	b, err := codec.Encode(e)
	if err != nil {
		return types.Signature{}, err
	}

	sig, err := signature.Sign(b, signer.URI)
	return types.NewSignature(sig), err
}

// Encode encodes the payload to Scale.
func (e *DynamicExtrinsicPayload) Encode(encoder scale.Encoder) error {
	if err := encoder.Encode(e.Method); err != nil {
		return err
	}

	if e.Era != nil {
		if err := encoder.Encode(e.Era); err != nil {
			return err
		}
	}

	if e.Nonce != nil {
		if err := encoder.Encode(e.Nonce); err != nil {
			return err
		}
	}

	if e.CheckMetadataMode != nil {
		if err := encoder.Encode(e.CheckMetadataMode); err != nil {
			return err
		}
	}

	if e.Tip != nil {
		if err := encoder.Encode(e.Tip); err != nil {
			return err
		}
	}

	if e.SpecVersion != nil {
		if err := encoder.Encode(e.SpecVersion); err != nil {
			return err
		}
	}

	if e.TransactionVersion != nil {
		if err := encoder.Encode(e.TransactionVersion); err != nil {
			return err
		}
	}

	if e.GenesisHash != nil {
		if err := encoder.Encode(e.GenesisHash); err != nil {
			return err
		}
	}

	if e.BlockHash != nil {
		if err := encoder.Encode(e.BlockHash); err != nil {
			return err
		}
	}

	if e.CheckMetadataHash != nil {
		if err := encoder.Encode(e.CheckMetadataHash); err != nil {
			return err
		}
	}

	return nil
}
