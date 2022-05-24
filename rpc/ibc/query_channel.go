package ibc

import (
	clienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	chantypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
)

func (i IBC) QueryChannelClient(
	height uint32,
	channelid,
	portid string,
) (
	*clienttypes.IdentifiedClientState,
	error,
) {
	var res clienttypes.IdentifiedClientState
	err := i.client.Call(&res, queryChannelClientMethod, height, channelid, portid)
	if err != nil {
		return &clienttypes.IdentifiedClientState{}, err
	}
	return &res, nil
}

func (i IBC) QueryConnectionChannels(
	height uint32,
	connectionid string,
) (
	*chantypes.QueryChannelsResponse,
	error,
) {
	var res *chantypes.QueryChannelsResponse
	err := i.client.Call(&res, queryConnectionChannelsMethod, height, connectionid)
	if err != nil {
		return &chantypes.QueryChannelsResponse{}, err
	}
	return res, nil
}

func (i IBC) QueryChannels() (
	*chantypes.QueryChannelsResponse,
	error,
) {
	var res *chantypes.QueryChannelsResponse
	err := i.client.Call(&res, queryChannelsMethod)
	if err != nil {
		return &chantypes.QueryChannelsResponse{}, err
	}
	return res, nil
}
