package types_test

import (
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v2/types"
	"github.com/stretchr/testify/assert"
)

func TestBeefySignature(t *testing.T) {
	empty := types.NewOptionBeefySignatureEmpty()
	assert.True(t, empty.IsNone())
	assert.False(t, empty.IsSome())

	sig := types.NewOptionBeefySignature(types.BeefySignature{})
	sig.SetNone()
	assert.True(t, sig.IsNone())
	sig.SetSome(types.BeefySignature{})
	assert.True(t, sig.IsSome())
	assertRoundtrip(t, sig)
}
