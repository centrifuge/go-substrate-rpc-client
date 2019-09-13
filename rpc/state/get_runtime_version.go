package state

import (
	"github.com/centrifuge/go-substrate-rpc-client/types"
)

func (s *State) GetRuntimeVersion(blockHash types.Hash) (*types.RuntimeVersion, error) {
	return s.getRuntimeVersion(&blockHash)
}

func (s *State) GetRuntimeVersionLatest() (*types.RuntimeVersion, error) {
	return s.getRuntimeVersion(nil)
}

func (s *State) getRuntimeVersion(blockHash *types.Hash) (*types.RuntimeVersion, error) {
	var runtimeVersion types.RuntimeVersion
	var err error
	if blockHash == nil {
		err = (*s.client).Call(&runtimeVersion, "state_getRuntimeVersion")
	} else {
		err = (*s.client).Call(&runtimeVersion, "state_getRuntimeVersion", *blockHash)
	}
	if err != nil {
		return &runtimeVersion, err
	}
	return &runtimeVersion, err
}

//func (s *State) getRuntimeVersion(blockHash *types.Hash) (*types.RuntimeVersion, error) {
//	runtimeVersion := types.NewRuntimeVersion()
//
//	var res string
//	var err error
//	if blockHash == nil {
//		err = (*s.client).Call(&res, "state_getRuntimeVersion")
//	} else {
//		err = (*s.client).Call(&res, "state_getRuntimeVersion", *blockHash)
//	}
//	if err != nil {
//		return runtimeVersion, err
//	}
//
//	err = types.DecodeFromHexString(res, runtimeVersion)
//	return runtimeVersion, err
//}
