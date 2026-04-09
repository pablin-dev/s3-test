package driver

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// --- Tests for get_all_records_paginated.go. ---.
func TestGetAllRecordsPaginated(t *testing.T) {
	ctx := context.Background() // Declare ctx
	tableName := "test-table"
	mockLastKey := map[string]types.AttributeValue{"ID": &types.AttributeValueMemberS{Value: "last-item"}}

	// Test case: successful scan
	expectedItems := []map[string]types.AttributeValue{
		{"ID": &types.AttributeValueMemberS{Value: "item1"}},
		{"ID": &types.AttributeValueMemberS{Value: "item2"}},
	}
	expectedLastEvaluatedKey := map[string]types.AttributeValue{"ID": &types.AttributeValueMemberS{Value: "next-key"}}

	mockClientSuccess := &mockDynamoDBClient{ // Use mock from mocks_test.go
		scanFunc: func(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
			if *params.TableName != tableName || !reflect.DeepEqual(params.ExclusiveStartKey, mockLastKey) { // reflect.DeepEqual requires exact type match, which we are using
				t.Errorf("Scan called with incorrect parameters")
			}
			return &dynamodb.ScanOutput{
				Items:            expectedItems,
				LastEvaluatedKey: expectedLastEvaluatedKey,
			}, nil
		},
	}
	driverSuccess := NewDynamoDBDriver(mockClientSuccess, tableName) // Pass mock client
	items, lastKey, err := driverSuccess.GetAllRecordsPaginated(ctx, mockLastKey)
	if err != nil {
		t.Errorf("GetAllRecordsPaginated failed: %v", err)
	}
	if !reflect.DeepEqual(items, expectedItems) {
		t.Errorf("Expected items %v, got %v", expectedItems, items)
	}
	if !reflect.DeepEqual(lastKey, expectedLastEvaluatedKey) {
		t.Errorf("Expected lastKey %v, got %v", expectedLastEvaluatedKey, lastKey)
	}

	// Test case: Scan returns an error
	mockClientError := &mockDynamoDBClient{ // Use mock from mocks_test.go
		scanFunc: func(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
			return nil, errors.New("simulated Scan error")
		},
	}
	driverError := NewDynamoDBDriver(mockClientError, tableName)
	_, _, err = driverError.GetAllRecordsPaginated(ctx, mockLastKey)
	if err == nil || err.Error() != "failed to scan DynamoDB: simulated Scan error" {
		t.Errorf("Expected DynamoDB Scan error, got: %v", err)
	}

	// Test case: nil client (should be handled by NewDynamoDBDriver)
	driverNilClient := NewDynamoDBDriver(nil, tableName) // Passing nil client
	_, _, err = driverNilClient.GetAllRecordsPaginated(ctx, mockLastKey)
	if err == nil || err.Error() != "dynamodb client not initialized" {
		t.Errorf("Expected error for nil client, got: %v", err)
	}
}
