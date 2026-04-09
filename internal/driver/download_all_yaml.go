package driver

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pablodev/s3-test/internal/entity"
)

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
