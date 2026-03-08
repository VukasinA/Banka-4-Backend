package repository

import "user-service/internal/model"

type EmployeeRepository interface {
	Create(employee *model.Employee) error
	GetByEmail(email string) (*model.Employee, error)
}
