package ibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (i IBC) QueryBalanceWithAddress(
	addr []byte,
) (
	sdk.Coins,
	error,
) {
	var res sdk.Coins
	err := i.client.Call(&res, "ibc_queryBalanceWithAddress", addr)
	if err != nil {
		return sdk.Coins{}, err
	}
	return res, nil
}
