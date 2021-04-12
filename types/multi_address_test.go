package types_test

import (
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v3/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
	"github.com/stretchr/testify/assert"
)

func TestNewMultiAddressFromAccountID(t *testing.T) {
	assertRoundtrip(t, types.NewMultiAddressFromAccountID(signature.TestKeyringPairAlice.PublicKey))

	_, err := types.NewMultiAddressFromHexAccountID("123!")
	assert.Error(t, err)

	addr, err := types.NewMultiAddressFromHexAccountID(types.HexEncodeToString(signature.TestKeyringPairAlice.PublicKey))
	assert.NoError(t, err)
	assertRoundtrip(t, addr)
	assertRoundtrip(t, types.MultiAddress{
		IsIndex: true,
		AsIndex: 100,
	})
	assertRoundtrip(t, types.MultiAddress{
		IsRaw: true,
		AsRaw: []byte{1, 2, 3},
	})
	assertRoundtrip(t, types.MultiAddress{
		IsAddress32: true,
		AsAddress32: [32]byte{},
	})
	assertRoundtrip(t, types.MultiAddress{
		IsAddress20: true,
		AsAddress20: [20]byte{},
	})
}
