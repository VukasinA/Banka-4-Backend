package client

import "context"

type TradingClient interface {
	TransferFunds(ctx context.Context, fromManagerID uint, toManagerID uint) (uint64, error)
}
