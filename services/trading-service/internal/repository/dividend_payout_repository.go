package repository

import (
	"context"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
)

type DividendPayoutRepository interface {
	Save(ctx context.Context, payout *model.DividendPayout) error

	FindAllByUserID(ctx context.Context, userID uint, ownerType model.OwnerType) ([]model.DividendPayout, error)

	FindAll(ctx context.Context) ([]model.DividendPayout, error)

	FindAllByStockID(ctx context.Context, stockID uint) ([]model.DividendPayout, error)
}
