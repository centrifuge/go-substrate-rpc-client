# Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
# Copyright (C) 2019  Centrifuge GmbH
#
# This file is part of Go Substrate RPC Client (GSRPC).
#
# GSRPC is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# GSRPC is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.

clean: 				##clean vendor folder. Should be run before a make install
	@echo 'cleaning previous /vendor folder'
	@rm -rf vendor/
	@echo 'done cleaning'

install: 			## installs dependencies
	@command -v dep >/dev/null 2>&1 || go get -u github.com/golang/dep/...
	@command -v golangci-lint >/dev/null 2>&1 || go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	@dep ensure

lint: 				## runs linters on go code
	@golangci-lint run

lint-fix: 			## runs linters on go code and automatically fixes issues
	@golangci-lint run --fix

test: 				## runs all tests in project against the RPC URL specified in the RPC_URL env variable or localhost
	@go test -v -race ./...

test-cover: 			## runs all tests in project against the RPC URL specified in the RPC_URL env variable
                                ## or localhost and report coverage
	@go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	@mv coverage.txt shared

test-dockerized: 		## runs all tests in a docker container against the Substrate Default Docker image
	@docker-compose build
	@docker-compose up --abort-on-container-exit

test-e2e-deployed: export RPC_URL?=wss://poc3-rpc.polkadot.io
test-e2e-deployed: 		## runs only end-to-end (e2e) tests against a deployed testnet (defaults to Alexander (wss://poc3-rpc.polkadot.io) if RPC_URL is not set)
	@docker build . -t gsrpc-test
	@docker run --rm -e RPC_URL gsrpc-test go test -v github.com/centrifuge/go-substrate-rpc-client/teste2e

run-substrate-docker: 		## runs the Substrate Default Docker image, this can be used to run the tests
	docker run -p 9933:9933 -p 9944:9944 -p 30333:30333 parity/substrate:latest-v1.0 --dev --rpc-external --ws-external

run-substrate-docker-v2: 	## runs the Substrate Default Docker image, this can be used to run the tests
	docker run -p 9933:9933 -p 9944:9944 -p 30333:30333 parity/substrate:latest --dev --rpc-external --ws-external

help: 				## shows this help
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST)

.PHONY: clean install lint lint-fix test test-dockerized run-substrate-docker run-substrate-docker-v2
