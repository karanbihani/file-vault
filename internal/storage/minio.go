package storage

import (
	"context"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Config holds the configuration for the MinIO client, read from environment variables.
type Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	UseSSL          bool
}

// Client is a wrapper around the MinIO client that provides our application's storage methods.
type Client struct {
	minioClient *minio.Client
	bucketName  string
}

// NewClient creates and initializes a new MinIO client.
// It also checks if the required bucket exists and creates it if it doesn't.
func NewClient(ctx context.Context, config Config) *Client {
	minioClient, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		log.Fatalf("Failed to initialize MinIO client: %v", err)
	}

	// Check if the bucket already exists.
	exists, err := minioClient.BucketExists(ctx, config.BucketName)
	if err != nil {
		log.Fatalf("Failed to check if MinIO bucket exists: %v", err)
	}

	// If the bucket doesn't exist, create it. This makes the setup self-healing.
	if !exists {
		err = minioClient.MakeBucket(ctx, config.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalf("Failed to create MinIO bucket: %v", err)
		}
		log.Printf("Successfully created bucket: %s\n", config.BucketName)
	}

	return &Client{
		minioClient: minioClient,
		bucketName:  config.BucketName,
	}
}

// Save uploads a file to the MinIO bucket.
func (c *Client) Save(ctx context.Context, objectName string, data io.Reader, size int64, contentType string) error {
	_, err := c.minioClient.PutObject(ctx, c.bucketName, objectName, data, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

// Get retrieves a file object from the MinIO bucket.
// CORRECTED: Return type is now *minio.Object to give access to the Stat() method.
func (c *Client) Get(ctx context.Context, objectName string) (*minio.Object, error) {
	return c.minioClient.GetObject(ctx, c.bucketName, objectName, minio.GetObjectOptions{})
}

// Delete removes a file object from the MinIO bucket.
func (c *Client) Delete(ctx context.Context, objectName string) error {
	return c.minioClient.RemoveObject(ctx, c.bucketName, objectName, minio.RemoveObjectOptions{})
}
