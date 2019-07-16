clean: ##clean vendor's folder. Should be run before a make install
	@echo 'cleaning previous /vendor folder'
	@rm -rf vendor/
	@echo 'done cleaning'

install-deps: ## Install Dependencies
	@command -v dep >/dev/null 2>&1 || go get -u github.com/golang/dep/...
	@dep ensure
	@curl -L https://git.io/vp6lP | sh -s ${GOMETALINTER_VERSION}
	@mv ./bin/* $(GOPATH)/bin/; rm -rf ./bin

lint-check: ## runs linters on go code
	@gometalinter --exclude=anchors/service.go  --disable-all --enable=golint --enable=goimports --enable=vet --enable=nakedret \
	--enable=staticcheck --vendor --skip=resources --skip=testingutils --skip=protobufs  --deadline=1m ./...;

format-go: ## formats go code
	@goimports -w .

build: install-deps