package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/pb"
)

type TradingServiceClient struct {
	client pb.TradingServiceClient
}

func NewTradingServiceClient(conn *grpc.ClientConn) *TradingServiceClient {
	return &TradingServiceClient{client: pb.NewTradingServiceClient(conn)}
}

func (c *TradingServiceClient) TransferFunds(ctx context.Context, fromManagerID uint, toManagerID uint) (uint64, error) {
	resp, err := c.client.TransferFunds(ctx, &pb.TransferFundsRequest{
		FromManagerId: uint64(fromManagerID),
		ToManagerId:   uint64(toManagerID),
	})
	if err != nil {
		return 0, fmt.Errorf("trading client TransferFunds: %w", err)
	}
	return resp.FundsTransferred, nil
}
