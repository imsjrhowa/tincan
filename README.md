# TinCan

A simple CLI tool for transferring files between machines using Amazon S3 as intermediary storage.

## Features

- Upload files to S3 from any machine
- Download files from S3 to any machine
- List files in your bucket
- Clean up (delete all files)
- Simple configuration

## Setup

1. **Install Go** (if not already installed)
2. **Build the binary:**
   ```bash
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

   Or create a config file at `~/.config/tincan.yaml`:
   ```yaml
   bucket_name: your-bucket-name
   aws_region: us-east-1
   ```

## Usage

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

## Building for Multiple Platforms

```bash
# Windows
GOOS=windows GOARCH=amd64 go build -o tincan.exe ./cmd/tincan

# macOS
GOOS=darwin GOARCH=amd64 go build -o tincan-mac ./cmd/tincan

# Linux
GOOS=linux GOARCH=amd64 go build -o tincan-linux ./cmd/tincan
```