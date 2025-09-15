BINARY_NAME=tincan

.PHONY: build clean test install

build:
	go build -o $(BINARY_NAME) ./cmd/tincan

build-all:
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME).exe ./cmd/tincan
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-mac ./cmd/tincan
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux ./cmd/tincan

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