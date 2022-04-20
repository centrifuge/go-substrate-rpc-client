package types

type CollatorID [32]U8
type CollatorSignature [64]U8

type CandidateDescriptor struct {
	ParachainID                  ParachainID
	RelayParent                  Hash
	CollatorID                   CollatorID
	PersistentValidationDataHash Hash
	PoVHash                      Hash
	ErasureRoot                  Hash
	CollatorSignature            CollatorSignature
	ParaHead                     Hash
	ValidationCodeHash           Hash
}

type CandidateReceipt struct {
	Descriptor CandidateDescriptor

	CommitmentsHash Hash
}
