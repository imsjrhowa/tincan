package s3client

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	embeddedcreds "tincan/internal/credentials"
)

type Client struct {
	s3Client   *s3.Client
	bucketName string
}

func New() (*Client, error) {
	var cfg aws.Config
	var err error
	var bucketName string

	// Check if credentials are embedded at build time
	if embeddedcreds.HasEmbeddedCredentials() {
		// Use embedded credentials
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(embeddedcreds.Region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				embeddedcreds.AccessKey,
				embeddedcreds.SecretKey,
				"",
			)),
		)
		bucketName = embeddedcreds.BucketName
	} else {
		// Fall back to environment variables or AWS credentials file
		cfg, err = config.LoadDefaultConfig(context.TODO())
		bucketName = os.Getenv("TINCAN_BUCKET")
		if bucketName == "" {
			return nil, fmt.Errorf("TINCAN_BUCKET environment variable is required when credentials are not embedded")
		}
	}

	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	return &Client{
		s3Client:   s3.NewFromConfig(cfg),
		bucketName: bucketName,
	}, nil
}

func (c *Client) Upload(filePath, key string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("unable to open file %q: %w", filePath, err)
	}
	defer file.Close()

	_, err = c.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("unable to upload %q to %q: %w", filePath, c.bucketName, err)
	}

	return nil
}

func (c *Client) Download(key, filePath string) error {
	result, err := c.s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("unable to download %q from %q: %w", key, c.bucketName, err)
	}
	defer result.Body.Close()

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("unable to create file %q: %w", filePath, err)
	}
	defer file.Close()

	_, err = io.Copy(file, result.Body)
	if err != nil {
		return fmt.Errorf("unable to write to file %q: %w", filePath, err)
	}

	return nil
}

func (c *Client) List() ([]string, error) {
	result, err := c.s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(c.bucketName),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to list objects in %q: %w", c.bucketName, err)
	}

	var files []string
	for _, obj := range result.Contents {
		files = append(files, *obj.Key)
	}

	return files, nil
}

func (c *Client) Delete(key string) error {
	_, err := c.s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("unable to delete %q from %q: %w", key, c.bucketName, err)
	}

	return nil
}