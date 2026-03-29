package extensions

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
	StorageWeightReclaimSignedExtension        SignedExtensionName = "StorageWeightReclaim"
	PrevalidateAttestsSignedExtension          SignedExtensionName = "PrevalidateAttests"
	CheckNetworkMembershipSignedExtension      SignedExtensionName = "CheckNetworkMembership"
	// AuthorizeCall is frame_system::AuthorizeCall (Polkadot SDK); empty payload fields for standard signed extrinsics.
	AuthorizeCallSignedExtension SignedExtensionName = "AuthorizeCall"
	// SetOrigin is frame_system::SetOrigin (transaction extension); no extra signing payload for normal signed transfers.
	SetOriginSignedExtension SignedExtensionName = "SetOrigin"
)
