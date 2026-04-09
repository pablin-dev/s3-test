package driver

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodb_types "github.com/aws/aws-sdk-go-v2/service/dynamodb/types" // Aliased
)

// --- Tests for store_result.go ---
func TestStoreResult(t *testing.T) {
	ctx := context.Background()
	tableName := "test-table"
	testID := "test-123"
	testVersion := 1
	testResult := map[string]string{"data": "value"}

	// Test case: successful storage
	mockClientSuccess := &mockDynamoDBClient{ // Use mock from mocks_test.go
		putItemFunc: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			if *params.TableName != tableName || params.Item["ID"].(*dynamodb_types.AttributeValueMemberS).Value != testID || params.Item["Version"].(*dynamodb_types.AttributeValueMemberN).Value != "1" {
				t.Errorf("PutItem called with incorrect parameters")
			}
			return &dynamodb.PutItemOutput{}, nil
		},
	}
	driverSuccess := NewDynamoDBDriver(mockClientSuccess, tableName) // Pass mock client
	err := driverSuccess.StoreResult(ctx, testID, testVersion, testResult)
	if err != nil {
		t.Errorf("StoreResult failed: %v", err)
	}

	// Test case: empty ID
	driverEmptyID := NewDynamoDBDriver(mockClientSuccess, tableName)
	err = driverEmptyID.StoreResult(ctx, "", testVersion, testResult)
	if err == nil || err.Error() != "cannot store result with empty ID" {
		t.Errorf("Expected error for empty ID, got: %v", err)
	}

	// Test case: nil client passed to driver constructor
	driverNilClient := NewDynamoDBDriver(nil, tableName) // Passing nil client
	err = driverNilClient.StoreResult(ctx, testID, testVersion, testResult)
	if err == nil || err.Error() != "dynamodb client not initialized" {
		t.Errorf("Expected error for nil client, got: %v", err)
	}

	// Test case: PutItem returns an error
	mockClientError := &mockDynamoDBClient{ // Use mock from mocks_test.go
		putItemFunc: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			return nil, errors.New("simulated PutItem error")
		},
	}
	driverError := NewDynamoDBDriver(mockClientError, tableName)
	err = driverError.StoreResult(ctx, testID, testVersion, testResult)
	if err == nil || err.Error() != "failed to store result in DynamoDB: simulated PutItem error" {
		t.Errorf("Expected DynamoDB PutItem error, got: %v", err)
	}
}
