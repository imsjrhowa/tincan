package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"tincan/pkg/s3client"
)

var downloadCmd = &cobra.Command{
	Use:   "download [filename]",
	Short: "Download a file from S3",
	Args:  cobra.ExactArgs(1),
	RunE:  runDownload,
}

func runDownload(cmd *cobra.Command, args []string) error {
	fileName := args[0]

	client, err := s3client.New()
	if err != nil {
		return fmt.Errorf("failed to create S3 client: %w", err)
	}

	// Check if file already exists locally
	if _, err := os.Stat(fileName); err == nil {
		fmt.Printf("File %s already exists. Overwrite? (y/N): ", fileName)
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Download cancelled")
			return nil
		}
	}

	fmt.Printf("Downloading %s...\n", fileName)

	if err := client.Download(fileName, fileName); err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}

	fmt.Printf("Successfully downloaded %s\n", fileName)
	return nil
}