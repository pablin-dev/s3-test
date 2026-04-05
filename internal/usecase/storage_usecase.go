package usecase

import (
	"context"
	"github.com/expr-lang/expr"
	"github.com/pablodev/s3-test/internal/adapter/repository"
	"github.com/pablodev/s3-test/internal/entity"
)

type StorageUseCase struct {
	s3 repository.S3Repository
	db repository.DynamoDBRepository
}

func NewStorageUseCase(s3 repository.S3Repository, db repository.DynamoDBRepository) *StorageUseCase {
	return &StorageUseCase{s3: s3, db: db}
}

func (u *StorageUseCase) Upload(ctx context.Context, file entity.YamlFile) error {
	return u.s3.UploadYaml(ctx, file)
}

func (u *StorageUseCase) ProcessAndSync(ctx context.Context) error {
	files, err := u.s3.DownloadAllYaml(ctx)
	if err != nil {
		return err
	}

	latest := make(map[string]entity.YamlFile)
	for _, f := range files {
		if cur, ok := latest[f.ID]; !ok || f.Version > cur.Version {
			latest[f.ID] = f
		}
	}

	for _, f := range latest {
		program, err := expr.Compile(f.Expression)
		if err != nil {
			return err
		}
		output, err := expr.Run(program, nil)
		if err != nil {
			return err
		}
		if err := u.db.StoreResult(ctx, f.ID, f.Version, output); err != nil {
			return err
		}
	}
	return nil
}
