package repository

import (
	"banking-service/internal/model"
	"context"
)

type CardRequestRepository interface {
	Create(ctx context.Context, request *model.CardRequest) error
	FindByAccountNumberClientIDAndCode(ctx context.Context, accountNumber string, clientID uint, code string) (*model.CardRequest, error)
	FindLatestPendingByAccountNumberAndClientID(ctx context.Context, accountNumber string, clientID uint) (*model.CardRequest, error)
	Update(ctx context.Context, request *model.CardRequest) error
}
