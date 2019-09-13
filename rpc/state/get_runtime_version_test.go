package state

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestState_GetRuntimeVersionLatest(t *testing.T) {
	runtimeVersion, err := state.GetRuntimeVersionLatest()
	assert.NoError(t, err)
	assert.Equal(t, "substrate-node", runtimeVersion.ImplName)
}

// TODO make test dynamic
//func TestState_GetRuntimeVersion(t *testing.T) {
//	bz, _ := hex.DecodeString("cc9ea640d4d4f4dd260b1cbb65cb275df995b056710265f9becdd7e6e1a7b9e0")
//
//	var bz32 [32]byte
//	copy(bz32[:], bz)
//
//	runtimeVersion, err := state.GetRuntimeVersion(types.NewHash(bz32))
//	assert.NoError(t, err)
//	assert.Equal(t, "system", runtimeVersion.RuntimeVersion.Modules[0].Name)
//}
