package grpc

import (
	"context"

	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/pb"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/service"
)

type TradingServiceServer struct {
	pb.UnimplementedTradingServiceServer
	investmentFundService *service.InvestmentFundService
}

func NewTradingServiceServer(investmentFundService *service.InvestmentFundService) *TradingServiceServer {
	return &TradingServiceServer{
		investmentFundService: investmentFundService,
	}
}

func (s *TradingServiceServer) TransferFunds(ctx context.Context, req *pb.TransferFundsRequest) (*pb.TransferFundsResponse, error) {
	count, err := s.investmentFundService.TransferFunds(ctx, uint(req.FromManagerId), uint(req.ToManagerId))
	if err != nil {
		return nil, err
	}
	return &pb.TransferFundsResponse{FundsTransferred: uint64(count)}, nil
}
