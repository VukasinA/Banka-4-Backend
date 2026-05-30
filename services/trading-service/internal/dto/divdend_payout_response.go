package dto

import "time"

// DividendPayoutResponse represents a single dividend payout for API responses.
// Maps to the table columns described in the spec:
// User | Stock | Quantity | GrossAmount | Tax | NetAmount | PaymentDate
type DividendPayoutResponse struct {
	DividendPayoutID uint      `json:"dividendPayoutId"`
	UserID           uint      `json:"userId"`
	OwnerType        string    `json:"ownerType"`
	Stock            string    `json:"stock"` // ticker symbol e.g. "AAPL"
	Quantity         float64   `json:"quantity"`
	GrossAmount      float64   `json:"grossAmount"`
	TaxAmount        float64   `json:"taxAmount"`
	NetAmount        float64   `json:"netAmount"`
	CurrencyCode     string    `json:"currencyCode"`
	AccountNumber    string    `json:"accountNumber"`
	PaymentDate      time.Time `json:"paymentDate"`
}

// ListDividendPayoutsResponse wraps a list of dividend payouts.
type ListDividendPayoutsResponse struct {
	Data []DividendPayoutResponse `json:"data"`
}
