package dto

type CreateFundRequest struct {
	Name                string  `json:"name"`
	Description         string  `json:"description"`
	MinimumContribution float64 `json:"minimumContribution"`
}
