package state

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestState_GetMetadataLatest(t *testing.T) {
	metadata, err := state.GetMetadataLatest()
	assert.NoError(t, err)
	assert.Equal(t, "system", metadata.Metadata.Modules[0].Name)
}

// TODO make test dynamic
//func TestState_GetMetadata(t *testing.T) {
//	bz, _ := hex.DecodeString("cc9ea640d4d4f4dd260b1cbb65cb275df995b056710265f9becdd7e6e1a7b9e0")
//
//	var bz32 [32]byte
//	copy(bz32[:], bz)
//
//	metadata, err := state.GetMetadata(types.NewHash(bz32))
//	assert.NoError(t, err)
//	assert.Equal(t, "system", metadata.Metadata.Modules[0].Name)
//}
