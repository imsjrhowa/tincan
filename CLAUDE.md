# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

TinCan is a CLI tool for transferring files between machines using Amazon S3 as intermediary storage. The application is built in Go using the Cobra CLI framework and AWS SDK v2.

## Architecture

- **cmd/tincan/**: Main CLI application entry point and command definitions
  - `main.go`: Root command setup and initialization
  - `upload.go`, `download.go`, `list.go`, `clean.go`: Individual command implementations
- **pkg/s3client/**: S3 operations abstraction layer - handles all AWS S3 interactions
- **internal/config/**: Configuration management using Viper - supports YAML files and environment variables

## Build and Development Commands

```bash
# Build for current platform
make build
# or
go build -o tincan ./cmd/tincan

# Build for all platforms (Windows, macOS, Linux)
make build-all

# Run tests
make test
# or
go test -v ./...

# Clean build artifacts
make clean

# Install to GOPATH/bin
make install

# Run directly without building
make run
# or
go run ./cmd/tincan
```

## Configuration

The application uses two configuration methods:
1. Environment variables with `TINCAN_` prefix (e.g., `TINCAN_BUCKET_NAME`)
2. YAML config file at `~/.config/tincan.yaml` or `./tincan.yaml`

Required: S3 bucket name via `TINCAN_BUCKET` env var or `bucket_name` in config file.

## Key Dependencies

- `github.com/spf13/cobra`: CLI framework
- `github.com/spf13/viper`: Configuration management
- `github.com/aws/aws-sdk-go-v2`: AWS SDK for S3 operations

## AWS Setup Required

The application requires AWS credentials and appropriate S3 permissions for GetObject, PutObject, DeleteObject, and ListBucket operations.