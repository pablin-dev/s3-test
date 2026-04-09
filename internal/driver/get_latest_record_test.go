package driver

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// --- Tests for get_latest_record.go. ---.
func TestGetLatestRecord(t *testing.T) {
	ctx := context.Background() // Declare ctx
	tableName := "test-table"
	testID := "some-id"

	// Test case: successful query for a record
	expectedItem := map[string]types.AttributeValue{"ID": &types.AttributeValueMemberS{Value: testID}, "Version": &types.AttributeValueMemberN{Value: "5"}}
	mockClientSuccess := &mockDynamoDBClient{ // Use mock from mocks_test.go
		queryFunc: func(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			if *params.TableName != tableName || params.ExpressionAttributeValues[":id"].(*types.AttributeValueMemberS).Value != testID {
				t.Errorf("Query called with incorrect parameters")
			}
			return &dynamodb.QueryOutput{
				Items: []map[string]types.AttributeValue{expectedItem},
			}, nil
		},
	}
	driverSuccess := NewDynamoDBDriver(mockClientSuccess, tableName) // Pass mock client
	item, err := driverSuccess.GetLatestRecord(ctx, testID)
	if err != nil {
		t.Errorf("GetLatestRecord failed: %v", err)
	}
	if !reflect.DeepEqual(item, expectedItem) {
		t.Errorf("Expected item %v, got %v", expectedItem, item)
	}

	// Test case: no record found for ID
	mockClientNotFound := &mockDynamoDBClient{ // Use mock from mocks_test.go
		queryFunc: func(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			return &dynamodb.QueryOutput{Items: []map[string]types.AttributeValue{}}, nil
		},
	}
	driverNotFound := NewDynamoDBDriver(mockClientNotFound, tableName) // Pass mock client
	item, err = driverNotFound.GetLatestRecord(ctx, testID)
	if err != nil {
		t.Errorf("GetLatestRecord failed: %v", err)
	}
	if item != nil {
		t.Errorf("Expected nil item when not found, got %v", item)
	}

	// Test case: Query returns an error
	mockClientError := &mockDynamoDBClient{ // Use mock from mocks_test.go
		queryFunc: func(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			return nil, errors.New("simulated Query error")
		},
	}
	driverError := NewDynamoDBDriver(mockClientError, tableName) // Pass mock client
	_, err = driverError.GetLatestRecord(ctx, testID)
	if err == nil || err.Error() != "failed to query latest record: simulated Query error" {
		t.Errorf("Expected DynamoDB Query error, got: %v", err)
	}

	// Test case: nil client (should be handled by NewDynamoDBDriver)
	driverNilClient := NewDynamoDBDriver(nil, tableName) // Passing nil client
	_, err = driverNilClient.GetLatestRecord(ctx, testID)
	if err == nil || err.Error() != "dynamodb client not initialized" {
		t.Errorf("Expected error for nil client, got: %v", err)
	}
}
