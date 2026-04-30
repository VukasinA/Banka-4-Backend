package repository

import (
	"context"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
)

type AssetOwnershipRepository interface {
	FindByIdentity(ctx context.Context, identityID uint, ownerType model.OwnerType) ([]model.AssetOwnership, error)
	Upsert(ctx context.Context, ownership *model.AssetOwnership) error
	// IncreaseReservedAmount atomically adds delta to reserved_amount for the
	// ownership row identified by identityID+ownerType+assetID. It is the ONLY
	// code path that may grow reserved_amount — it is called exclusively when
	// an OTC deal (option contract) is finalised.
	IncreaseReservedAmount(ctx context.Context, identityID uint, ownerType model.OwnerType, assetID uint, delta float64) error
}
