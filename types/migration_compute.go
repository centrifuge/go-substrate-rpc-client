package types

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
)

// MigrationCompute is an enum describing how a migration was computed.
type MigrationCompute struct {
	IsSigned bool
	IsAuto   bool
}

func (m *MigrationCompute) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch b {
	case 0:
		m.IsSigned = true
	case 1:
		m.IsAuto = true
	}

	return nil
}
