package extrinsic

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/extrinsic/extensions"
)

type DynamicExtrinsicSignature struct {
	Signer            types.MultiAddress
	Signature         types.MultiSignature
	Era               *types.ExtrinsicEra
	Nonce             *types.UCompact
	Tip               *types.UCompact
	CheckMetadataMode *extensions.CheckMetadataMode
}

func (d DynamicExtrinsicSignature) Encode(encoder scale.Encoder) error {
	if err := encoder.Encode(d.Signer); err != nil {
		return err
	}

	if err := encoder.Encode(d.Signature); err != nil {
		return err
	}

	if d.Era != nil {
		if err := encoder.Encode(*d.Era); err != nil {
			return err
		}
	}

	if d.Nonce != nil {
		if err := encoder.Encode(*d.Nonce); err != nil {
			return err
		}
	}

	if d.Tip != nil {
		if err := encoder.Encode(*d.Tip); err != nil {
			return err
		}
	}

	if d.CheckMetadataMode != nil {
		if err := encoder.Encode(*d.CheckMetadataMode); err != nil {
			return err
		}
	}

	return nil
}
