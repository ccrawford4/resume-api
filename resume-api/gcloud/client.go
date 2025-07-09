package gcloud

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"resume-api/parser"

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
func (gc *GoogleCloudClient) DownloadFile(bucketName, objectName, destDir string) error {
	rc, err := gc.storageClient.Bucket(bucketName).Object(objectName).NewReader(gc.ctx)
	if err != nil {
		return fmt.Errorf("failed to create object reader: %w", err)
	}
	defer rc.Close()

	localPath := filepath.Join(destDir, objectName)
	outFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, rc)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}
	return nil
}

func (gc *GoogleCloudClient) DownloadAndParseAllPDFs(bucketName string) (map[string]string, error) {
	files, err := gc.ListBucketFiles(bucketName)
	if err != nil {
		return nil, err
	}
	destDir := "../resumes"
	os.MkdirAll(destDir, 0755)
	for _, file := range files {
		if filepath.Ext(file) == ".pdf" {
			if err := gc.DownloadFile(bucketName, file, destDir); err != nil {
				fmt.Printf("Failed to download %s: %v\n", file, err)
			}
		}
	}
	return parser.ParseResume(), nil
}
