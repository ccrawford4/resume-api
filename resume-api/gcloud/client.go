package gcloud

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"resume-api/parser"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type GoogleCloudClient struct {
	storageClient *storage.Client
	ctx           context.Context
}

// createStorageClient creates a new storage client using service account credentials
func NewClient(ctx context.Context) (*GoogleCloudClient, error) {
	// Path to the service account key file
	keyFile := "../keys.json"

	// Check if the key file exists
	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("service account key file not found: %s", keyFile)
	}

	// Create client with service account credentials
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(keyFile))
	if err != nil {
		return nil, fmt.Errorf("failed to create storage client with service account: %w", err)
	}

	return &GoogleCloudClient{
		storageClient: client,
		ctx:           ctx,
	}, nil
}

// ListBucketFiles fetches all file names in the specified GCS bucket
func (gc *GoogleCloudClient) ListBucketFiles(bucketName string) ([]string, error) {
	bucket := gc.storageClient.Bucket(bucketName)
	it := bucket.Objects(gc.ctx, nil)

	var files []string
	for {
		attrs, err := it.Next()
		if err != nil {
			fmt.Printf("err: %s\n", err.Error())
			break
		}
		// return nil, fmt.Errorf("error iterating bucket objects: %w", err)
		files = append(files, attrs.Name)
	}
	return files, nil
}

// DownloadFile downloads a file from GCS to the local resumes directory
func (gc *GoogleCloudClient) DownloadFile(bucketName, objectName string) (*parser.File, error) {
	// Get the reader for the objeect
	rc, err := gc.storageClient.Bucket(bucketName).Object(objectName).NewReader(gc.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create object reader: %w", err)
	}
	defer rc.Close()

	// Get the presigned URL for the object
	url, err := gc.storageClient.Bucket(bucketName).SignedURL(objectName, &storage.SignedURLOptions{
		Method:  "GET",
		Expires: time.Now().Add(1 * time.Hour), // URL valid for 1 hour
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get signed URL: %w", err)
	}

	// Read the file content
	fileContent, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("failed to read object content: %w", err)
	}

	// Return the file object
	return &parser.File{
		Name:    objectName,
		Content: fileContent,
		URL:     url,
	}, nil
}

func (gc *GoogleCloudClient) DownloadAllPDFs(bucketName string) ([]*parser.File, error) {
	files, err := gc.ListBucketFiles(bucketName)
	if err != nil {
		return nil, err
	}

	var fileObjects []*parser.File
	for _, file := range files {
		if filepath.Ext(file) == ".pdf" {
			newFile, err := gc.DownloadFile(bucketName, file)
			if err != nil {
				return nil, fmt.Errorf("error downloading file %s: %w", file, err)
			}
			fileObjects = append(fileObjects, newFile)
		}
	}

	return fileObjects, nil
}
