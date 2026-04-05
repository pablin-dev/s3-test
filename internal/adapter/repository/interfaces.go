package repository

import (
	"context"
	"github.com/pablodev/s3-test/internal/entity"
)

type S3Repository interface {
	UploadYaml(ctx context.Context, file entity.YamlFile) error
	DownloadAllYaml(ctx context.Context) ([]entity.YamlFile, error)
}

type DynamoDBRepository interface {
	StoreResult(ctx context.Context, id string, version int, result interface{}) error
}
