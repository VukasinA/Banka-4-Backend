package service

import (
	"context"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/client"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/dto"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/repository"
)

type ProfitService struct {
	repo       repository.ProfitRepository
	userClient client.UserServiceClient
}

func NewProfitService(
	repo repository.ProfitRepository,
	userClient client.UserServiceClient,
) *ProfitService {
	return &ProfitService{
		repo:       repo,
		userClient: userClient,
	}
}
func (s *ProfitService) GetActuaryProfits(
	ctx context.Context,
	page, pageSize int32,
	firstName, lastName string,
) ([]dto.ActuaryProfitResponse, error) {

	// 1. user-service (identitet)
	resp, err := s.userClient.GetAllActuaries(ctx, page, pageSize, firstName, lastName)
	if err != nil {
		return nil, err
	}

	// 2. izvuci IDs
	ids := make([]uint64, 0, len(resp.Actuaries))
	for _, a := range resp.Actuaries {
		ids = append(ids, a.Id)
	}

	// 3. profit iz trading DB
	profitMap, err := s.repo.GetProfitByUserIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	result := make([]dto.ActuaryProfitResponse, 0, len(resp.Actuaries))

	for _, a := range resp.Actuaries {

		profit := profitMap[a.Id]

		result = append(result, dto.ActuaryProfitResponse{
			FirstName: a.FirstName,
			LastName:  a.LastName,
			Role:      "supervisor",
			ProfitRSD: profit,
		})
	}

	return result, nil
}
func (s *ProfitService) GetFundPositions(
	ctx context.Context,
) ([]dto.FundPositionResponse, error) {

	funds, err := s.repo.GetAllInvestmentFunds(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]dto.FundPositionResponse, 0, len(funds))

	for _, f := range funds {

		// sum invested amounts
		var total float64
		for _, p := range f.Positions {
			total += p.TotalInvestedAmount
		}

		bankPct := 10.0
		bankValue := total * bankPct / 100
		profit := bankValue * 0.2

		managerName := ""
		if manager, err := s.userClient.GetEmployeeById(ctx, uint64(f.ManagerID)); err == nil {
			managerName = manager.FullName
		}

		result = append(result, dto.FundPositionResponse{
			FundName:       f.Name,
			ManagerName:    managerName,
			BankSharePct:   bankPct,
			BankShareValue: bankValue,
			Profit:         profit,
		})
	}

	return result, nil
}
