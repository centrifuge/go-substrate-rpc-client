package ibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (i IBC) QueryBalanceWithAddress(addr []byte) (*sdk.Coin, error) {
	var res *sdk.Coin
	err := i.client.Call(&res, "ibc_queryBalanceWithAddress", addr)
	if err != nil {
		return &sdk.Coin{}, err
	}
	return res, nil
}
