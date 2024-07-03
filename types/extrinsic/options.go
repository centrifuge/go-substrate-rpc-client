package extrinsic

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/extrinsic/extensions"
)

type SignedFieldValues map[SignedFieldName]any

type SigningOption func(vals SignedFieldValues)

func WithEra(era types.ExtrinsicEra, blockHash types.Hash) SigningOption {
	return func(vals SignedFieldValues) {
		vals[EraSignedField] = era
		vals[BlockHashSignedField] = blockHash
	}
}

func WithNonce(nonce types.UCompact) SigningOption {
	return func(vals SignedFieldValues) {
		vals[NonceSignedField] = nonce
	}
}

func WithMetadataMode(mode extensions.CheckMetadataMode, metadataHash extensions.CheckMetadataHash) SigningOption {
	return func(vals SignedFieldValues) {
		vals[CheckMetadataHashModeSignedField] = mode
		vals[CheckMetadataHashSignedField] = metadataHash
	}
}

func WithTip(tip types.UCompact) SigningOption {
	return func(vals SignedFieldValues) {
		vals[TipSignedField] = tip
	}
}

func WithSpecVersion(specVersion types.U32) SigningOption {
	return func(vals SignedFieldValues) {
		vals[SpecVersionSignedField] = specVersion
	}
}

func WithTransactionVersion(transactionVersion types.U32) SigningOption {
	return func(vals SignedFieldValues) {
		vals[TransactionVersionSignedField] = transactionVersion
	}
}

func WithGenesisHash(genesisHash types.Hash) SigningOption {
	return func(vals SignedFieldValues) {
		vals[GenesisHashSignedField] = genesisHash
	}
}
