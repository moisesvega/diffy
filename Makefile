# Directory containing the Makefile.
PROJECT_ROOT = $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

# Directories containing independent Go modules.
MODULE_DIRS = .

.PHONY: all
all: lint build test

.PHONY: build
build:
	go build ./...

.PHONY: test
test:
	@$(foreach dir,$(MODULE_DIRS),(cd $(dir) && go test -race ./...) &&) true

.PHONY: cover
cover:
	@$(foreach dir,$(MODULE_DIRS), ( \
		cd $(dir) && \
		go test -race -coverprofile=cover.out.tmp -coverpkg=./... ./... \
		&& cat cover.out.tmp | grep -v "mock.go" > cover.out \
		&& go tool cover -html=cover.out -o cover.html) &&) true

.PHONY: tidy
tidy:
	@$(foreach dir,$(MODULES),(cd $(dir) && go mod tidy) &&) true

.PHONY: lint
lint: golangci-lint

.PHONY: golangci-lint
golangci-lint:
	@$(foreach mod,$(MODULE_DIRS), \
		(cd $(mod) && \
		echo "[lint] golangci-lint: $(mod)" && \
		golangci-lint run --path-prefix $(mod)) &&) true