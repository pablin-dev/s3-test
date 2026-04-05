package usecase

import (
	"github.com/pablodev/s3-test/internal/entity"
	"github.com/pablodev/s3-test/internal/repository"
)

// CustomDomainUseCase handles business logic for custom domains.
type CustomDomainUseCase struct {
	repo repository.CustomDomainRepository
}

// NewCustomDomainUseCase creates a new instance of CustomDomainUseCase.
func NewCustomDomainUseCase(repo repository.CustomDomainRepository) *CustomDomainUseCase {
	return &CustomDomainUseCase{repo: repo}
}

// SaveDomain saves a custom domain.
func (u *CustomDomainUseCase) SaveDomain(domain *entity.CustomDomain) error {
	return u.repo.Save(domain)
}

// GetDomainByID retrieves a custom domain by ID.
func (u *CustomDomainUseCase) GetDomainByID(id string) (*entity.CustomDomain, error) {
	return u.repo.GetByID(id)
}
