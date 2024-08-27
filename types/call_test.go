package types_test

import (
	. "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCall(t *testing.T) {
	c := Call{CallIndex{6, 1}, Args{0, 0, 0}}

	enc, err := codec.EncodeToHex(c)
	assert.NoError(t, err)
	assert.Equal(t, "0x0601000000", enc)
}

func TestNewCallV4(t *testing.T) {
	addr, err := NewAddressFromHexAccountID("0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48")
	assert.NoError(t, err)

	c, err := NewCall(ExamplaryMetadataV4, "balances.transfer", addr, NewUCompactFromUInt(1000))
	assert.NoError(t, err)

	enc, err := codec.EncodeToHex(c)
	assert.NoError(t, err)

	assert.Equal(t, "0x0300ff8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48a10f", enc)
}

func TestNewCallV7(t *testing.T) {
	c, err := NewCall(&exampleMetadataV7, "Module2.my function", U8(3))
	assert.NoError(t, err)

	enc, err := codec.EncodeToHex(c)
	assert.NoError(t, err)

	assert.Equal(t, "0x010003", enc)
}

func TestNewCallV8(t *testing.T) {
	c, err := NewCall(&exampleMetadataV8, "Module2.my function", U8(3))
	assert.NoError(t, err)

	enc, err := codec.EncodeToHex(c)
	assert.NoError(t, err)

	assert.Equal(t, "0x010003", enc)
}
