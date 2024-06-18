package types

import (
	"errors"
	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
)

type CheckMetadataMode byte

var (
	CheckMetadataModeDisabled CheckMetadataMode = 0
	CheckMetadataModeEnabled  CheckMetadataMode = 1
)

func (m CheckMetadataMode) Encode(encoder scale.Encoder) error {
	switch m {
	case 0:
		return encoder.PushByte(0)
	case 1:
		return encoder.PushByte(1)
	default:
		return errors.New("unsupported check metadata mode")
	}
}

type CheckMetadataHash struct {
	Hash Option[H256]
}

func (c CheckMetadataHash) Encode(encoder *scale.Encoder) error {
	if c.Hash.HasValue() {
		_, hash := c.Hash.Unwrap()

		return encoder.Encode(hash)
	}

	return encoder.PushByte(0)
}
