package service

import (
	"errors"
	"user-service/internal/dto"
	"user-service/internal/model"
	"user-service/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type EmployeeService struct {
	repo repository.EmployeeRepository
}

func NewEmployeeService(repo repository.EmployeeRepository) *EmployeeService {
	return &EmployeeService{repo: repo}
}

func (s *EmployeeService) Register(dto dto.UserCreateDTO) (*model.Employee, error) {

	if _, err := s.repo.GetByEmail(dto.Email); err == nil {
		return nil, errors.New("email already in use")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	employee := model.Employee{
		FirstName:   dto.FirstName,
		LastName:    dto.LastName,
		Gender:      dto.Gender,
		DateOfBirth: dto.DateOfBirth,
		Email:       dto.Email,
		PhoneNumber: dto.PhoneNumber,
		Address:     dto.Address,
		Username:    dto.Username,
		Password:    string(hashedPassword),
		Department:  dto.Department,
		PositionID:  dto.PositionID,
		Active:      true,
	}

	if err := s.repo.Create(&employee); err != nil {
		return nil, err
	}

	return &employee, nil
}
