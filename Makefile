#@IgnoreInspection BashAddShebang
export ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))
export CGO_ENABLED=0
export GO111MODULE=on
export GOPROXY=https://proxy.golang.org,direct

export ENV=development

.PHONY: all format lint build test

all: build lint test

.which-go:
	@which go > /dev/null || (echo "install go from https://golang.org/dl/" & exit 1)

format: .which-go
	gofmt -s -w $(ROOT)

.which-lint:
	@which golangci-lint > /dev/null || (echo "install golangci-lint from https://github.com/golangci/golangci-lint" & exit 1)

lint: .which-lint
	golangci-lint run

build: .which-go
	go build -v -o $(ROOT)/bin/api -ldflags="-s -w" $(ROOT)/cmd/api/*.go

test: .which-go
	CGO_ENABLED=1 go test -race -coverprofile=coverage.txt -covermode=atomic $(ROOT)/...
