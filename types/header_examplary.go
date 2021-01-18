// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
//
// Copyright 2019 Centrifuge GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build test

package types

var ExamplaryHeader = Header{ParentHash: Hash{0xf6, 0xe0, 0xf9, 0xd, 0x27, 0x57, 0x8a, 0x18, 0xfc, 0x22, 0xcf, 0x10, 0x83, 0x3d, 0x7f, 0xb2, 0x86, 0xf9, 0x78, 0xce, 0x4e, 0x41, 0x1e, 0x87, 0xb2, 0x19, 0xc1, 0xd4, 0x98, 0xf6, 0xd3, 0x2c}, Number: 0xa, StateRoot: Hash{0x2e, 0x5e, 0x1e, 0xa0, 0x5f, 0xcd, 0xec, 0xca, 0xa, 0xa, 0xf0, 0xe9, 0x8d, 0xc5, 0xc3, 0x20, 0xcb, 0x62, 0x13, 0xad, 0xc2, 0x1f, 0x4f, 0xad, 0x2c, 0xf, 0xfc, 0x74, 0x7a, 0x18, 0x64, 0xc}, ExtrinsicsRoot: Hash{0xc4, 0x5f, 0x66, 0x5d, 0x47, 0xe4, 0x8e, 0x54, 0xf, 0x1c, 0x89, 0xd8, 0x7e, 0x8, 0x79, 0xb8, 0x0, 0x53, 0x34, 0xf0, 0x6e, 0x53, 0x7d, 0xaa, 0x7c, 0xe8, 0xab, 0x51, 0xcc, 0x12, 0x39, 0x1b}, Digest: Digest{DigestItem{IsOther: false, AsOther: Bytes(nil), IsAuthoritiesChange: false, AsAuthoritiesChange: []AuthorityID(nil), IsChangesTrieRoot: false, AsChangesTrieRoot: Hash{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, IsSealV0: false, AsSealV0: SealV0{Signer: 0x0, Signature: Signature{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}}, IsConsensus: false, AsConsensus: Consensus{ConsensusEngineID: 0x0, Bytes: Bytes(nil)}, IsSeal: false, AsSeal: Seal{ConsensusEngineID: 0x0, Bytes: Bytes(nil)}, IsPreRuntime: true, AsPreRuntime: PreRuntime{ConsensusEngineID: 0x45424142, Bytes: Bytes{0x2, 0x0, 0x0, 0x0, 0x0, 0xd9, 0xe8, 0x33, 0x1f, 0x0, 0x0, 0x0, 0x0}}}, DigestItem{IsOther: false, AsOther: Bytes(nil), IsAuthoritiesChange: false, AsAuthoritiesChange: []AuthorityID(nil), IsChangesTrieRoot: false, AsChangesTrieRoot: Hash{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, IsSealV0: false, AsSealV0: SealV0{Signer: 0x0, Signature: Signature{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}}, IsConsensus: false, AsConsensus: Consensus{ConsensusEngineID: 0x0, Bytes: Bytes(nil)}, IsSeal: true, AsSeal: Seal{ConsensusEngineID: 0x45424142, Bytes: Bytes{0xaa, 0x3f, 0x2f, 0x1, 0x50, 0x15, 0x4a, 0x9a, 0x6c, 0xe6, 0xb8, 0xf5, 0x18, 0x28, 0x8b, 0xc2, 0x92, 0xb2, 0x1, 0xf8, 0x36, 0x32, 0xdc, 0xf9, 0xeb, 0xd, 0x2d, 0x0, 0x5a, 0x38, 0xde, 0x2a, 0xa3, 0x67, 0x45, 0x31, 0xc7, 0x4, 0x46, 0x4e, 0xe6, 0x76, 0x88, 0x76, 0x83, 0x68, 0xba, 0xb8, 0x40, 0x11, 0x7, 0x6a, 0x35, 0xe8, 0xe6, 0xdd, 0x4a, 0xf5, 0x9d, 0xb4, 0x15, 0x5, 0x2f, 0x8d}}, IsPreRuntime: false, AsPreRuntime: PreRuntime{ConsensusEngineID: 0x0, Bytes: Bytes(nil)}}}} //nolint:lll
