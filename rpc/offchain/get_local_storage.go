package offchain

import (
	"fmt"

	"github.com/snowfork/go-substrate-rpc-client/v3/types"
)

// StorageKind ...
type StorageKind string

const (
	// Persistent storage
	Persistent StorageKind = "PERSISTENT"
	// Local storage
	Local StorageKind = "LOCAL"
)

// LocalStorageGet retrieves the stored data
func (c *Offchain) LocalStorageGet(kind StorageKind, key []byte) (*types.StorageDataRaw, error) {
	var res string

	err := c.client.Call(&res, "offchain_localStorageGet", kind, fmt.Sprintf("%#x", key))
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	b, err := types.HexDecodeString(res)
	if err != nil {
		return nil, err
	}

	data := types.NewStorageDataRaw(b)
	return &data, nil
}

// LocalStorageSet saves the data
func (c *Offchain) LocalStorageSet(kind StorageKind, key []byte, value []byte) error {
	var res string

	err := c.client.Call(&res, "offchain_localStorageSet", kind, fmt.Sprintf("%#x", key), fmt.Sprintf("%#x", value))
	if err != nil {
		return err
	}

	return nil
}
