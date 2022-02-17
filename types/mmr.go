package types

import "encoding/json"

// GenerateMmrProofResponse contains the generate proof rpc response
type GenerateMmrProofResponse struct {
	BlockHash H256
	Leaf      MmrLeaf
	Proof     MmrProof
}

// GenerateMmrBatchProofResponse contains the generate batch proof rpc response
type GenerateMmrBatchProofResponse struct {
	BlockHash H256
	Leaves    []LeafWithIndex
	Proof     MmrProof
}

type LeafWithIndex struct {
	Index U64
	Leaf  MmrLeaf
}

// UnmarshalJSON fills u with the JSON encoded byte array given by b
func (d *GenerateMmrProofResponse) UnmarshalJSON(bz []byte) error {
	var tmp struct {
		BlockHash string `json:"blockHash"`
		Leaf      string `json:"leaf"`
		Proof     string `json:"proof"`
	}
	if err := json.Unmarshal(bz, &tmp); err != nil {
		return err
	}
	err := DecodeFromHexString(tmp.BlockHash, &d.BlockHash)
	if err != nil {
		return err
	}
	var encodedLeaf MMREncodableOpaqueLeaf
	err = DecodeFromHexString(tmp.Leaf, &encodedLeaf)
	if err != nil {
		return err
	}
	err = DecodeFromBytes(encodedLeaf, &d.Leaf)
	if err != nil {
		return err
	}
	err = DecodeFromHexString(tmp.Proof, &d.Proof)
	if err != nil {
		return err
	}
	return nil
}

// UnmarshalJSON fills u with the JSON encoded byte array given by b
func (d *GenerateMmrBatchProofResponse) UnmarshalJSON(bz []byte) error {
	var tmp struct {
		BlockHash string `json:"blockHash"`
		Leaves    string `json:"leaves"`
		Proof     string `json:"proof"`
	}
	if err := json.Unmarshal(bz, &tmp); err != nil {
		return err
	}
	err := DecodeFromHexString(tmp.BlockHash, &d.BlockHash)
	if err != nil {
		return err
	}
	var encodedLeaf MMREncodableOpaqueLeaf
	err = DecodeFromHexString(tmp.Leaves, &encodedLeaf)
	if err != nil {
		return err
	}
	err = DecodeFromBytes(encodedLeaf, &d.Leaves)
	if err != nil {
		return err
	}
	err = DecodeFromHexString(tmp.Proof, &d.Proof)
	if err != nil {
		return err
	}
	return nil
}

type MMREncodableOpaqueLeaf Bytes

// MmrProof is a MMR proof
type MmrProof struct {
	// The index of the leaf the proof is for.
	LeafIndex U64
	// Number of leaves in MMR, when the proof was generated.
	LeafCount U64
	// Proof elements (hashes of siblings of inner nodes on the path to the leaf).
	Items []H256
}

type MmrLeaf struct {
	Version               MMRLeafVersion
	ParentNumberAndHash   ParentNumberAndHash
	BeefyNextAuthoritySet BeefyNextAuthoritySet
	ParachainHeads        H256
}

type MMRLeafVersion U8

type ParentNumberAndHash struct {
	ParentNumber U32
	Hash         H256
}

type BeefyNextAuthoritySet struct {
	// ID
	ID U64
	// Number of validators in the set.
	Len U32
	// Merkle Root Hash build from BEEFY uncompressed AuthorityIds.
	Root H256
}
