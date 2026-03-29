package dto

import "time"

type BaseListingResponse struct {
	ListingID         uint    `json:"listingId"`
	Ticker            string  `json:"ticker"`
	Name              string  `json:"name"`
	Exchange          string  `json:"exchange"`
	Price             float64 `json:"price"`
	Ask               float64 `json:"ask"`
	Bid               float64 `json:"bid"`
	Change            float64 `json:"change"`
	Volume            uint    `json:"volume"`
	MaintenanceMargin float64 `json:"maintenanceMargin"`
	InitialMarginCost float64 `json:"initialMarginCost"`
}

type StockResponse struct {
	BaseListingResponse
	OutstandingShares float64 `json:"outstandingShares"`
	DividendYield     float64 `json:"dividendYield"`
}

type FuturesResponse struct {
	BaseListingResponse
	SettlementDate time.Time `json:"settlementDate"`
	ContractSize   float64   `json:"contractSize"`
	ContractUnit   string    `json:"contractUnit"`
}

type ForexResponse struct {
	ForexPairID       uint    `json:"forexPairId"`
	Ticker            string  `json:"ticker"`
	Base              string  `json:"base"`
	Quote             string  `json:"quote"`
	Price             float64 `json:"price"`
	Ask               float64 `json:"ask"`
	Bid               float64 `json:"bid"`
	Change            float64 `json:"change"`
	Volume            uint    `json:"volume"`
	MaintenanceMargin float64 `json:"maintenanceMargin"`
	InitialMarginCost float64 `json:"initialMarginCost"`
}

type PaginatedResponse[T any] struct {
	Data     []T   `json:"data"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
}
