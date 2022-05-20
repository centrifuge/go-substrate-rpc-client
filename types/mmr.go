package types

import (
	"encoding/json"
)

// GenerateMmrProofResponse contains the generate proof rpc response
type GenerateMmrProofResponse struct {
	BlockHash H256
	Leaf      MmrLeaf
	Proof     MmrProof
}

// GenerateMmrBatchProofResponse contains the generate batch proof rpc response
type GenerateMmrBatchProofResponse struct {
	BlockHash H256
	Leaves    []MmrLeaf
	Proof     MmrBatchProof
}

type OpaqueLeafWithIndex struct {
	Leaf  []byte
	Index uint64
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

	var opaqueLeaves [][]byte
	err = DecodeFromHexString(tmp.Leaves, &opaqueLeaves)
	if err != nil {
		return err
	}
	for _, leaf := range opaqueLeaves {
		var mmrLeaf MmrLeaf
		err := DecodeFromBytes(leaf, &mmrLeaf)
		if err != nil {
			return err
		}
		d.Leaves = append(d.Leaves, mmrLeaf)
	}
	err = DecodeFromHexString(tmp.Proof, &d.Proof)
	if err != nil {
		return err
	}
	return nil
}

type MMREncodableOpaqueLeaf []byte

// MmrProof is a MMR proof
type MmrProof struct {
	// The index of the leaf the proof is for.
	LeafIndex U64
	// Number of leaves in MMR, when the proof was generated.
	LeafCount U64
	// Proof elements (hashes of siblings of inner nodes on the path to the leaf).
	Items []H256
}

// MmrProof is a MMR proof
type MmrBatchProof struct {
	// The index of the leaf the proof is for.
	LeafIndex []U64
	// Number of leaves in MMR, when the proof was generated.
	LeafCount U64
	// Proof elements (hashes of siblings of inner nodes on the path to the leaf).
	Items []H256
}

type MmrLeaf struct {
	Version               uint8
	ParentNumberAndHash   ParentNumberAndHash
	BeefyNextAuthoritySet BeefyNextAuthoritySet
	ParachainHeads        [32]byte
}

type MMRLeafVersion U8

type ParentNumberAndHash struct {
	ParentNumber uint32
	Hash         [32]byte
}

type BeefyNextAuthoritySet struct {
	// ID
	ID uint64
	// Number of validators in the set.
	Len uint32
	// Merkle Root Hash build from BEEFY uncompressed AuthorityIds.
	Root [32]byte
}
