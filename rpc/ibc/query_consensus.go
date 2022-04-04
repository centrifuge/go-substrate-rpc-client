package ibc

import (
	clienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	ibcexported "github.com/cosmos/ibc-go/v3/modules/core/exported"
)

func (i IBC) QueryConsensusState(height ibcexported.Height) (clienttypes.QueryConsensusStateResponse, error) {
	var res clienttypes.QueryConsensusStateResponse
	err := i.client.Call(&res, "ibc_queryConsensusState", height)
	if err != nil {
		return clienttypes.QueryConsensusStateResponse{}, err
	}
	return res, nil
}
