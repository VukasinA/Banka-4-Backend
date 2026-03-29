package service

import (
	"context"

	commonErrors "github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/errors"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/dto"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/repository"
)

type ListingService struct {
	listingRepo repository.ListingRepository
	futuresRepo repository.FuturesRepository
	forexRepo   repository.ForexRepository
}

func NewListingService(
	listingRepo repository.ListingRepository,
	futuresRepo repository.FuturesRepository,
	forexRepo repository.ForexRepository,
) *ListingService {
	return &ListingService{
		listingRepo: listingRepo,
		futuresRepo: futuresRepo,
		forexRepo:   forexRepo,
	}
}

// --- Helpers ---

func latestDaily(infos []model.ListingDailyPriceInfo) *model.ListingDailyPriceInfo {
	if len(infos) == 0 {
		return nil
	}
	return &infos[0]
}

func baseResponse(l model.Listing, daily *model.ListingDailyPriceInfo) dto.BaseListingResponse {
	r := dto.BaseListingResponse{
		ListingID:         l.ListingID,
		Ticker:            l.Ticker,
		Name:              l.Name,
		Exchange:          l.ExchangeMIC,
		Price:             l.Price,
		Ask:               l.Ask,
		MaintenanceMargin: l.MaintenanceMargin,
		InitialMarginCost: l.MaintenanceMargin * 1.1,
	}
	if daily != nil {
		r.Bid = daily.Bid
		r.Change = daily.Change
		r.Volume = daily.Volume
	}
	return r
}

func toFilter(q dto.ListingQuery) (repository.ListingFilter, error) {
	f := repository.ListingFilter{
		Search:    q.Search,
		Exchange:  q.Exchange,
		PriceMin:  q.PriceMin,
		PriceMax:  q.PriceMax,
		AskMin:    q.AskMin,
		AskMax:    q.AskMax,
		BidMin:    q.BidMin,
		BidMax:    q.BidMax,
		VolumeMin: q.VolumeMin,
		VolumeMax: q.VolumeMax,
		SortBy:    q.SortBy,
		SortDir:   q.SortDir,
		Page:      q.Page,
		PageSize:  q.PageSize,
	}
	sd, err := q.ParseSettlementDate()
	if err != nil {
		return f, err
	}
	f.SettlementDate = sd
	return f, nil
}

// --- Stocks ---

func (s *ListingService) GetStocks(ctx context.Context, q dto.ListingQuery) (dto.PaginatedResponse[dto.StockResponse], error) {
	filter, err := toFilter(q)
	if err != nil {
		return dto.PaginatedResponse[dto.StockResponse]{}, commonErrors.BadRequestErr("invalid settlement_date format")
	}

	listings, total, err := s.listingRepo.FindStocks(ctx, filter)
	if err != nil {
		return dto.PaginatedResponse[dto.StockResponse]{}, commonErrors.InternalErr(err)
	}

	data := make([]dto.StockResponse, len(listings))
	for i, l := range listings {
		daily := latestDaily(l.DailyPriceInfos)
		var outstandingShares, dividendYield float64
		if l.Stock != nil {
			outstandingShares = l.Stock.OutstandingShares
			dividendYield = l.Stock.DividendYield
		}
		data[i] = dto.StockResponse{
			BaseListingResponse: baseResponse(l, daily),
			OutstandingShares:   outstandingShares,
			DividendYield:       dividendYield,
		}
	}

	return dto.PaginatedResponse[dto.StockResponse]{
		Data:     data,
		Total:    total,
		Page:     q.Page,
		PageSize: q.PageSize,
	}, nil
}

// --- Futures ---

func (s *ListingService) GetFutures(ctx context.Context, q dto.ListingQuery) (dto.PaginatedResponse[dto.FuturesResponse], error) {
	filter, err := toFilter(q)
	if err != nil {
		return dto.PaginatedResponse[dto.FuturesResponse]{}, commonErrors.BadRequestErr("invalid settlement_date format")
	}

	listings, total, err := s.listingRepo.FindFutures(ctx, filter)
	if err != nil {
		return dto.PaginatedResponse[dto.FuturesResponse]{}, commonErrors.InternalErr(err)
	}

	// batch fetch futures contracts po tickerima
	tickers := make([]string, len(listings))
	for i, l := range listings {
		tickers[i] = l.Ticker
	}

	contracts, err := s.futuresRepo.FindByTickers(ctx, tickers)
	if err != nil {
		return dto.PaginatedResponse[dto.FuturesResponse]{}, commonErrors.InternalErr(err)
	}

	contractMap := make(map[string]model.FuturesContract)
	for _, fc := range contracts {
		contractMap[fc.Ticker] = fc
	}

	data := make([]dto.FuturesResponse, len(listings))
	for i, l := range listings {
		daily := latestDaily(l.DailyPriceInfos)
		fc := contractMap[l.Ticker]
		data[i] = dto.FuturesResponse{
			BaseListingResponse: baseResponse(l, daily),
			SettlementDate:      fc.SettlementDate,
			ContractSize:        fc.ContractSize,
			ContractUnit:        fc.ContractUnit,
		}
	}

	return dto.PaginatedResponse[dto.FuturesResponse]{
		Data:     data,
		Total:    total,
		Page:     q.Page,
		PageSize: q.PageSize,
	}, nil
}

// --- Forex ---

func (s *ListingService) GetForex(ctx context.Context, q dto.ListingQuery) (dto.PaginatedResponse[dto.ForexResponse], error) {
	filter, _ := toFilter(q)

	pairs, total, err := s.forexRepo.FindAll(ctx, filter)
	if err != nil {
		return dto.PaginatedResponse[dto.ForexResponse]{}, commonErrors.InternalErr(err)
	}

	data := make([]dto.ForexResponse, len(pairs))
	for i, p := range pairs {
		data[i] = dto.ForexResponse{
			ForexPairID: p.ForexPairID,
			Ticker:      p.Base + "/" + p.Quote,
			Base:        p.Base,
			Quote:       p.Quote,
			Price:       p.Rate,
		}
	}

	return dto.PaginatedResponse[dto.ForexResponse]{
		Data:     data,
		Total:    total,
		Page:     q.Page,
		PageSize: q.PageSize,
	}, nil
}
