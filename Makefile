BINARY_NAME=tincan
VERSION?=1.0.0
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT) -X main.BuildDate=$(BUILD_DATE)"

.PHONY: build clean test install

build:
	go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/tincan

build-all:
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME).exe ./cmd/tincan
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-mac ./cmd/tincan
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-linux ./cmd/tincan

clean:
	go clean
	rm -f $(BINARY_NAME) $(BINARY_NAME).exe $(BINARY_NAME)-mac $(BINARY_NAME)-linux

test:
	go test -v ./...

install:
	go install ./cmd/tincan

run:
	go run ./cmd/tincan

.DEFAULT_GOAL := build