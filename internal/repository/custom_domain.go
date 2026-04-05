package repository

import "github.com/pablodev/s3-test/internal/entity"

// CustomDomainRepository defines the persistence operations for CustomDomain.
type CustomDomainRepository interface {
	Save(domain *entity.CustomDomain) error
	GetByID(id string) (*entity.CustomDomain, error)
}
