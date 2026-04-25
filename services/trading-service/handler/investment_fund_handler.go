package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/errors"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/dto"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/service"
)

type InvestmentFundHandler struct {
	service *service.InvestmentFundService
}

func NewInvestmentFundHandler(service *service.InvestmentFundService) *InvestmentFundHandler {
	return &InvestmentFundHandler{service: service}
}

// CreateFund godoc
// @Summary Create a new investment fund
// @Description Supervisor creates a new investment fund. An RSD account is automatically created for the fund.
// @Tags investment-funds
// @Accept json
// @Produce json
// @Param request body dto.CreateFundRequest true "Fund details"
// @Success 201 {object} dto.CreateFundResponse
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Router /api/investment-funds [post]
func (h *InvestmentFundHandler) CreateFund(c *gin.Context) {
	var req dto.CreateFundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequestErr(err.Error()))
		return
	}

	fund, err := h.service.CreateFund(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, fund)
}
