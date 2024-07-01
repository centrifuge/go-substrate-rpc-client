package extrinsic

import (
	"fmt"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

// SignedExtensionName is the type for the signed extension name present in the metadata.
type SignedExtensionName string

const (
	CheckNonZeroSenderSignedExtension          SignedExtensionName = "CheckNonZeroSender"
	CheckMortalitySignedExtension              SignedExtensionName = "CheckMortality"
	CheckEraSignedExtension                    SignedExtensionName = "CheckEra"
	CheckNonceSignedExtension                  SignedExtensionName = "CheckNonce"
	ChargeTransactionPaymentSignedExtension    SignedExtensionName = "ChargeTransactionPayment"
	ChargeAssetTxPaymentSignedExtension        SignedExtensionName = "ChargeAssetTxPayment"
	CheckMetadataHashSignedExtension           SignedExtensionName = "CheckMetadataHash"
	CheckSpecVersionSignedExtension            SignedExtensionName = "CheckSpecVersion"
	CheckTxVersionSignedExtension              SignedExtensionName = "CheckTxVersion"
	CheckGenesisSignedExtension                SignedExtensionName = "CheckGenesis"
	CheckWeightSignedExtension                 SignedExtensionName = "CheckWeight"
	PreBalanceTransferExtensionSignedExtension SignedExtensionName = "PreBalanceTransferExtension"
)

type DynamicExtrinsicValidationFn func(payload *DynamicExtrinsicPayload, signature *DynamicExtrinsicSignature) bool

// DynamicExtrinsicValidationFns holds a DynamicExtrinsicValidationFn for the supported signed extension types.
var DynamicExtrinsicValidationFns = map[SignedExtensionName]DynamicExtrinsicValidationFn{
	CheckMortalitySignedExtension: func(payload *DynamicExtrinsicPayload, signature *DynamicExtrinsicSignature) bool {
		return payload.Era != nil && payload.BlockHash != nil && signature.Era != nil
	},
	CheckEraSignedExtension: func(payload *DynamicExtrinsicPayload, signature *DynamicExtrinsicSignature) bool {
		return payload.Era != nil && payload.BlockHash != nil && signature.Era != nil
	},
	CheckNonceSignedExtension: func(payload *DynamicExtrinsicPayload, signature *DynamicExtrinsicSignature) bool {
		return payload.Nonce != nil && signature.Nonce != nil
	},
	ChargeTransactionPaymentSignedExtension: func(payload *DynamicExtrinsicPayload, signature *DynamicExtrinsicSignature) bool {
		return payload.Tip != nil && signature.Tip != nil
	},
	CheckMetadataHashSignedExtension: func(payload *DynamicExtrinsicPayload, signature *DynamicExtrinsicSignature) bool {
		return payload.CheckMetadataMode != nil && payload.CheckMetadataHash != nil && signature.CheckMetadataMode != nil
	},
	CheckSpecVersionSignedExtension: func(payload *DynamicExtrinsicPayload, _ *DynamicExtrinsicSignature) bool {
		return payload.SpecVersion != nil
	},
	CheckTxVersionSignedExtension: func(payload *DynamicExtrinsicPayload, _ *DynamicExtrinsicSignature) bool {
		return payload.TransactionVersion != nil
	},
	CheckGenesisSignedExtension: func(payload *DynamicExtrinsicPayload, _ *DynamicExtrinsicSignature) bool {
		return payload.GenesisHash != nil
	},
	// There's nothing that we can check in the payload or signature in the following cases, however, these are added to
	// ensure that the extension is acknowledged and that the check is performed successfully.
	CheckNonZeroSenderSignedExtension: func(payload *DynamicExtrinsicPayload, signature *DynamicExtrinsicSignature) bool {
		return true
	},
	CheckWeightSignedExtension: func(payload *DynamicExtrinsicPayload, signature *DynamicExtrinsicSignature) bool {
		return true
	},
	PreBalanceTransferExtensionSignedExtension: func(payload *DynamicExtrinsicPayload, signature *DynamicExtrinsicSignature) bool {
		return true
	},
}

// checkDynamicExtrinsicData confirms that the DynamicExtrinsicPayload and DynamicExtrinsicSignature contain all the necessary
// signed extensions that are expected, based on the metadata.
func checkDynamicExtrinsicData(meta *types.Metadata, payload *DynamicExtrinsicPayload, signature *DynamicExtrinsicSignature) error {
	for _, signedExtension := range meta.AsMetadataV14.Extrinsic.SignedExtensions {
		signedExtensionType, ok := meta.AsMetadataV14.EfficientLookup[signedExtension.Type.Int64()]

		if !ok {
			return fmt.Errorf("signed extension type %d is not defined", signedExtension.Type.Int64())
		}

		signedExtensionName := signedExtensionType.Path[len(signedExtensionType.Path)-1]

		validationFn, ok := DynamicExtrinsicValidationFns[SignedExtensionName(signedExtensionName)]

		if !ok {
			return fmt.Errorf("signed extension '%s' is not supported", signedExtensionName)
		}

		if valid := validationFn(payload, signature); !valid {
			return fmt.Errorf("check for signed extension '%s' failed", signedExtensionName)
		}
	}

	return nil
}
