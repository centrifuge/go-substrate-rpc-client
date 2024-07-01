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

type DynamicExtrinsicValidationFn func(payload *DynamicExtrinsicPayload, signature *DynamicExtrinsicSignature) error

// DynamicExtrinsicValidationFns holds a DynamicExtrinsicValidationFn for the supported signed extension types.
var DynamicExtrinsicValidationFns = map[SignedExtensionName]DynamicExtrinsicValidationFn{
	CheckMortalitySignedExtension: func(payload *DynamicExtrinsicPayload, signature *DynamicExtrinsicSignature) error {
		if payload.Era != nil && payload.BlockHash != nil && signature.Era != nil {
			return nil
		}

		return validationErr(CheckMortalitySignedExtension, "WithEra")
	},
	CheckEraSignedExtension: func(payload *DynamicExtrinsicPayload, signature *DynamicExtrinsicSignature) error {
		if payload.Era != nil && payload.BlockHash != nil && signature.Era != nil {
			return nil
		}

		return validationErr(CheckEraSignedExtension, "WithEra")
	},
	CheckNonceSignedExtension: func(payload *DynamicExtrinsicPayload, signature *DynamicExtrinsicSignature) error {
		if payload.Nonce != nil && signature.Nonce != nil {
			return nil
		}

		return validationErr(CheckNonceSignedExtension, "WithNonce")
	},
	ChargeTransactionPaymentSignedExtension: func(payload *DynamicExtrinsicPayload, signature *DynamicExtrinsicSignature) error {
		if payload.Tip != nil && signature.Tip != nil {
			return nil
		}

		return validationErr(ChargeTransactionPaymentSignedExtension, "WithTip")
	},
	CheckMetadataHashSignedExtension: func(payload *DynamicExtrinsicPayload, signature *DynamicExtrinsicSignature) error {
		if payload.CheckMetadataMode != nil && payload.CheckMetadataHash != nil && signature.CheckMetadataMode != nil {
			return nil
		}

		return validationErr(CheckMetadataHashSignedExtension, "WithMetadataMode")
	},
	CheckSpecVersionSignedExtension: func(payload *DynamicExtrinsicPayload, _ *DynamicExtrinsicSignature) error {
		if payload.SpecVersion != nil {
			return nil
		}

		return validationErr(CheckSpecVersionSignedExtension, "WithSpecVersion")

	},
	CheckTxVersionSignedExtension: func(payload *DynamicExtrinsicPayload, _ *DynamicExtrinsicSignature) error {
		if payload.TransactionVersion != nil {
			return nil
		}

		return validationErr(CheckTxVersionSignedExtension, "WithTransactionVersion")
	},
	CheckGenesisSignedExtension: func(payload *DynamicExtrinsicPayload, _ *DynamicExtrinsicSignature) error {
		if payload.GenesisHash != nil {
			return nil
		}

		return validationErr(CheckGenesisSignedExtension, "WithGenesisHash")
	},
	// There's nothing that we can check in the payload or signature in the following cases, however, these are added to
	// ensure that the extension is acknowledged and that the check is performed successfully.
	CheckNonZeroSenderSignedExtension: func(payload *DynamicExtrinsicPayload, signature *DynamicExtrinsicSignature) error {
		return nil
	},
	CheckWeightSignedExtension: func(payload *DynamicExtrinsicPayload, signature *DynamicExtrinsicSignature) error {
		return nil
	},
	PreBalanceTransferExtensionSignedExtension: func(payload *DynamicExtrinsicPayload, signature *DynamicExtrinsicSignature) error {
		return nil
	},
}

func validationErr(signedExtensionName SignedExtensionName, optFnName string) error {
	return fmt.Errorf("the signed extension check for '%s' failed, please make sure that the '%s' opt is provided", signedExtensionName, optFnName)
}

// checkDynamicExtrinsicData confirms that the DynamicExtrinsicPayload and DynamicExtrinsicSignature contain all the necessary
// data for the signed extensions that are expected based on the metadata.
func checkDynamicExtrinsicData(meta *types.Metadata, payload *DynamicExtrinsicPayload, signature *DynamicExtrinsicSignature) error {
	for _, signedExtension := range meta.AsMetadataV14.Extrinsic.SignedExtensions {
		signedExtensionType, ok := meta.AsMetadataV14.EfficientLookup[signedExtension.Type.Int64()]

		if !ok {
			return fmt.Errorf("signed extension type '%d' is not defined", signedExtension.Type.Int64())
		}

		signedExtensionName := signedExtensionType.Path[len(signedExtensionType.Path)-1]

		validationFn, ok := DynamicExtrinsicValidationFns[SignedExtensionName(signedExtensionName)]

		if !ok {
			return fmt.Errorf("signed extension '%s' is not supported", signedExtensionName)
		}

		if err := validationFn(payload, signature); err != nil {
			return err
		}
	}

	return nil
}
