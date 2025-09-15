package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"tincan/pkg/s3client"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove all files from S3 bucket",
	RunE:  runClean,
}

func runClean(cmd *cobra.Command, args []string) error {
	client, err := s3client.New()
	if err != nil {
		return fmt.Errorf("failed to create S3 client: %w", err)
	}

	files, err := client.List()
	if err != nil {
		return fmt.Errorf("failed to list files: %w", err)
	}

	if len(files) == 0 {
		fmt.Println("No files to clean")
		return nil
	}

	fmt.Printf("The following %d files will be deleted:\n", len(files))
	for _, file := range files {
		fmt.Printf("  - %s\n", file)
	}
	fmt.Print("\nAre you sure you want to delete these files? (y/N): ")

	var response string
	fmt.Scanln(&response)
	if response != "y" && response != "Y" {
		fmt.Println("Clean cancelled")
		return nil
	}

	for _, file := range files {
		if err := client.Delete(file); err != nil {
			fmt.Printf("Failed to delete %s: %v\n", file, err)
		} else {
			fmt.Printf("Deleted %s\n", file)
		}
	}

	fmt.Println("Clean completed")
	return nil
}