package driver

import (
	"testing"
	// Import aws.
)

// --- Tests for dynamodb_driver.go. ---.
func TestNewDynamoDBDriver(t *testing.T) {
	tableName := "test-table"
	mockClient := &mockDynamoDBClient{}                // Create mock client
	driver := NewDynamoDBDriver(mockClient, tableName) // Pass mock client

	if driver == nil {
		t.Fatal("NewDynamoDBDriver() returned nil")
	}
	if driver.tableName != tableName {
		t.Errorf("Expected table name %q, got %q", tableName, driver.tableName)
	}
	// Check if the client is the injected mock
	if driver.client == nil {
		t.Error("NewDynamoDBDriver client is nil")
	}
}
