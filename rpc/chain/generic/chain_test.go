package generic

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/client/mocks"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGenericDefaultChain_GetBlock(t *testing.T) {
	clientMock := mocks.NewClient(t)

	genericChain := NewDefaultChain(clientMock)

	blockHash := types.Hash{1, 2, 3}

	blockJustification := []byte{4, 5, 6}

	clientMock.On("CallContext", context.Background(), mock.Anything, mock.Anything, mock.Anything).
		Run(
			func(args mock.Arguments) {
				target, ok := args.Get(1).(**DefaultGenericSignedBlock)
				assert.True(t, ok)

				*target = new(DefaultGenericSignedBlock)
				(**target).Justification = blockJustification

				method, ok := args.Get(2).(string)
				assert.True(t, ok)
				assert.Equal(t, getBlockMethod, method)

				hashStr, ok := args.Get(3).(string)
				assert.True(t, ok)

				encodedBlockHash := hexutil.Encode(blockHash[:])
				assert.Equal(t, encodedBlockHash, hashStr)
			},
		).Return(nil)

	res, err := genericChain.GetBlock(blockHash)
	assert.NoError(t, err)
	assert.Equal(t, blockJustification, res.GetJustification())
}

func TestGenericDefaultChain_GetBlock_BlockCallError(t *testing.T) {
	clientMock := mocks.NewClient(t)

	genericChain := NewDefaultChain(clientMock)

	blockHash := types.Hash{1, 2, 3}

	clientMock.On("CallContext", context.Background(), mock.Anything, mock.Anything, mock.Anything).
		Run(
			func(args mock.Arguments) {
				_, ok := args.Get(1).(**DefaultGenericSignedBlock)
				assert.True(t, ok)

				method, ok := args.Get(2).(string)
				assert.True(t, ok)
				assert.Equal(t, getBlockMethod, method)

				hashStr, ok := args.Get(3).(string)
				assert.True(t, ok)

				encodedBlockHash := hexutil.Encode(blockHash[:])
				assert.Equal(t, encodedBlockHash, hashStr)
			},
		).Return(errors.New("error"))

	res, err := genericChain.GetBlock(blockHash)
	assert.ErrorIs(t, err, ErrGetBlockCall)
	assert.Nil(t, res)
}

func TestGenericDefaultChain_GetBlockLatest(t *testing.T) {
	clientMock := mocks.NewClient(t)

	genericChain := NewDefaultChain(clientMock)

	blockJustification := []byte{4, 5, 6}

	clientMock.On("CallContext", context.Background(), mock.Anything, mock.Anything, mock.Anything).
		Run(
			func(args mock.Arguments) {
				target, ok := args.Get(1).(**DefaultGenericSignedBlock)
				assert.True(t, ok)

				*target = new(DefaultGenericSignedBlock)
				(**target).Justification = blockJustification

				method, ok := args.Get(2).(string)
				assert.True(t, ok)
				assert.Equal(t, getBlockMethod, method)
			},
		).Return(nil)

	res, err := genericChain.GetBlockLatest()
	assert.NoError(t, err)
	assert.Equal(t, blockJustification, res.Justification)
}

func TestGenericDefaultChain_GetBlockLatest_BlockCallError(t *testing.T) {
	clientMock := mocks.NewClient(t)

	genericChain := NewDefaultChain(clientMock)

	clientMock.On("CallContext", context.Background(), mock.Anything, mock.Anything, mock.Anything).
		Run(
			func(args mock.Arguments) {
				_, ok := args.Get(1).(**DefaultGenericSignedBlock)
				assert.True(t, ok)

				method, ok := args.Get(2).(string)
				assert.True(t, ok)
				assert.Equal(t, getBlockMethod, method)
			},
		).Return(errors.New("error"))

	res, err := genericChain.GetBlockLatest()
	assert.ErrorIs(t, err, ErrGetBlockCall)
	assert.Nil(t, res)
}

type testRunner interface {
	runTest(t *testing.T)
}

type testBlockData[
	A any,
	S any,
	P any,
] struct{}

func (td *testBlockData[A, S, P]) getGenericTestTypeNames() string {
	var (
		a A
		s S
		p P
	)

	return fmt.Sprintf("%T - %T - %T", a, s, p)
}

func (td *testBlockData[A, S, P]) getTestName(testName string) string {
	return fmt.Sprintf("%s with: %s", testName, td.getGenericTestTypeNames())
}

func (td *testBlockData[A, S, P]) runTest(t *testing.T) {
	t.Run(td.getTestName("Get block"), func(t *testing.T) {
		clientMock := mocks.NewClient(t)

		genericChain := NewChain[A, S, P, *SignedBlock[A, S, P]](clientMock)

		blockHash := types.Hash{1, 2, 3}

		blockJustification := []byte{4, 5, 6}

		clientMock.On("CallContext", context.Background(), mock.Anything, mock.Anything, mock.Anything).
			Run(
				func(args mock.Arguments) {
					target, ok := args.Get(1).(**SignedBlock[A, S, P])
					assert.True(t, ok)

					*target = new(SignedBlock[A, S, P])
					(*target).Justification = blockJustification

					method, ok := args.Get(2).(string)
					assert.True(t, ok)
					assert.Equal(t, getBlockMethod, method)

					hashStr, ok := args.Get(3).(string)
					assert.True(t, ok)

					encodedBlockHash := hexutil.Encode(blockHash[:])
					assert.Equal(t, encodedBlockHash, hashStr)
				},
			).Return(nil)

		res, err := genericChain.GetBlock(blockHash)
		assert.NoError(t, err)
		assert.Equal(t, blockJustification, res.GetJustification())
	})

	t.Run(td.getTestName("Get block - block call error"), func(t *testing.T) {
		clientMock := mocks.NewClient(t)

		genericChain := NewChain[A, S, P, *SignedBlock[A, S, P]](clientMock)

		blockHash := types.Hash{1, 2, 3}

		clientMock.On("CallContext", context.Background(), mock.Anything, mock.Anything, mock.Anything).
			Run(
				func(args mock.Arguments) {
					_, ok := args.Get(1).(**SignedBlock[A, S, P])
					assert.True(t, ok)

					method, ok := args.Get(2).(string)
					assert.True(t, ok)
					assert.Equal(t, getBlockMethod, method)

					hashStr, ok := args.Get(3).(string)
					assert.True(t, ok)

					encodedBlockHash := hexutil.Encode(blockHash[:])
					assert.Equal(t, encodedBlockHash, hashStr)
				},
			).Return(errors.New("error"))

		res, err := genericChain.GetBlock(blockHash)
		assert.ErrorIs(t, err, ErrGetBlockCall)
		assert.Nil(t, res)
	})

	t.Run(td.getTestName("Get latest block"), func(t *testing.T) {
		clientMock := mocks.NewClient(t)

		genericChain := NewChain[A, S, P, *SignedBlock[A, S, P]](clientMock)

		blockJustification := []byte{4, 5, 6}

		clientMock.On("CallContext", context.Background(), mock.Anything, mock.Anything, mock.Anything).
			Run(
				func(args mock.Arguments) {
					target, ok := args.Get(1).(**SignedBlock[A, S, P])
					assert.True(t, ok)

					*target = new(SignedBlock[A, S, P])
					(*target).Justification = blockJustification

					method, ok := args.Get(2).(string)
					assert.True(t, ok)
					assert.Equal(t, getBlockMethod, method)
				},
			).Return(nil)

		res, err := genericChain.GetBlockLatest()
		assert.NoError(t, err)
		assert.Equal(t, blockJustification, res.GetJustification())
	})

	t.Run(td.getTestName("Get latest block - block call error"), func(t *testing.T) {
		clientMock := mocks.NewClient(t)

		genericChain := NewChain[A, S, P, *SignedBlock[A, S, P]](clientMock)

		clientMock.On("CallContext", context.Background(), mock.Anything, mock.Anything, mock.Anything).
			Run(
				func(args mock.Arguments) {
					_, ok := args.Get(1).(**SignedBlock[A, S, P])
					assert.True(t, ok)

					method, ok := args.Get(2).(string)
					assert.True(t, ok)
					assert.Equal(t, getBlockMethod, method)
				},
			).Return(errors.New("error"))

		res, err := genericChain.GetBlockLatest()
		assert.ErrorIs(t, err, ErrGetBlockCall)
		assert.Nil(t, res)
	})
}

func TestGenericChain(t *testing.T) {
	testRunners := []testRunner{
		&testBlockData[
			types.MultiAddress,
			types.MultiSignature,
			DefaultPaymentFields,
		]{},
		&testBlockData[
			types.MultiAddress,
			types.MultiSignature,
			PaymentFieldsWithAssetID,
		]{},
		&testBlockData[
			[20]byte,
			[65]byte,
			DefaultPaymentFields,
		]{},
	}

	for _, test := range testRunners {
		test.runTest(t)
	}
}
