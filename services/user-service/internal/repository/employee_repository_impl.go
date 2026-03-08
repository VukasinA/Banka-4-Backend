package repository

import (
	"user-service/internal/model"

	"gorm.io/gorm"
)

type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) Create(employee *model.Employee) error {
	return r.db.Create(employee).Error
}

func (r *employeeRepository) GetByEmail(email string) (*model.Employee, error) {
	var employee model.Employee

	err := r.db.
		Where("email = ?", email).
		First(&employee).Error

	if err != nil {
		return nil, err
	}

	return &employee, nil
}
