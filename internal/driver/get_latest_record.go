package driver

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// GetLatestRecord fetches the latest version for a given ID using Query.
func (d *DynamoDBDriver) GetLatestRecord(ctx context.Context, id string) (map[string]types.AttributeValue, error) {
	if d.client == nil {
		return nil, fmt.Errorf("dynamodb client not initialized")
	}

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
