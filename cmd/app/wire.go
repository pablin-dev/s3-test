//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/google/wire"
	"github.com/pablodev/s3-test/internal/adapter/repository"
	"github.com/pablodev/s3-test/internal/driver"
	"github.com/pablodev/s3-test/internal/usecase"
)

var superSet = wire.NewSet(
	LoadConfig,
	ProvideAWSConfig,
	ProvideS3Bucket,
	ProvideDynamoDBTable,
	driver.NewS3Driver,
	driver.NewDynamoDBDriver,
	wire.Bind(new(repository.S3Repository), new(*driver.S3Driver)),
	wire.Bind(new(repository.DynamoDBRepository), new(*driver.DynamoDBDriver)),
	usecase.NewStorageUseCase,
)

func InitializeApp(ctx context.Context) (*usecase.StorageUseCase, error) {
	wire.Build(
		superSet,
		wire.Value("config.local.yaml"),
	)
	return nil, nil
}
