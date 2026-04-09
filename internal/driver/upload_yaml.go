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

	fmt.Printf("Uploaded file %s to S3 bucket %s", file.ID, d.bucket)
	return nil
}
