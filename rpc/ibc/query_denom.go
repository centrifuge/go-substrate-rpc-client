package ibc

import (
	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
)

func (i IBC) QueryDenomTrace(
	denom string,
) (
	transfertypes.QueryDenomTraceResponse,
	error,
) {
	var res transfertypes.QueryDenomTraceResponse
	err := i.client.Call(&res, queryDenomTraceMethod, denom)
	if err != nil {
		return transfertypes.QueryDenomTraceResponse{}, err
	}
	return res, nil
}

func (i IBC) QueryDenomTraces(
	offset,
	limit uint64,
	height uint32,
) (
	*transfertypes.QueryDenomTracesResponse,
	error,
) {
	var res *transfertypes.QueryDenomTracesResponse
	err := i.client.Call(&res, queryDenomTracesMethod, offset, limit, height)
	if err != nil {
		return &transfertypes.QueryDenomTracesResponse{}, err
	}
	return res, nil
}
