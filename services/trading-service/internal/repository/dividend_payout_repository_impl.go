package repository

import (
	"context"

	commondb "github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/db"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
	"gorm.io/gorm"
)

type dividendPayoutRepository struct {
	db *gorm.DB
}

func NewDividendPayoutRepository(db *gorm.DB) DividendPayoutRepository {
	return &dividendPayoutRepository{db: db}
}

func (r *dividendPayoutRepository) Save(ctx context.Context, payout *model.DividendPayout) error {
	return commondb.DBFromContext(ctx, r.db).Create(payout).Error
}

func (r *dividendPayoutRepository) FindAllByUserID(ctx context.Context, userID uint, ownerType model.OwnerType) ([]model.DividendPayout, error) {
	var payouts []model.DividendPayout
	err := commondb.DBFromContext(ctx, r.db).
		Where("user_id = ? AND owner_type = ?", userID, ownerType).
		Preload("Stock").
		Preload("Stock.Asset").
		Order("payment_date DESC").
		Find(&payouts).Error
	return payouts, err
}

func (r *dividendPayoutRepository) FindAll(ctx context.Context) ([]model.DividendPayout, error) {
	var payouts []model.DividendPayout
	err := commondb.DBFromContext(ctx, r.db).
		Preload("Stock").
		Preload("Stock.Asset").
		Order("payment_date DESC").
		Find(&payouts).Error
	return payouts, err
}

func (r *dividendPayoutRepository) FindAllByStockID(ctx context.Context, stockID uint) ([]model.DividendPayout, error) {
	var payouts []model.DividendPayout
	err := commondb.DBFromContext(ctx, r.db).
		Where("stock_id = ?", stockID).
		Preload("Stock").
		Preload("Stock.Asset").
		Order("payment_date DESC").
		Find(&payouts).Error
	return payouts, err
}
