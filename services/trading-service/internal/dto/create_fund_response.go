package dto

import "time"

type CreateFundResponse struct {
	FundID              uint      `json:"fund_id"`
	Name                string    `json:"name"`
	Description         string    `json:"description"`
	MinimumContribution float64   `json:"minimum_contribution"`
	ManagerID           uint      `json:"manager_id"`
	LiquidAssets        float64   `json:"liquid_assets"`
	AccountNumber       string    `json:"account_number"`
	CreatedAt           time.Time `json:"created_at"`
}
