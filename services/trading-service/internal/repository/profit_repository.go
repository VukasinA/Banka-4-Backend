package repository

import (
	"context"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
)

type ProfitRepository interface {
	GetAllInvestmentFunds(ctx context.Context) ([]model.InvestmentFund, error)
	GetProfitByUserIDs(ctx context.Context, userIDs []uint64) (map[uint64]float64, error)
}
