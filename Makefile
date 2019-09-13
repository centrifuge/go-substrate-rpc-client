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
	@go test ./...

.PHONY: clean install lint lint-fix test
