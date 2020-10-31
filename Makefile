.PHONY: default
default: all ;

format:
	go fmt .
	goimports -w .

lint:
	golangci-lint run
	go vet ./...

test: lint
	go test -v ./...

test-nolint:
	go test -v ./...

all: format lint test
