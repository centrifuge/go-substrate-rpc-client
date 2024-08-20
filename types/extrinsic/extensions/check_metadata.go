package extensions

import (
	"errors"

	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

type CheckMetadataMode byte

var (
	CheckMetadataModeDisabled CheckMetadataMode
	CheckMetadataModeEnabled  CheckMetadataMode = 1
)

func (m CheckMetadataMode) Encode(encoder scale.Encoder) error {
	switch m {
	case CheckMetadataModeDisabled:
		return encoder.PushByte(0)
	case CheckMetadataModeEnabled:
		return encoder.PushByte(1)
	default:
		return errors.New("unsupported check metadata mode")
	}
}

type CheckMetadataHash struct {
	Hash types.Option[types.H256]
}

func (c CheckMetadataHash) Encode(encoder *scale.Encoder) error {
	if c.Hash.HasValue() {
		_, hash := c.Hash.Unwrap()

		return encoder.Encode(hash)
	}

	return encoder.PushByte(0)
}
