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
		size := formatBytes(file.Size)
		date := file.LastModified.Format("2006-01-02 15:04:05")
		fmt.Printf("  %-40s %10s  %s\n", file.Name, size, date)
	}

	return nil
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}