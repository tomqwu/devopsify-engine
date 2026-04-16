.PHONY: build test lint fmt coverage clean docker-build

BINARY_CLI=bin/dne
BINARY_SERVER=bin/dne-server
MODULE=github.com/deepnative/engine

build: fmt
	go build -o $(BINARY_CLI) ./cmd/dne
	go build -o $(BINARY_SERVER) ./cmd/dne-server

test:
	go test ./... -race -count=1

lint:
	golangci-lint run ./...

fmt:
	gofmt -w .
	goimports -w .

coverage:
	go test ./... -race -coverprofile=coverage.out -covermode=atomic
	go tool cover -html=coverage.out -o coverage.html

clean:
	rm -rf bin/ coverage.out coverage.html

docker-build:
	docker build -t deepnative/engine:latest .
