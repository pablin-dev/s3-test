package driver

import (
	"context"
	"errors"
	"testing"

	// Import aws for aws.Equal
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pablodev/s3-test/internal/entity"
)

// --- Tests for upload_yaml.go ---
func TestUploadYaml(t *testing.T) {
	ctx := context.Background() // Declare ctx
	bucket := "test-bucket"
	testYamlFile := entity.YamlFile{
		ID:         "my-yaml-file",
		Expression: "key: value", // Removed newline
	}

	// Test case: successful upload
	mockClientSuccess := &mockS3Client{ // Use mock from mocks_test.go
		putObjectFunc: func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
			if *params.Bucket != bucket || *params.Key != "my-yaml-file.yaml" {
				t.Errorf("PutObject called with incorrect parameters")
			}
			// Read body to verify content
			bodyBytes := make([]byte, 1024)
			n, err := params.Body.Read(bodyBytes)
			if err != nil && err.Error() != "EOF" { // Ignore EOF on empty read
				t.Errorf("Error reading body: %v", err)
			}
			body := string(bodyBytes[:n])
			if body != testYamlFile.Expression {
				// Using %s to avoid potential issues with %q and newlines/special chars.
				t.Errorf("Expected body %s, got %s", testYamlFile.Expression, body)
			}
			return &s3.PutObjectOutput{}, nil
		},
	}
	driverSuccess := NewS3Driver(mockClientSuccess, bucket) // Pass mock client
	err := driverSuccess.UploadYaml(ctx, testYamlFile)
	if err != nil {
		t.Errorf("UploadYaml failed: %v", err)
	}

	// Test case: empty ID
	driverEmptyID := NewS3Driver(mockClientSuccess, bucket)
	err = driverEmptyID.UploadYaml(ctx, entity.YamlFile{ID: "", Expression: "content"})
	if err == nil || err.Error() != "cannot upload file with empty ID" {
		t.Errorf("Expected error for empty ID, got: %v", err)
	}

	// Test case: nil client passed to driver constructor
	driverNilClient := NewS3Driver(nil, bucket) // Passing nil client
	err = driverNilClient.UploadYaml(ctx, testYamlFile)
	if err == nil || err.Error() != "s3 client not initialized" {
		t.Errorf("Expected error for nil client, got: %v", err)
	}

	// Test case: PutObject returns an error
	mockClientError := &mockS3Client{ // Use mock from mocks_test.go
		putObjectFunc: func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
			return nil, errors.New("simulated PutObject error")
		},
	}
	driverError := NewS3Driver(mockClientError, bucket) // Pass mock client
	err = driverError.UploadYaml(ctx, testYamlFile)
	if err == nil || err.Error() != "failed to upload to S3: simulated PutObject error" {
		t.Errorf("Expected S3 PutObject error, got: %v", err)
	}
}
