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
	"encoding/hex"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// const subkeyCmd = "/Users/philipstanislaus/go/src/github.com/paritytech/substrate/target/debug/subkey"
const subkeyCmd = "subkey"

type KeyringPair struct {
	// URI is the derivation path for the private key in subkey
	URI string
	// Address is an SS58 address
	Address string
	// PublicKey
	PublicKey []byte
}

var TestKeyringPairAlice = KeyringPair{
	URI:       "//Alice",
	PublicKey: []byte{0xd4, 0x35, 0x93, 0xc7, 0x15, 0xfd, 0xd3, 0x1c, 0x61, 0x14, 0x1a, 0xbd, 0x4, 0xa9, 0x9f, 0xd6, 0x82, 0x2c, 0x85, 0x58, 0x85, 0x4c, 0xcd, 0xe3, 0x9a, 0x56, 0x84, 0xe7, 0xa5, 0x6d, 0xa2, 0x7d}, //nolint:lll
	Address:   "5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY",
}

// Sign signs data with the private key under the given derivation path, returning the signature. Requires the subkey
// command to be in path
func Sign(data []byte, privateKeyURI string) ([]byte, error) {
	// use "subkey" command for signature
	cmd := exec.Command(subkeyCmd, "sign", privateKeyURI, "--hex")

	// data to stdin
	dataHex := hex.EncodeToString(data)
	cmd.Stdin = strings.NewReader(dataHex)

	log.Printf("echo -n \"%v\" | %v sign %v --hex", dataHex, subkeyCmd, privateKeyURI)

	// execute the command, get the output
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to sign with subkey: %v", err.Error())
	}

	// remove line feed
	if len(out) > 0 && out[len(out)-1] == 10 {
		out = out[:len(out)-1]
	}

	outStr := string(out)

	dec, err := hex.DecodeString(outStr)

	return dec, err
}

// Verify verifies data using the provided signature and the key under the derivation path. Requires the subkey
// command to be in path
func Verify(data []byte, sig []byte, privateKeyURI string) (bool, error) {
	// hexify the sig
	sigHex := hex.EncodeToString(sig)

	// use "subkey" command for signature
	cmd := exec.Command(subkeyCmd, "verify", sigHex, privateKeyURI, "--hex")

	// data to stdin
	dataHex := hex.EncodeToString(data)
	cmd.Stdin = strings.NewReader(dataHex)

	log.Printf("echo -n \"%v\" | %v verify %v %v --hex", dataHex, subkeyCmd, sigHex, privateKeyURI)

	// execute the command, get the output
	out, err := cmd.Output()
	if err != nil {
		log.Fatal("Failed to verify with subkey", err.Error())
	}

	// remove line feed
	if len(out) > 0 && out[len(out)-1] == 10 {
		out = out[:len(out)-1]
	}

	outStr := string(out)
	valid := outStr == "Signature verifies correctly."
	return valid, nil
}
