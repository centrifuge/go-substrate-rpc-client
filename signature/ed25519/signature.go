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

package ed25519

import (
	"errors"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/vedhavyas/go-subkey/v2"
	"github.com/vedhavyas/go-subkey/v2/ed25519"
	"golang.org/x/crypto/blake2b"
)

// signature.KeyringPairFromSecret creates KeyPair based on seed/phrase and network
// Leave network empty for default behavior
func KeyringPairFromSecret(seedOrPhrase string, network uint16) (signature.KeyringPair, error) {
	scheme := ed25519.Scheme{}
	kyr, err := subkey.DeriveKeyPair(scheme, seedOrPhrase)
	if err != nil {
		return signature.KeyringPair{}, err
	}

	ss58Address := kyr.SS58Address(network)

	var pk = kyr.Public()

	return signature.KeyringPair{
		KeyType:   1,
		URI:       seedOrPhrase,
		Address:   ss58Address,
		PublicKey: pk,
	}, nil
}

// Sign signs data with the private key under the given derivation path, returning the signature. Requires the subkey
// command to be in path
func Sign(data []byte, privateKeyURI string) ([]byte, error) {
	// if data is longer than 256 bytes, hash it first
	if len(data) > 256 {
		h := blake2b.Sum256(data)
		data = h[:]
	}

	scheme := ed25519.Scheme{}
	kyr, err := subkey.DeriveKeyPair(scheme, privateKeyURI)
	if err != nil {
		return nil, err
	}

	signature, err := kyr.Sign(data)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

// Verify verifies data using the provided signature and the key under the derivation path. Requires the subkey
// command to be in path
func Verify(data []byte, sig []byte, privateKeyURI string) (bool, error) {
	// if data is longer than 256 bytes, hash it first
	if len(data) > 256 {
		h := blake2b.Sum256(data)
		data = h[:]
	}

	scheme := ed25519.Scheme{}
	kyr, err := subkey.DeriveKeyPair(scheme, privateKeyURI)
	if err != nil {
		return false, err
	}

	if len(sig) != 64 {
		return false, errors.New("wrong signature length")
	}

	v := kyr.Verify(data, sig)

	return v, nil
}
