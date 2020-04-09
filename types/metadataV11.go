package types

import (
	"context"
	"github.com/centrifuge/go-substrate-rpc-client/scale"
)

// Modelled after packages/types/src/Metadata/v10/toV11.ts
type MetadataV11 struct {
	MetadataV10
	Extrinsic ExtrinsicV11
}

// Modelled after packages/types/src/Metadata/v10/toV11.ts
type ExtrinsicV11 struct {
	Version          uint8
	SignedExtensions []string
}

func (m *MetadataV11) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&m.MetadataV10)
	if err != nil {
		return err
	}
	return decoder.Decode(&m.Extrinsic)
}

func (m MetadataV11) Encode(ctx context.Context, encoder scale.Encoder) error {
	err := encoder.Encode(ctx, m.MetadataV10)
	if err != nil {
		return err
	}
	return encoder.Encode(ctx, m.Extrinsic)
}

func (e *ExtrinsicV11) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&e.Version)
	if err != nil {
		return err
	}

	return decoder.Decode(&e.SignedExtensions)
}

func (e ExtrinsicV11) Encode(ctx context.Context, encoder scale.Encoder) error {
	err := encoder.Encode(ctx, e.Version)
	if err != nil {
		return err
	}

	return encoder.Encode(ctx, e.SignedExtensions)
}
