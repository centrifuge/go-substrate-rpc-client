package extrinsic

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/extrinsic/extensions"
)

// DynamicExtrinsicSigningOption defines a function that mutates the DynamicExtrinsicPayload and DynamicExtrinsicSignature of
// a DynamicExtrinsic.
type DynamicExtrinsicSigningOption func(*DynamicExtrinsicPayload, *DynamicExtrinsicSignature)

func WithEra(era types.ExtrinsicEra, blockHash types.Hash) DynamicExtrinsicSigningOption {
	return func(payload *DynamicExtrinsicPayload, signature *DynamicExtrinsicSignature) {
		payload.Era = &era
		payload.BlockHash = &blockHash
		signature.Era = &era
	}
}

func WithNonce(nonce types.UCompact) DynamicExtrinsicSigningOption {
	return func(payload *DynamicExtrinsicPayload, signature *DynamicExtrinsicSignature) {
		payload.Nonce = &nonce
		signature.Nonce = &nonce
	}
}

func WithMetadataMode(mode extensions.CheckMetadataMode, metadataHash extensions.CheckMetadataHash) DynamicExtrinsicSigningOption {
	return func(payload *DynamicExtrinsicPayload, signature *DynamicExtrinsicSignature) {
		payload.CheckMetadataMode = &mode
		payload.CheckMetadataHash = &metadataHash
		signature.CheckMetadataMode = &mode
	}
}

func WithTip(tip types.UCompact) DynamicExtrinsicSigningOption {
	return func(payload *DynamicExtrinsicPayload, signature *DynamicExtrinsicSignature) {
		payload.Tip = &tip
		signature.Tip = &tip
	}
}

func WithSpecVersion(specVersion types.U32) DynamicExtrinsicSigningOption {
	return func(payload *DynamicExtrinsicPayload, signature *DynamicExtrinsicSignature) {
		payload.SpecVersion = &specVersion
	}
}

func WithTransactionVersion(transactionVersion types.U32) DynamicExtrinsicSigningOption {
	return func(payload *DynamicExtrinsicPayload, signature *DynamicExtrinsicSignature) {
		payload.TransactionVersion = &transactionVersion
	}
}

func WithGenesisHash(genesisHash types.Hash) DynamicExtrinsicSigningOption {
	return func(payload *DynamicExtrinsicPayload, signature *DynamicExtrinsicSignature) {
		payload.GenesisHash = &genesisHash
	}
}
