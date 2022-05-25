package mmr

import (
	"github.com/ComposableFi/go-substrate-rpc-client/v4/client"
	"github.com/ComposableFi/go-substrate-rpc-client/v4/types"
)

// GenerateProof retrieves a MMR proof and leaf for the specified leave index, at the given blockHash (useful to query a
// proof at an earlier block, likely with another MMR root)
func (c *MMR) GenerateProof(leafIndex uint64, blockHash types.Hash) (types.GenerateMmrProofResponse, error) {
	return c.generateProof(leafIndex, &blockHash)
}

// GenerateProof retrieves a MMR proof and leaf for the specified leave index, at the given blockHash (useful to query a
// proof at an earlier block, likely with another MMR root)
func (c *MMR) GenerateBatchProof(indices []uint64, blockHash types.Hash) (types.GenerateMmrBatchProofResponse, error) {
	return c.generateBatchProof(indices, &blockHash)
}

// GenerateProofLatest retrieves the latest MMR proof and leaf for the specified leave index
func (c *MMR) GenerateProofLatest(leafIndex uint64) (types.GenerateMmrProofResponse, error) {
	return c.generateProof(leafIndex, nil)
}

func (c *MMR) generateProof(leafIndex uint64, blockHash *types.Hash) (types.GenerateMmrProofResponse, error) {
	var res types.GenerateMmrProofResponse
	err := client.CallWithBlockHash(c.client, &res, "mmr_generateProof", blockHash, leafIndex)
	if err != nil {
		return types.GenerateMmrProofResponse{}, err
	}

	return res, nil
}

func (c *MMR) generateBatchProof(indices []uint64, blockHash *types.Hash) (types.GenerateMmrBatchProofResponse, error) {
	var res types.GenerateMmrBatchProofResponse
	err := client.CallWithBlockHash(c.client, &res, "mmr_generateBatchProof", blockHash, indices)
	if err != nil {
		return types.GenerateMmrBatchProofResponse{}, err
	}

	return res, nil
}
