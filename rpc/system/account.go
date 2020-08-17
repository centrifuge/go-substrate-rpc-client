package system

import "github.com/mailchain/go-substrate-rpc-client/types"

func (c *System) AccountNextIndex(address string) (types.U32, error) {
	var res types.U32

	if err := c.client.Call(&res, "system_accountNextIndex", address); err != nil {
		return 0, err
	}

	return res, nil
}
