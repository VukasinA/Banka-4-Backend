package model

import "time"

type InvestmentFund struct {
	FundID              uint      `gorm:"primaryKey;autoIncrement"`
	Name                string    `gorm:"not null;size:255;uniqueIndex"`
	Description         string    `gorm:"not null;size:1000"`
	MinimumContribution float64   `gorm:"not null;default:0"`
	ManagerID           uint      `gorm:"not null"`
	LiquidAssets        float64   `gorm:"not null;default:0"`
	AccountNumber       string    `gorm:"not null;size:50;uniqueIndex"`
	CreatedAt           time.Time `gorm:"not null"`

	Listings  []Listing            `gorm:"many2many:fund_listings;"`
	Positions []ClientFundPosition `gorm:"foreignKey:FundID"`
}

type ClientFundPosition struct {
	PositionID          uint            `gorm:"primaryKey;autoIncrement"`
	FundID              uint            `gorm:"not null;index"`
	Fund                *InvestmentFund `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	ClientID            uint            `gorm:"not null;index"`
	TotalInvestedAmount float64         `gorm:"not null;default:0"`
	UpdatedAt           time.Time       `gorm:"not null"`
}
