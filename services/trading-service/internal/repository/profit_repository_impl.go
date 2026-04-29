package repository

import (
	"context"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
	"gorm.io/gorm"
)

type profitRepository struct {
	db *gorm.DB
}

func NewProfitRepository(db *gorm.DB) ProfitRepository {
	return &profitRepository{db: db}
}
func (r *profitRepository) GetAllInvestmentFunds(ctx context.Context) ([]model.InvestmentFund, error) {
	var funds []model.InvestmentFund

	err := r.db.WithContext(ctx).
		Preload("Positions").
		Find(&funds).Error

	return funds, err
}
func (r *profitRepository) GetProfitByUserIDs(
	ctx context.Context,
	userIDs []uint64,
) (map[uint64]float64, error) {

	type row struct {
		UserID uint
		Profit float64
	}

	var rows []row

	err := r.db.WithContext(ctx).
		Model(&model.OrderTransaction{}).
		Select("orders.order_owner_user_id as user_id, SUM(order_transactions.total_price - order_transactions.commission) as profit").
		Joins("JOIN orders ON orders.order_id = order_transactions.order_id").
		Where("orders.order_owner_user_id IN ?", userIDs).
		Group("orders.order_owner_user_id").
		Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	result := make(map[uint64]float64)

	for _, r := range rows {
		result[uint64(r.UserID)] = r.Profit
	}

	return result, nil
}
