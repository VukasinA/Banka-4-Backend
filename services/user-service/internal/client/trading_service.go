package client

import (
	"context"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/user-service/internal/config"
)

func NewTradingServiceConnection(lc fx.Lifecycle, cfg *config.Configuration) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		cfg.TradingServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return conn.Close()
		},
	})

	return conn, nil
}
