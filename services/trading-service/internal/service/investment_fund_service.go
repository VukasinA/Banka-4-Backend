package service

import (
	"context"
	"time"

	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/auth"
	commonErrors "github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/errors"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/client"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/dto"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/repository"
)

type InvestmentFundService struct {
	fundRepo      repository.InvestmentFundRepository
	ownershipRepo repository.AssetOwnershipRepository
	listingRepo   repository.ListingRepository
	bankingClient client.BankingClient
}

func NewInvestmentFundService(
	fundRepo repository.InvestmentFundRepository,
	ownershipRepo repository.AssetOwnershipRepository,
	listingRepo repository.ListingRepository,
	bankingClient client.BankingClient,
) *InvestmentFundService {
	return &InvestmentFundService{
		fundRepo:      fundRepo,
		ownershipRepo: ownershipRepo,
		listingRepo:   listingRepo,
		bankingClient: bankingClient,
	}
}

func (s *InvestmentFundService) sumSecuritiesValue(ctx context.Context, fundID uint) (float64, error) {
	ownerships, err := s.ownershipRepo.FindByUserId(ctx, fundID, model.OwnerTypeFund)
	if err != nil {
		return 0, err
	}
	if len(ownerships) == 0 {
		return 0, nil
	}

	assetIDs := make([]uint, len(ownerships))
	for i, o := range ownerships {
		assetIDs[i] = o.AssetID
	}

	listings, err := s.listingRepo.FindByAssetIDs(ctx, assetIDs)
	if err != nil {
		return 0, err
	}

	priceByAsset := make(map[uint]float64, len(listings))
	for _, l := range listings {
		priceByAsset[l.AssetID] = l.Price
	}

	var total float64
	for _, o := range ownerships {
		total += o.Amount * priceByAsset[o.AssetID]
	}
	return total, nil
}

func (s *InvestmentFundService) GetAllFunds(ctx context.Context, query dto.ListFundsQuery) (*dto.ListFundsResponse, error) {
	funds, total, err := s.fundRepo.FindAll(ctx, query.Name, query.SortBy, query.SortDir, query.Page, query.PageSize)
	if err != nil {
		return nil, commonErrors.InternalErr(err)
	}

	result := make([]dto.FundSummaryResponse, len(funds))
	for i, fund := range funds {
		secVal, err := s.sumSecuritiesValue(ctx, fund.FundID)
		if err != nil {
			return nil, commonErrors.InternalErr(err)
		}
		result[i] = dto.ToFundSummaryResponse(fund, secVal)
	}

	return &dto.ListFundsResponse{
		Data:     result,
		Total:    total,
		Page:     query.Page,
		PageSize: query.PageSize,
	}, nil
}

func (s *InvestmentFundService) GetActuaryFunds(ctx context.Context, managerID uint) ([]dto.ActuaryFundResponse, error) {
	funds, err := s.fundRepo.FindByManagerID(ctx, managerID)
	if err != nil {
		return nil, commonErrors.InternalErr(err)
	}

	result := make([]dto.ActuaryFundResponse, len(funds))
	for i, fund := range funds {
		secVal, err := s.sumSecuritiesValue(ctx, fund.FundID)
		if err != nil {
			return nil, commonErrors.InternalErr(err)
		}
		result[i] = dto.ToActuaryFundResponse(fund, secVal)
	}

	return result, nil
}

func (s *InvestmentFundService) CreateFund(ctx context.Context, req dto.CreateFundRequest) (*dto.CreateFundResponse, error) {
	authCtx := auth.GetAuthFromContext(ctx)
	if authCtx == nil {
		return nil, commonErrors.UnauthorizedErr("not authenticated")
	}

	if authCtx.IdentityType != auth.IdentityEmployee {
		return nil, commonErrors.ForbiddenErr("only employees can create investment funds")
	}

	if authCtx.EmployeeID == nil {
		return nil, commonErrors.UnauthorizedErr("employee identity missing")
	}

	managerID := *authCtx.EmployeeID

	existing, err := s.fundRepo.FindByName(ctx, req.Name)
	if err != nil {
		return nil, commonErrors.InternalErr(err)
	}
	if existing != nil {
		return nil, commonErrors.ConflictErr("fund name is already taken")
	}

	accountNumber, err := s.bankingClient.CreateFundAccount(ctx, req.Name, uint64(managerID))
	if err != nil {
		return nil, commonErrors.InternalErr(err)
	}

	fund := &model.InvestmentFund{
		Name:                req.Name,
		Description:         req.Description,
		MinimumContribution: req.MinimumContribution,
		ManagerID:           managerID,
		LiquidAssets:        0,
		AccountNumber:       accountNumber,
		CreatedAt:           time.Now(),
	}

	if err := s.fundRepo.Create(ctx, fund); err != nil {
		return nil, commonErrors.InternalErr(err)
	}

	return &dto.CreateFundResponse{
		FundID:              fund.FundID,
		Name:                fund.Name,
		Description:         fund.Description,
		MinimumContribution: fund.MinimumContribution,
		ManagerID:           fund.ManagerID,
		LiquidAssets:        fund.LiquidAssets,
		AccountNumber:       fund.AccountNumber,
		CreatedAt:           fund.CreatedAt,
	}, nil
}
