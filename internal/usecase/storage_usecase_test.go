package usecase

import (
	"context"
	"github.com/pablodev/s3-test/internal/adapter/repository/mocks"
	"github.com/pablodev/s3-test/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestStorageUseCase_ProcessAndSync(t *testing.T) {
	mockS3 := new(mocks.S3Repository)
	mockDB := new(mocks.DynamoDBRepository)
	uc := NewStorageUseCase(mockS3, mockDB)

	files := []entity.YamlFile{
		{ID: "1", Version: 1, Expression: "1+1"},
		{ID: "1", Version: 2, Expression: "1+2"},
	}

	mockS3.On("DownloadAllYaml", mock.Anything).Return(files, nil)
	mockDB.On("StoreResult", mock.Anything, "1", 2, mock.Anything).Return(nil)

	err := uc.ProcessAndSync(context.Background())
	assert.NoError(t, err)
	mockS3.AssertExpectations(t)
	mockDB.AssertExpectations(t)
}
