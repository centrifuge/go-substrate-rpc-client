package ibc

import (
	chantypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
)

func (i IBC) QueryPackets(
	channelid,
	portid string,
	seqs []uint64,
) (
	[]chantypes.Packet,
	error,
) {
	var res []chantypes.Packet
	err := i.client.Call(&res, queryPacketsMethod, channelid, portid, seqs)
	if err != nil {
		return []chantypes.Packet{}, err
	}
	return res, nil
}

func (i IBC) QueryPacketCommitments(
	height uint64,
	channelid,
	portid string) (
	*chantypes.QueryPacketCommitmentsResponse,
	error,
) {
	var res *chantypes.QueryPacketCommitmentsResponse
	err := i.client.Call(&res, queryPacketCommitmentsMethod, height, channelid, portid)
	if err != nil {
		return &chantypes.QueryPacketCommitmentsResponse{}, err
	}
	return res, nil
}

func (i IBC) QueryPacketAcknowledgements(
	height uint32,
	channelid,
	portid string,
) (
	*chantypes.QueryPacketAcknowledgementsResponse,
	error,
) {
	var res *chantypes.QueryPacketAcknowledgementsResponse
	err := i.client.Call(&res, queryPacketAcknowledgementsMethod, height, channelid, portid)
	if err != nil {
		return &chantypes.QueryPacketAcknowledgementsResponse{}, err
	}
	return res, nil
}

func (i IBC) QueryUnreceivedPackets(
	height uint32,
	channelid,
	portid string,
	seqs []uint64,
) (
	[]uint64, error,
) {
	var res []uint64
	err := i.client.Call(&res, queryUnreceivedPacketsMethod, height, channelid, portid, seqs)
	if err != nil {
		return []uint64{}, err
	}
	return res, nil
}

func (i IBC) QueryUnreceivedAcknowledgements(
	height uint32,
	channelid,
	portid string,
	seqs []uint64,
) (
	[]uint64,
	error,
) {
	var res []uint64
	err := i.client.Call(&res, queryUnreceivedAcknowledgementMethod, height, channelid, portid, seqs)
	if err != nil {
		return []uint64{}, err
	}
	return res, nil
}

func (i IBC) QueryNextSeqRecv(
	height uint32,
	channelid,
	portid string,
) (
	*chantypes.QueryNextSequenceReceiveResponse,
	error,
) {
	var res *chantypes.QueryNextSequenceReceiveResponse
	err := i.client.Call(&res, queryNextSeqRecvMethod, height, channelid, portid)
	if err != nil {
		return &chantypes.QueryNextSequenceReceiveResponse{}, err
	}
	return res, nil
}

func (i IBC) QueryPacketCommitment(
	height int64,
	channelid,
	portid string,
) (
	*chantypes.QueryPacketCommitmentResponse,
	error,
) {
	var res *chantypes.QueryPacketCommitmentResponse
	err := i.client.Call(&res, queryPacketCommitmentMethod, height, channelid, portid)
	if err != nil {
		return &chantypes.QueryPacketCommitmentResponse{}, err
	}
	return res, nil
}

func (i IBC) QueryPacketAcknowledgement(
	height uint32,
	channelid,
	portid string,
	seq uint64,
) (
	*chantypes.QueryPacketAcknowledgementResponse,
	error,
) {
	var res *chantypes.QueryPacketAcknowledgementResponse
	err := i.client.Call(&res, queryPacketAcknowledgementMethod, height, channelid, portid, seq)
	if err != nil {
		return &chantypes.QueryPacketAcknowledgementResponse{}, err
	}
	return res, nil
}

func (i IBC) QueryPacketReceipt(
	height uint32,
	channelid,
	portid string,
	seq uint64,
) (
	*chantypes.QueryPacketReceiptResponse,
	error,
) {
	var res *chantypes.QueryPacketReceiptResponse
	err := i.client.Call(&res, queryPacketReceiptMethod, height, channelid, portid, seq)
	if err != nil {
		return &chantypes.QueryPacketReceiptResponse{}, err
	}
	return res, nil
}
