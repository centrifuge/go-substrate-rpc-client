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

const (
	Alice = "//Alice"
)

// Sign signs data, returning the signature. Requires the subkey command to be in path
func Sign(data []byte) ([]byte, error) {
	// use "subkey" command for signature
	cmd := exec.Command("subkey", "sign", Alice)

	// data to stdin
	dataHex := hex.EncodeToString(data)
	cmd.Stdin = strings.NewReader(dataHex)

	log.Printf("echo \"%v\" | subkey sign %v", dataHex, Alice)

	// execute the command, get the output
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("Failed to sign with subkey: %v", err.Error())
	}

	// remove line feed
	if len(out) > 0 && out[len(out)-1] == 10 {
		out = out[:len(out)-1]
	}

	outStr := string(out)

	dec, err := hex.DecodeString(outStr)

	return dec, err
}

// Verify verifies data using the provided signature. Requires the subkey command to be in path
func Verify(data []byte, sig []byte) (bool, error) {
	// hexify the sig
	sigHex := hex.EncodeToString(sig)

	// use "subkey" command for signature
	cmd := exec.Command("subkey", "verify", sigHex, Alice)

	// data to stdin
	dataHex := hex.EncodeToString(data)
	cmd.Stdin = strings.NewReader(dataHex)

	log.Printf("echo \"%v\" | subkey verify %v %v", dataHex, sigHex, Alice)

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
