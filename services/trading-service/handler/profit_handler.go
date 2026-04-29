package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/service"
)

type ProfitHandler struct {
	service *service.ProfitService
}

func NewProfitHandler(service *service.ProfitService) *ProfitHandler {
	return &ProfitHandler{service: service}
}

// GetActuaryProfits godoc
// @Summary Get actuary profits
// @Description Returns paginated list of actuaries with their profits (agents and supervisors)
// @Tags profit
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param first_name query string false "Filter by first name"
// @Param last_name query string false "Filter by last name"
// @Success 200 {array} dto.ActuaryProfitResponse
// @Failure 400 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /api/profit/actuaries [get]
func (h *ProfitHandler) GetActuaryProfits(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page_size"})
		return
	}

	firstName := c.Query("first_name")
	lastName := c.Query("last_name")

	res, err := h.service.GetActuaryProfits(
		c.Request.Context(),
		int32(page),
		int32(pageSize),
		firstName,
		lastName,
	)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetFundPositions godoc
// @Summary Get investment fund positions
// @Description Returns all investment funds with bank share, manager info and profit calculation
// @Tags profit
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {array} dto.FundPositionResponse
// @Failure 500 {object} errors.AppError
// @Router /api/profit/funds [get]
func (h *ProfitHandler) GetFundPositions(c *gin.Context) {
	res, err := h.service.GetFundPositions(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}