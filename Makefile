# Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
# Copyright (C) 2019  Philip Stanislaus, Philip Stehlik, Vimukthi Wickramasinghe
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

clean: ##clean vendor's folder. Should be run before a make install
	@echo 'cleaning previous /vendor folder'
	@rm -rf vendor/
	@echo 'done cleaning'

install: ## Install Dependencies
	@command -v dep >/dev/null 2>&1 || go get -u github.com/golang/dep/...
	@command -v golangci-lint >/dev/null 2>&1 || go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	@dep ensure

lint: ## runs linters on go code
	@golangci-lint run

lint-fix: ## runs linters on go code and automatically fixes issues
	@golangci-lint run --fix

test: ## runs all tests in project
	@go test -v ./...

test-cover: ## runs all tests in project and report coverage
	@go test -cover -count=1 ./...

test-dockerized: ## runs the tests in a docker container against the Substrate Default Docker image
	@docker-compose build
	@docker-compose up --abort-on-container-exit

run-substrate-docker: ## runs the Substrate Default Docker image
	docker run -p 9933:9933 -p 9944:9944 -p 30333:30333 parity/substrate:latest-v1.0 --dev --rpc-external --ws-external

.PHONY: clean install lint lint-fix test test-dockerized run-substrate-docker
