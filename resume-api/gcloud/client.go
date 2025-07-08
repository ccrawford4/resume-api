package gcloud

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"resume-api/parser"

	"cloud.google.com/go/storage"
)

// ListBucketFiles fetches all file names in the specified GCS bucket
func ListBucketFiles(bucketName string) ([]string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage client: %w", err)
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)
	it := bucket.Objects(ctx, nil)

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
func DownloadFile(bucketName, objectName, destDir string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create storage client: %w", err)
	}
	defer client.Close()

	rc, err := client.Bucket(bucketName).Object(objectName).NewReader(ctx)
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

func DownloadAndParseAllPDFs(bucketName string) (map[string]string, error) {
	files, err := ListBucketFiles(bucketName)
	if err != nil {
		return nil, err
	}
	destDir := "../resumes"
	os.MkdirAll(destDir, 0755)
	for _, file := range files {
		if filepath.Ext(file) == ".pdf" {
			if err := DownloadFile(bucketName, file, destDir); err != nil {
				fmt.Printf("Failed to download %s: %v\n", file, err)
			}
		}
	}
	return parser.ParseResume(), nil
}
