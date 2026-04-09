package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pablodev/s3-test/internal/driver"
)

// Wrapper for S3 client.
type S3ClientWrapper struct {
	*s3.Client
}

func (w *S3ClientWrapper) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	return w.Client.PutObject(ctx, params, optFns...)
}
func (w *S3ClientWrapper) ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	return w.Client.ListObjectsV2(ctx, params, optFns...)
}

func ProvideS3ClientAPI(client *s3.Client) driver.S3ClientAPI {
	return &S3ClientWrapper{Client: client}
}

// Wrapper for DynamoDB client.
type DynamoDBClientWrapper struct {
	*dynamodb.Client
}

func (w *DynamoDBClientWrapper) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return w.Client.PutItem(ctx, params, optFns...)
}
func (w *DynamoDBClientWrapper) Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	return w.Client.Scan(ctx, params, optFns...)
}
func (w *DynamoDBClientWrapper) Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	return w.Client.Query(ctx, params, optFns...)
}

func ProvideDynamoDBClientAPI(client *dynamodb.Client) driver.DynamoDBClientAPI {
	return &DynamoDBClientWrapper{Client: client}
}
