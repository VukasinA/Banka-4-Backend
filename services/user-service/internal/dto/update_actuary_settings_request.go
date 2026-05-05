package dto

type UpdateActuarySettingsRequest struct {
	Limit        *float64 `json:"limit" binding:"omitempty,gte=0"`
	NeedApproval *bool    `json:"need_approval"`
	IsAgent      *bool    `json:"is_agent"`
	IsSupervisor *bool    `json:"is_supervisor"`
}
