package types

import "github.com/centrifuge/go-substrate-rpc-client/v4/scale"

// pub struct ClassMetadata<DepositBalance, StringLimit: Get<u32>> {
//	/// The balance deposited for this metadata.
//	///
//	/// This pays for the data stored in this struct.
//	pub(super) deposit: DepositBalance,
//	/// General information concerning this asset. Limited in length by `StringLimit`. This will
//	/// generally be either a JSON dump or the hash of some JSON which can be found on a
//	/// hash-addressable global publication system such as IPFS.
//	pub(super) data: BoundedVec<u8, StringLimit>,
//	/// Whether the asset metadata may be changed by a non Force origin.
//	pub(super) is_frozen: bool,
//}

const (
	// CentrifugeChainStringLimit is defined in the centrifuge-chain.
	CentrifugeChainStringLimit = 256
)

type ClassMetadata struct {
	Deposit  U128
	Data     [CentrifugeChainStringLimit]U8
	IsFrozen bool
}

func (c *ClassMetadata) Decode(decoder scale.Decoder) error {
	if err := decoder.Decode(&c.Deposit); err != nil {
		return err
	}

	if err := decoder.Decode(&c.Data); err != nil {
		return err
	}

	return decoder.Decode(&c.IsFrozen)
}

func (c ClassMetadata) Encode(encoder scale.Encoder) error {
	if err := encoder.Encode(c.Deposit); err != nil {
		return err
	}

	if err := encoder.Encode(c.Data); err != nil {
		return err
	}

	return encoder.Encode(c.IsFrozen)
}
