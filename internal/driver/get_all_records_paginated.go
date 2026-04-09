package driver

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// GetAllRecordsPaginated fetches all records (or latest versions) from DynamoDB.
func (d *DynamoDBDriver) GetAllRecordsPaginated(ctx context.Context, lastKey map[string]types.AttributeValue) ([]map[string]types.AttributeValue, map[string]types.AttributeValue, error) {
	if d.client == nil {
		return nil, nil, fmt.Errorf("dynamodb client not initialized")
	}

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
