# Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
#
# Copyright 2019 Centrifuge GmbH
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

clean: 				##clean vendor folder. Should be run before a make install
	@echo 'cleaning previous /vendor folder'
	@rm -rf vendor/
	@echo 'done cleaning'

install: 			## installs dependencies
	@command -v dep >/dev/null 2>&1 || curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	@command -v golangci-lint >/dev/null 2>&1 || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.25.0
	@dep ensure

lint: 				## runs linters on go code
	@golangci-lint run

lint-fix: 			## runs linters on go code and automatically fixes issues
	@golangci-lint run --fix

test: 				## runs all tests in project against the RPC URL specified in the RPC_URL env variable or localhost while excluding gethrpc
	@go test -race -count=1 `go list ./... | grep -v '/gethrpc'`

test-cover: 			## runs all tests in project against the RPC URL specified in the RPC_URL env variable or localhost and report coverage
	@go test -race -coverprofile=coverage.txt -covermode=atomic `go list ./... | grep -v '/gethrpc'`
	@mv coverage.txt shared

test-dockerized: 		## runs all tests in a docker container against the Substrate Default Docker image
	@docker-compose build
	@docker-compose up --abort-on-container-exit

test-e2e-deployed: export RPC_URL?=wss://serinus-5.kusama.network
test-e2e-deployed: 		## runs only end-to-end (e2e) tests against a deployed testnet (defaults to Kusama CC2 (wss://serinus-5.kusama.network) if RPC_URL is not set)
	@docker build . -t gsrpc-test
	@docker run --rm -e RPC_URL -e TEST_PRIV_KEY gsrpc-test go test -v github.com/centrifuge/go-substrate-rpc-client/teste2e

run-substrate-docker: 		## runs the Substrate 1.0 Default Docker image, this can be used to run the tests
	docker run -p 9933:9933 -p 9944:9944 -p 30333:30333 parity/substrate:latest-v1.0 --dev --rpc-external --ws-external

run-substrate-docker-v2: 	## runs the Substrate 2.0 Default Docker image, this can be used to run the tests
	docker run -p 9933:9933 -p 9944:9944 -p 30333:30333 parity/substrate:v2.0.0-rc6 --dev --rpc-external --ws-external

help: 				## shows this help
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST)

.PHONY: clean install lint lint-fix test test-dockerized run-substrate-docker run-substrate-docker-v2
