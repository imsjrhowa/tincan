package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"tincan/pkg/s3client"
)

var uploadCmd = &cobra.Command{
	Use:   "upload [file]",
	Short: "Upload a file to S3",
	Args:  cobra.ExactArgs(1),
	RunE:  runUpload,
}

func runUpload(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filePath)
	}

	client, err := s3client.New()
	if err != nil {
		return fmt.Errorf("failed to create S3 client: %w", err)
	}

	fileName := filepath.Base(filePath)
	fmt.Printf("Uploading %s...\n", fileName)

	if err := client.Upload(filePath, fileName); err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	fmt.Printf("Successfully uploaded %s\n", fileName)
	return nil
}