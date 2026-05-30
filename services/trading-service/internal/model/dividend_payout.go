package model

import "time"

// DividendPayout represents a single dividend payment to an owner of a stock.
// Actuaries (OwnerTypeActuary) are exempt from the 15% capital gains tax —
// their dividend goes directly into bank profit.
type DividendPayout struct {
	DividendPayoutID uint      `gorm:"primaryKey;autoIncrement"`
	UserID           uint      `gorm:"not null;index"`
	OwnerType        OwnerType `gorm:"not null;size:10"`
	StockID          uint      `gorm:"not null;index"`
	Stock            Stock
	Quantity         float64   `gorm:"not null"`
	PricePerShare    float64   `gorm:"not null"`
	GrossAmount      float64   `gorm:"not null"`
	TaxAmount        float64   `gorm:"not null;default:0"`
	NetAmount        float64   `gorm:"not null"`
	CurrencyCode     string    `gorm:"not null;size:10"`
	AccountNumber    string    `gorm:"not null"`
	PaymentDate      time.Time `gorm:"not null;index"`
}
