package driver

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

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

	fmt.Printf("Stored result for ID %s in DynamoDB table %s", id, d.tableName)
	return nil
}
