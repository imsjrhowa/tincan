package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"tincan/pkg/s3client"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List files in S3 bucket",
	RunE:  runList,
}

func runList(cmd *cobra.Command, args []string) error {
	client, err := s3client.New()
	if err != nil {
		return fmt.Errorf("failed to create S3 client: %w", err)
	}

	files, err := client.List()
	if err != nil {
		return fmt.Errorf("failed to list files: %w", err)
	}

	if len(files) == 0 {
		fmt.Println("No files found in bucket")
		return nil
	}

	fmt.Println("Files in bucket:")
	for _, file := range files {
		fmt.Printf("  %s\n", file)
	}

	return nil
}