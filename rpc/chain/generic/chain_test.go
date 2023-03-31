package generic

import (
	"errors"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/client/mocks"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGenericChain_GetBlock(t *testing.T) {
	clientMock := mocks.NewClient(t)

	genericChain := NewChain[DefaultSignedBlock](clientMock)

	blockHash := types.Hash{1, 2, 3}

	blockJustification := []byte{4, 5, 6}

	clientMock.On("Call", mock.Anything, mock.Anything, mock.Anything).
		Run(
			func(args mock.Arguments) {
				target, ok := args.Get(0).(*DefaultSignedBlock)
				assert.True(t, ok)

				target.Justification = blockJustification

				method, ok := args.Get(1).(string)
				assert.True(t, ok)
				assert.Equal(t, getBlockMethod, method)

				hashStr, ok := args.Get(2).(string)
				assert.True(t, ok)

				encodedBlockHash := hexutil.Encode(blockHash[:])
				assert.Equal(t, encodedBlockHash, hashStr)
			},
		).Return(nil)

	res, err := genericChain.GetBlock(blockHash)
	assert.NoError(t, err)
	assert.Equal(t, blockJustification, res.Justification)
}

func TestGenericChain_GetBlock_BlockCallError(t *testing.T) {
	clientMock := mocks.NewClient(t)

	genericChain := NewChain[DefaultSignedBlock](clientMock)

	blockHash := types.Hash{1, 2, 3}

	clientMock.On("Call", mock.Anything, mock.Anything, mock.Anything).
		Run(
			func(args mock.Arguments) {
				_, ok := args.Get(0).(*DefaultSignedBlock)
				assert.True(t, ok)

				method, ok := args.Get(1).(string)
				assert.True(t, ok)
				assert.Equal(t, getBlockMethod, method)

				hashStr, ok := args.Get(2).(string)
				assert.True(t, ok)

				encodedBlockHash := hexutil.Encode(blockHash[:])
				assert.Equal(t, encodedBlockHash, hashStr)
			},
		).Return(errors.New("error"))

	res, err := genericChain.GetBlock(blockHash)
	assert.ErrorIs(t, err, ErrGetBlockCall)
	assert.Equal(t, DefaultSignedBlock{}, res)
}

func TestGenericChain_GetBlockLatest(t *testing.T) {
	clientMock := mocks.NewClient(t)

	genericChain := NewChain[DefaultSignedBlock](clientMock)

	blockJustification := []byte{4, 5, 6}

	clientMock.On("Call", mock.Anything, mock.Anything, mock.Anything).
		Run(
			func(args mock.Arguments) {
				target, ok := args.Get(0).(*DefaultSignedBlock)
				assert.True(t, ok)

				target.Justification = blockJustification

				method, ok := args.Get(1).(string)
				assert.True(t, ok)
				assert.Equal(t, getBlockMethod, method)
			},
		).Return(nil)

	res, err := genericChain.GetBlockLatest()
	assert.NoError(t, err)
	assert.Equal(t, blockJustification, res.Justification)
}

func TestGenericChain_GetBlockLatest_BlockCallError(t *testing.T) {
	clientMock := mocks.NewClient(t)

	genericChain := NewChain[DefaultSignedBlock](clientMock)

	clientMock.On("Call", mock.Anything, mock.Anything, mock.Anything).
		Run(
			func(args mock.Arguments) {
				_, ok := args.Get(0).(*DefaultSignedBlock)
				assert.True(t, ok)

				method, ok := args.Get(1).(string)
				assert.True(t, ok)
				assert.Equal(t, getBlockMethod, method)
			},
		).Return(errors.New("error"))

	res, err := genericChain.GetBlockLatest()
	assert.ErrorIs(t, err, ErrGetBlockCall)
	assert.Equal(t, DefaultSignedBlock{}, res)
}
