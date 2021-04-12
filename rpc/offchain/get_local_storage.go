package offchain

import (
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
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

	kb, err := types.EncodeToHexString(key)
	if err != nil {
		return nil, fmt.Errorf("failed to encode key: %w", err)
	}

	err = c.client.Call(&res, "offchain_localStorageGet", kind, kb)
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

	kb, err := types.EncodeToHexString(key)
	if err != nil {
		return fmt.Errorf("failed to encode key: %w", err)
	}

	vb, err := types.EncodeToHexString(value)
	if err != nil {
		return fmt.Errorf("failed to encode value: %w", err)
	}

	err = c.client.Call(&res, "offchain_localStorageSet", kind, kb, vb)
	if err != nil {
		return err
	}

	return nil
}
