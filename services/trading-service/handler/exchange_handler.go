package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/dto"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/service"
)

type ExchangeHandler struct {
	service *service.ExchangeService
}

func NewExchangeHandler(service *service.ExchangeService) *ExchangeHandler {
	return &ExchangeHandler{service: service}
}

// GetAll godoc
// @Summary Get all exchanges
// @Description Returns a list of all stock exchanges
// @Tags exchange
// @Produce json
// @Success 200 {array} dto.ExchangeResponse
// @Failure 500 {object} map[string]string
// @Router /api/exchange [get]
func (h *ExchangeHandler) GetAll(c *gin.Context) {
	exchanges, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToExchangeResponseList(exchanges))
}

// ToggleTradingEnabled godoc
// @Summary Toggle trading enabled for an exchange
// @Description Enables or disables trading time enforcement for a specific exchange (for testing purposes)
// @Tags exchange
// @Produce json
// @Param micCode path string true "Exchange MIC code"
// @Success 200 {object} dto.ExchangeResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/exchange/{micCode}/toggle [patch]
func (h *ExchangeHandler) ToggleTradingEnabled(c *gin.Context) {
	micCode := c.Param("micCode")

	exchange, err := h.service.ToggleTradingEnabled(c.Request.Context(), micCode)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToExchangeResponse(*exchange))
}
