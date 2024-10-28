# Directory containing the Makefile.
PROJECT_ROOT = $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

MOCKGEN = bin/mockgen

.PHONY: all
all: build lint test

.PHONY: build
build:
	go build ./...

.PHONY: test
test:
	go test ./...

.PHONY: cover
cover:
	#	ignore mock.go files
	go test -race -coverprofile=cover.out.tmp -coverpkg=./... ./... \
	&& cat cover.out.tmp | grep -v "mock.go" > cover.out \
	&& go tool cover -html=cover.out -o cover.html

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: lint
lint:
	golangci-lint run

.PHONY: generate
generate: $(TOOLS)
	go generate -x ./...
	make -C doc generate

$(MOCKGEN): go.mod
	go install go.uber.org/mock/mockgen