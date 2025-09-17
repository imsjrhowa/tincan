# TinCan

A simple CLI tool for transferring files between machines using Amazon S3 as intermediary storage.

## Features

- Upload files to S3 from any machine
- Download files from S3 to any machine
- List files in your bucket
- Clean up (delete all files)
- **Web Interface** - GUI for easy file management
- **Embedded Credentials** - Build portable executables with credentials baked in
- Simple configuration

## Setup

1. **Install Go** (if not already installed)

2. **Build the binary:**
   ```bash
   # Using Makefile (recommended)
   make build

   # Or manually
   go build -o tincan ./cmd/tincan
   ```

3. **Configure AWS credentials** (one of these methods):
   - AWS credentials file (`~/.aws/credentials`)
   - Environment variables (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`)
   - IAM roles (if running on EC2)

4. **Set your S3 bucket:**
   ```bash
   export TINCAN_BUCKET=your-bucket-name
   ```

   Or create a config file at `~/.config/tincan.yaml` or `./tincan.yaml`:
   ```yaml
   bucket_name: your-bucket-name
   aws_region: us-east-1
   ```

## Usage

### Command Line Interface

```bash
# Upload a file
tincan upload myfile.txt

# Download a file
tincan download myfile.txt

# List all files
tincan list

# Clean up all files
tincan clean
```

### Web Interface

Start the web server for a GUI experience:

```bash
# Start web interface
tincan web
```

Then open your browser to `http://localhost:8080` for:
- Drag & drop file uploads
- Browse and download files
- Delete operations with confirmation
- Real-time file listing

## AWS Permissions

Your IAM user needs these S3 permissions:
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "s3:GetObject",
                "s3:PutObject",
                "s3:DeleteObject",
                "s3:ListBucket"
            ],
            "Resource": [
                "arn:aws:s3:::your-bucket-name",
                "arn:aws:s3:::your-bucket-name/*"
            ]
        }
    ]
}
```

## Building Options

### Using Makefile (Recommended)

```bash
# Build for current platform
make build

# Build for all platforms (Windows, macOS, Linux)
make build-all

# Run tests
make test

# Clean build artifacts
make clean

# Install to GOPATH/bin
make install

# Run directly without building
make run
```

### Manual Build (requires environment variables/config)

```bash
# Current platform
go build -o tincan ./cmd/tincan

# All platforms manually
GOOS=windows GOARCH=amd64 go build -o tincan.exe ./cmd/tincan
GOOS=darwin GOARCH=amd64 go build -o tincan-mac ./cmd/tincan
GOOS=linux GOARCH=amd64 go build -o tincan-linux ./cmd/tincan

# Run tests
go test -v ./...
```

### Embedded Credentials Build (portable, no setup required)

For maximum portability, you can embed AWS credentials directly into the executable:

1. **Copy the template:**
   ```bash
   cp build-with-creds.bat.template build-with-creds.bat
   ```

2. **Edit `build-with-creds.bat` with your credentials:**
   ```bat
   set ACCESS_KEY=your-aws-access-key
   set SECRET_KEY=your-aws-secret-key
   set REGION=your-aws-region
   set BUCKET=your-bucket-name
   ```

3. **Build:**
   ```bash
   ./build-with-creds.bat
   ```

This creates `tincan-embedded.exe` - a completely portable executable that:
- Contains all AWS credentials
- Requires no environment variables
- Needs no configuration files
- Can be copied to any machine and run immediately

**Security Note:** The `build-with-creds.bat` file (with real credentials) is automatically ignored by Git.

## Architecture

TinCan is built using Go with the following structure:

- **cmd/tincan/**: Main CLI application entry point and command definitions
  - `main.go`: Root command setup and initialization
  - Individual command implementations (`upload.go`, `download.go`, `list.go`, `clean.go`, etc.)
- **pkg/s3client/**: S3 operations abstraction layer - handles all AWS S3 interactions
- **internal/config/**: Configuration management using Viper - supports YAML files and environment variables

## Dependencies

- `github.com/spf13/cobra`: CLI framework
- `github.com/spf13/viper`: Configuration management
- `github.com/aws/aws-sdk-go-v2`: AWS SDK for S3 operations