package driver

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pablodev/s3-test/internal/entity"
)

// S3Driver implements storage operations using AWS S3.
type S3Driver struct {
	client *s3.Client
	bucket string
}

// NewS3Driver creates a new instance of S3Driver with an injected client.
func NewS3Driver(cfg aws.Config, bucket string) *S3Driver {
	return &S3Driver{
		client: s3.NewFromConfig(cfg),
		bucket: bucket,
	}
}

// UploadYaml uploads a YAML file content to S3.
func (d *S3Driver) UploadYaml(ctx context.Context, file entity.YamlFile) error {
	if d.client == nil {
		return errors.New("s3 client not initialized")
	}
	if file.ID == "" {
		return errors.New("cannot upload file with empty ID")
	}

	_, err := d.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(fmt.Sprintf("%s.yaml", file.ID)),
		Body:   strings.NewReader(file.Expression),
	})
	if err != nil {
		return fmt.Errorf("failed to upload to S3: %w", err)
	}

	fmt.Printf("Uploaded file %s to S3 bucket %s\n", file.ID, d.bucket)
	return nil
}

// DownloadAllYaml downloads all YAML files from S3.
func (d *S3Driver) DownloadAllYaml(ctx context.Context) ([]entity.YamlFile, error) {
	if d.client == nil {
		return nil, errors.New("s3 client not initialized")
	}
	output, err := d.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(d.bucket),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list S3 objects: %w", err)
	}

	var files []entity.YamlFile
	for _, obj := range output.Contents {
		files = append(files, entity.YamlFile{
			ID: *obj.Key,
		})
	}
	return files, nil
}
