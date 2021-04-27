package mmr

import (
	"github.com/snowfork/go-substrate-rpc-client/v2/client"
	"github.com/snowfork/go-substrate-rpc-client/v2/types"
)

// GenerateProof retrieves a MMR proof and leaf for the specified leave index, at the given blockHash (useful to query a
// proof at an earlier block, likely with antoher MMR root)
func (c *MMR) GenerateProof(leafIndex uint64, blockHash types.Hash) (types.GenerateMMRProofResponse, error) {
	return c.generateProof(leafIndex, &blockHash)
}

// GenerateProofLatest retrieves the latest MMR proof and leaf for the specified leave index
func (c *MMR) GenerateProofLatest(leafIndex uint64) (types.GenerateMMRProofResponse, error) {
	return c.generateProof(leafIndex, nil)
}

func (c *MMR) generateProof(leafIndex uint64, blockHash *types.Hash) (types.GenerateMMRProofResponse, error) {
	var res types.GenerateMMRProofResponse
	err := client.CallWithBlockHash(c.client, &res, "mmr_generateProof", blockHash, leafIndex)
	if err != nil {
		return types.GenerateMMRProofResponse{}, err
	}

	return res, nil
}
