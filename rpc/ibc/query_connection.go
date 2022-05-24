package ibc

import (
	conntypes "github.com/cosmos/ibc-go/v3/modules/core/03-connection/types"
)

func (i IBC) QueryConnection(
	height int64,
	connectionID string,
) (
	*conntypes.QueryConnectionResponse,
	error,
) {
	var res *conntypes.QueryConnectionResponse
	err := i.client.Call(&res, queryConnectionMethod, height, connectionID)
	if err != nil {
		return &conntypes.QueryConnectionResponse{}, err
	}
	return res, nil
}

func (i IBC) QueryConnections() (
	*conntypes.QueryConnectionsResponse,
	error,
) {
	var res *conntypes.QueryConnectionsResponse
	err := i.client.Call(&res, queryConnectionsMethod)
	if err != nil {
		return &conntypes.QueryConnectionsResponse{}, err
	}
	return res, nil
}

func (i IBC) QueryConnectionsUsingClient(
	height int64,
	clientid string,
) (
	*conntypes.QueryConnectionsResponse,
	error,
) {
	var res *conntypes.QueryConnectionsResponse
	err := i.client.Call(&res, queryConnectionUsingClientMethod, height, clientid)
	if err != nil {
		return &conntypes.QueryConnectionsResponse{}, err
	}
	return res, nil
}
