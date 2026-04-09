package driver

import (
	"testing"
)

// --- Test for s3_driver.go. ---.
func TestNewS3Driver(t *testing.T) {
	bucket := "test-bucket"
	mockClient := &mockS3Client{}             // Create mock client
	driver := NewS3Driver(mockClient, bucket) // Pass mock client

	if driver == nil {
		t.Fatal("NewS3Driver() returned nil")
	}
	if driver.bucket != bucket {
		t.Errorf("Expected bucket name %q, got %q", bucket, driver.bucket)
	}
	// Check if the client is the injected mock
	if driver.client == nil {
		t.Error("NewS3Driver client is nil")
	}
}
