package ibc

import (
	clienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
)

func (i IBC) QueryClientStateResponse(height int64, srcClientId string) (*clienttypes.QueryClientStateResponse, error) {
	var res *clienttypes.QueryClientStateResponse
	err := i.client.Call(&res, "ibc_queryClientState", height, srcClientId)
	if err != nil {
		return &clienttypes.QueryClientStateResponse{}, err
	}
	return res, nil
}

func (i IBC) QueryClientConsensusState(clientid string, revisionHeight, revisionNumber uint64, latestConsensusState bool) (*clienttypes.QueryConsensusStateResponse, error) {
	var res *clienttypes.QueryConsensusStateResponse
	err := i.client.Call(&res, "ibc_queryClientConsensusState", clientid, revisionHeight, revisionNumber, latestConsensusState)
	if err != nil {
		return &clienttypes.QueryConsensusStateResponse{}, err
	}
	return res, nil
}
func (i IBC) QueryUpgradedClient(height int64) (*clienttypes.QueryClientStateResponse, error) {
	var res *clienttypes.QueryClientStateResponse
	err := i.client.Call(&res, "ibc_queryUpgradedClient", height)
	if err != nil {
		return &clienttypes.QueryClientStateResponse{}, err
	}
	return res, nil
}

func (i IBC) QueryUpgradedConsState(height int64) (*clienttypes.QueryConsensusStateResponse, error) {
	var res *clienttypes.QueryConsensusStateResponse
	err := i.client.Call(&res, "ibc_queryUpgradedConnectionState", height)
	if err != nil {
		return &clienttypes.QueryConsensusStateResponse{}, err
	}
	return res, nil
}

func (i IBC) QueryClients() (*[]clienttypes.IdentifiedClientStates, error) {
	var res *[]clienttypes.IdentifiedClientStates
	err := i.client.Call(&res, "ibc_queryClients")
	if err != nil {
		return &[]clienttypes.IdentifiedClientStates{}, err
	}
	return res, nil
}
