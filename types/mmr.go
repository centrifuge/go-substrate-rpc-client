package types

import (
	"encoding/json"
)

// GenerateMMRProofResponse contains the generate proof rpc response
type GenerateMMRProofResponse struct {
	BlockHash H256
	Leaf      MMRLeaf
	Proof     MMRProof
}

// UnmarshalJSON fills u with the JSON encoded byte array given by b
func (d *GenerateMMRProofResponse) UnmarshalJSON(bz []byte) error {
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

// MarshalJSON returns a JSON encoded byte array of u
// func (d GenerateMMRProofResponse) MarshalJSON() ([]byte, error) {
// 	logs := make([]string, len(d))
// 	var err error
// 	for i, di := range d {
// 		logs[i], err = EncodeToHexString(di)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	return json.Marshal(struct {
// 		Logs []string `json:"logs"`
// 	}{
// 		Logs: logs,
// 	})
// }

type MMREncodableOpaqueLeaf Bytes

// MMRProof is a MMR proof
type MMRProof struct {
	// The index of the leaf the proof is for.
	LeafIndex U64
	// Number of leaves in MMR, when the proof was generated.
	LeafCount U64
	// Proof elements (hashes of siblings of inner nodes on the path to the leaf).
	Items []H256
}

type MMRLeaf struct {
	Version               MMRLeafVersion
	ParentNumberAndHash   ParentNumberAndHash
	BeefyNextAuthoritySet BeefyNextAuthoritySet
	ParachainHeads        H256
}

type MMRLeafVersion U8

type ParentNumberAndHash struct {
	ParentNumber U32
	Hash         Hash
}

type BeefyNextAuthoritySet struct {
	// ID
	ID U64
	// Number of validators in the set.
	Len U32
	// Merkle Root Hash build from BEEFY uncompressed AuthorityIds.
	Root H256
}
