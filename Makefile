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
