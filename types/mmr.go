package types

// GenerateMMRProofResponse contains the generate proof rpc response
type GenerateMMRProofResponse struct {
	blockHash H256
	leaf      MMRLeaf
	proof     MMRProof
}

// MMRProof is a mmr proof
type MMRProof struct {
	// The index of the leaf the proof is for.
	leafIndex U64
	// Number of leaves in MMR, when the proof was generated.
	leafCount U64
	// Proof elements (hashes of siblings of inner nodes on the path to the leaf).
	items []H256
}

type MMRLeaf struct {
	parentNumberAndHash   ParentNumberAndHash
	parachainHeads        H256
	beefyNextAuthoritySet BeefyNextAuthoritySet
}

type ParentNumberAndHash struct {
	parentNumber ParentNumber
	hash         [32]U8
}

// TODO: The MMRLeaf is a Vec<u8>, so double-scale encoded which messes this first variable, the ParentNumber up.
type ParentNumber [6]U8

type BeefyNextAuthoritySet struct {
	id U64
	// Number of validators in the set.
	len U32
	// Merkle Root Hash build from BEEFY uncompressed AuthorityIds.
	root H256
}
