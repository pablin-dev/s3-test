package driver

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws" // Import aws for aws.Equal
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/pablodev/s3-test/internal/entity"
)

// --- Tests for download_all_yaml.go. ---.
func TestDownloadAllYaml(t *testing.T) {
	ctx := context.Background() // Declare ctx
	bucket := "test-bucket"

	// Test case: successful download
	expectedObjects := []types.Object{
		{Key: aws.String("file1.yaml")},
		{Key: aws.String("file2.yaml")},
	}
	expectedFiles := []entity.YamlFile{
		{ID: "file1.yaml"},
		{ID: "file2.yaml"},
	}
	mockClientSuccess := &mockS3Client{ // Use mock from mocks_test.go
		listObjectsV2Func: func(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
			if *params.Bucket != bucket {
				t.Errorf("ListObjectsV2 called with incorrect parameters")
			}
			return &s3.ListObjectsV2Output{Contents: expectedObjects}, nil
		},
	}
	driverSuccess := NewS3Driver(mockClientSuccess, bucket) // Pass mock client
	files, err := driverSuccess.DownloadAllYaml(ctx)
	if err != nil {
		t.Errorf("DownloadAllYaml failed: %v", err)
	}
	if len(files) != len(expectedFiles) {
		t.Fatalf("Expected %d files, got %d", len(expectedFiles), len(files))
	}
	for i := range files {
		if files[i] != expectedFiles[i] {
			t.Errorf("Expected file %v, got %v", expectedFiles[i], files[i])
		}
	}

	// Test case: ListObjectsV2 returns an error
	mockClientError := &mockS3Client{ // Use mock from mocks_test.go
		listObjectsV2Func: func(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
			return nil, errors.New("simulated ListObjectsV2 error")
		},
	}
	driverError := NewS3Driver(mockClientError, bucket) // Pass mock client
	_, err = driverError.DownloadAllYaml(ctx)
	if err == nil || err.Error() != "failed to list S3 objects: simulated ListObjectsV2 error" {
		t.Errorf("Expected S3 ListObjectsV2 error, got: %v", err)
	}

	// Test case: nil client (should be handled by NewS3Driver)
	driverNilClient := NewS3Driver(nil, bucket) // Passing nil client
	_, err = driverNilClient.DownloadAllYaml(ctx)
	if err == nil || err.Error() != "s3 client not initialized" {
		t.Errorf("Expected error for nil client, got: %v", err)
	}
}
