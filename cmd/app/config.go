package main

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gopkg.in/yaml.v3"
)

type Config struct {
	AWS struct {
		Region        string `yaml:"region"`
		S3Bucket      string `yaml:"s3_bucket"`
		DynamoDBTable string `yaml:"dynamodb_table"`
	} `yaml:"aws"`
}

type S3Bucket string
type DynamoDBTable string

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	return &cfg, err
}

func ProvideAWSConfig(ctx context.Context, cfg *Config) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx, config.WithRegion(cfg.AWS.Region))
}

func ProvideS3Bucket(cfg *Config) S3Bucket { return S3Bucket(cfg.AWS.S3Bucket) }
func ProvideDynamoDBTable(cfg *Config) DynamoDBTable { return DynamoDBTable(cfg.AWS.DynamoDBTable) }

func ProvideS3Client(cfg aws.Config) *s3.Client {
	return s3.NewFromConfig(cfg)
}

func ProvideDynamoDBClient(cfg aws.Config) *dynamodb.Client {
	return dynamodb.NewFromConfig(cfg)
}
