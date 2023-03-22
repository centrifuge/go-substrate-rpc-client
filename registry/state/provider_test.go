package state

import (
	"errors"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"

	"github.com/centrifuge/go-substrate-rpc-client/v4/rpc/state/mocks"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

func TestProvider_GetLatestMetadata(t *testing.T) {
	stateRPCMock := mocks.NewState(t)

	provider := NewProvider(stateRPCMock)

	testMeta := &types.Metadata{}

	stateRPCMock.On("GetMetadataLatest").
		Return(testMeta, nil).
		Once()

	res, err := provider.GetLatestMetadata()
	assert.NoError(t, err)
	assert.Equal(t, testMeta, res)

	stateRPCError := errors.New("error")

	stateRPCMock.On("GetMetadataLatest").
		Return(nil, stateRPCError).
		Once()

	res, err = provider.GetLatestMetadata()
	assert.ErrorIs(t, err, stateRPCError)
	assert.Nil(t, res)
}

func TestProvider_GetMetadata(t *testing.T) {
	stateRPCMock := mocks.NewState(t)

	provider := NewProvider(stateRPCMock)

	testHash := types.Hash{}
	testMeta := &types.Metadata{}

	stateRPCMock.On("GetMetadata", testHash).
		Return(testMeta, nil).
		Once()

	res, err := provider.GetMetadata(testHash)
	assert.NoError(t, err)
	assert.Equal(t, testMeta, res)

	stateRPCError := errors.New("error")

	stateRPCMock.On("GetMetadata", testHash).
		Return(nil, stateRPCError).
		Once()

	res, err = provider.GetMetadata(testHash)
	assert.ErrorIs(t, err, stateRPCError)
	assert.Nil(t, res)

	types.NewMetadataV14()
}

func TestProvider_GetStorageEvents(t *testing.T) {
	stateRPCMock := mocks.NewState(t)

	provider := NewProvider(stateRPCMock)

	testHash := types.Hash{}

	var testMeta types.Metadata

	err := codec.DecodeFromHex(types.MetadataV14Data, &testMeta)
	assert.NoError(t, err)

	storageKey, err := types.CreateStorageKey(&testMeta, storagePrefix, storageMethod, nil)
	assert.NoError(t, err)

	storageData := &types.StorageDataRaw{}

	stateRPCMock.On("GetStorageRaw", storageKey, testHash).
		Return(storageData, nil).
		Once()

	res, err := provider.GetStorageEvents(&testMeta, testHash)
	assert.NoError(t, err)
	assert.Equal(t, storageData, res)

	stateRPCError := errors.New("error")

	stateRPCMock.On("GetStorageRaw", storageKey, testHash).
		Return(nil, stateRPCError).
		Once()

	res, err = provider.GetStorageEvents(&testMeta, testHash)
	assert.ErrorIs(t, err, stateRPCError)
	assert.Nil(t, res)
}
