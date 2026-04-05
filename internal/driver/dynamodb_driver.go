package driver

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DynamoDBDriver implements storage operations using AWS DynamoDB.
type DynamoDBDriver struct {
	client    *dynamodb.Client
	tableName string
}

// NewDynamoDBDriver creates a new instance of DynamoDBDriver with an injected client.
func NewDynamoDBDriver(cfg aws.Config, tableName string) *DynamoDBDriver {
	return &DynamoDBDriver{
		client:    dynamodb.NewFromConfig(cfg),
		tableName: tableName,
	}
}

// StoreResult stores a result object in DynamoDB.
func (d *DynamoDBDriver) StoreResult(ctx context.Context, id string, version int, result interface{}) error {
	if d.client == nil {
		return errors.New("dynamodb client not initialized")
	}
	if id == "" {
		return errors.New("cannot store result with empty ID")
	}

	_, err := d.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(d.tableName),
		Item: map[string]types.AttributeValue{
			"ID":      &types.AttributeValueMemberS{Value: id},
			"Version": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", version)},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to store result in DynamoDB: %w", err)
	}

	fmt.Printf("Stored result for ID %s in DynamoDB table %s\n", id, d.tableName)
	return nil
}

// GetAllRecordsPaginated fetches all records (or latest versions) from DynamoDB.
func (d *DynamoDBDriver) GetAllRecordsPaginated(ctx context.Context, lastKey map[string]types.AttributeValue) ([]map[string]types.AttributeValue, map[string]types.AttributeValue, error) {
	input := &dynamodb.ScanInput{
		TableName:         aws.String(d.tableName),
		ExclusiveStartKey: lastKey,
	}

	result, err := d.client.Scan(ctx, input)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to scan DynamoDB: %w", err)
	}

	return result.Items, result.LastEvaluatedKey, nil
}

// GetLatestRecord fetches the latest version for a given ID using Query.
func (d *DynamoDBDriver) GetLatestRecord(ctx context.Context, id string) (map[string]types.AttributeValue, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(d.tableName),
		KeyConditionExpression: aws.String("ID = :id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":id": &types.AttributeValueMemberS{Value: id},
		},
		ScanIndexForward: aws.Bool(false), // Sort by Version descending
		Limit:            aws.Int32(1),    // Get only the latest
	}

	result, err := d.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query latest record: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, nil
	}

	return result.Items[0], nil
}
