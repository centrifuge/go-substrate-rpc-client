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

package signature

import (
	"errors"
	"fmt"
	"os"

	"github.com/vedhavyas/go-subkey/common"
	subkey "github.com/vedhavyas/go-subkey/sr25519"
	"golang.org/x/crypto/blake2b"
)

type KeyringPair struct {
	// URI is the derivation path for the private key in subkey
	URI string
	// Address is an SS58 address
	Address string
	// PublicKey
	PublicKey []byte
}

// KeyringPairFromSecret creates KeyPair based on seed/phrase and network
// Leave network empty for default behavior
func KeyringPairFromSecret(seedOrPhrase, network string) (KeyringPair, error) {
	if network == "" {
		network = "substrate"
	}

	kyr, err := subkey.KeyRingFromURI(seedOrPhrase)
	if err != nil {
		return KeyringPair{}, err
	}

	ss58Address, err := kyr.SS58Address(network, common.SS58Checksum)
	if err != nil {
		return KeyringPair{}, err
	}

	var pk = kyr.Public()

	return KeyringPair{
		URI:       seedOrPhrase,
		Address:   ss58Address,
		PublicKey: pk[:],
	}, nil
}

var TestKeyringPairAlice = KeyringPair{
	URI:       "//Alice",
	PublicKey: []byte{0xd4, 0x35, 0x93, 0xc7, 0x15, 0xfd, 0xd3, 0x1c, 0x61, 0x14, 0x1a, 0xbd, 0x4, 0xa9, 0x9f, 0xd6, 0x82, 0x2c, 0x85, 0x58, 0x85, 0x4c, 0xcd, 0xe3, 0x9a, 0x56, 0x84, 0xe7, 0xa5, 0x6d, 0xa2, 0x7d}, //nolint:lll
	Address:   "5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY",
}

// Sign signs data with the private key under the given derivation path, returning the signature. Requires the subkey
// command to be in path
func Sign(data []byte, privateKeyURI string) ([]byte, error) {
	// if data is longer than 256 bytes, hash it first
	if len(data) > 256 {
		h := blake2b.Sum256(data)
		data = h[:]
	}

	kyr, err := subkey.KeyRingFromURI(privateKeyURI)
	if err != nil {
		return nil, err
	}

	signature, err := kyr.Sign(data)
	if err != nil {
		return nil, err
	}

	return signature[:], nil
}

// Verify verifies data using the provided signature and the key under the derivation path. Requires the subkey
// command to be in path
func Verify(data []byte, sig []byte, privateKeyURI string) (bool, error) {
	// if data is longer than 256 bytes, hash it first
	if len(data) > 256 {
		h := blake2b.Sum256(data)
		data = h[:]
	}

	kyr, err := subkey.KeyRingFromURI(privateKeyURI)
	if err != nil {
		return false, err
	}

	if len(sig) != 64 {
		return false, errors.New("wrong signature length")
	}
	var sig64 [64]byte
	copy(sig64[:], sig)

	v := kyr.Verify(data, sig64)

	return v, nil
}

// LoadKeyringPairFromEnv looks up whether the env variable TEST_PRIV_KEY is set and is not empty and tries to use its
// content as a private phrase, seed or URI to derive a key ring pair. Panics if the private phrase, seed or URI is
// not valid or the keyring pair cannot be derived
// Loads Network from TEST_NETWORK variable
// Leave TEST_NETWORK empty or unset for default
func LoadKeyringPairFromEnv() (kp KeyringPair, ok bool) {
	network := os.Getenv("TEST_NETWORK")
	priv, ok := os.LookupEnv("TEST_PRIV_KEY")
	if !ok || priv == "" {
		return kp, false
	}
	kp, err := KeyringPairFromSecret(priv, network)
	if err != nil {
		panic(fmt.Errorf("cannot load keyring pair from env or use fallback: %v", err))
	}
	return kp, true
}
