package system

import "github.com/mailchain/go-substrate-rpc-client/types"

func (c *System) AccountNextIndex(address string) (types.U64, error) {
	var res types.U64

	if err := c.client.Call(&res, "system_accountNextIndex", address); err != nil {
		return 0, err
	}

	return res, nil
}
