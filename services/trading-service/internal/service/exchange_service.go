package service

import (
	"context"

	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/errors"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/repository"
)

type ExchangeService struct {
	repo *repository.ExchangeRepository
}

func NewExchangeService(repo *repository.ExchangeRepository) *ExchangeService {
	return &ExchangeService{repo: repo}
}

func (s *ExchangeService) GetAll(ctx context.Context) ([]model.Exchange, error) {
	exchanges, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, errors.InternalErr(err)
	}
	return exchanges, nil
}

func (s *ExchangeService) ToggleTradingEnabled(ctx context.Context, micCode string) (*model.Exchange, error) {
	exchange, err := s.repo.ToggleTradingEnabled(ctx, micCode)
	if err != nil {
		return nil, errors.InternalErr(err)
	}
	if exchange == nil {
		return nil, errors.NotFoundErr("exchange not found: " + micCode)
	}
	return exchange, nil
}
