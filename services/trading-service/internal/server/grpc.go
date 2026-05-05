package server

import (
	"context"
	"fmt"
	"net"

	"go.uber.org/fx"
	"google.golang.org/grpc"

	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/pb"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/config"
	tradinggrpc "github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/grpc"
)

func NewGRPCServer(lc fx.Lifecycle, cfg *config.Configuration, tradingServer *tradinggrpc.TradingServiceServer) {
	grpcServer := grpc.NewServer()
	pb.RegisterTradingServiceServer(grpcServer, tradingServer)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			addr := fmt.Sprintf(":%s", cfg.GrpcPort)
			lis, err := net.Listen("tcp", addr)
			if err != nil {
				return fmt.Errorf("grpc listen on %s: %w", addr, err)
			}
			go func() {
				if err := grpcServer.Serve(lis); err != nil {
					fmt.Printf("trading gRPC server error: %v\n", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			grpcServer.GracefulStop()
			return nil
		},
	})
}
