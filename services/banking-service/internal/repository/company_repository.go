package repository

import (
	"banking-service/internal/model"
	"context"
)

type CompanyRepository interface {
	Create(ctx context.Context, company *model.Company) error
	WorkCodeExists(ctx context.Context, id uint) (bool, error)
	RegistrationNumberExists(ctx context.Context, registrationNumber string) (bool, error)
	TaxNumberExists(ctx context.Context, taxNumber string) (bool, error)
}
