package handler

import (
	"net/http"
	"strconv"

	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/auth"
	"github.com/RAF-SI-2025/Banka-4-Backend/common/pkg/errors"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/dto"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/service"
	"github.com/gin-gonic/gin"
)

type DividendHandler struct {
	dividendService *service.DividendPayoutService
}

func NewDividendHandler(dividendService *service.DividendPayoutService) *DividendHandler {
	return &DividendHandler{dividendService: dividendService}
}

// GetAllDividendPayouts godoc
// @Summary List all dividend payouts
// @Description Returns all dividend payout records. Restricted to supervisors.
// @Tags dividends
// @Produce json
// @Success 200 {object} dto.ListDividendPayoutsResponse
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Security BearerAuth
// @Router /api/dividends [get]
func (h *DividendHandler) GetAllDividendPayouts(c *gin.Context) {
	ctx := c.Request.Context()

	payouts, err := h.dividendService.GetAllPayouts(ctx)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.ListDividendPayoutsResponse{
		Data: mapPayouts(payouts),
	})
}

// GetClientDividendPayouts godoc
// @Summary List dividend payouts for a client
// @Description Returns dividend payout history for a specific client. The client may only view their own payouts.
// @Tags dividends
// @Produce json
// @Param clientId path int true "Client ID"
// @Success 200 {object} dto.ListDividendPayoutsResponse
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Security BearerAuth
// @Router /api/client/{clientId}/dividends [get]
func (h *DividendHandler) GetClientDividendPayouts(c *gin.Context) {
	ctx := c.Request.Context()

	clientIDStr := c.Param("clientId")
	clientID, err := strconv.ParseUint(clientIDStr, 10, 64)
	if err != nil {
		_ = c.Error(errors.BadRequestErr("invalid clientId"))
		return
	}

	callerID, ok := authIdentityID(c)
	if !ok {
		_ = c.Error(errors.UnauthorizedErr("not authenticated"))
		return
	}

	if callerID != uint(clientID) {
		_ = c.Error(errors.ForbiddenErr("access denied"))
		return
	}

	payouts, err := h.dividendService.GetPayoutsForUser(ctx, uint(clientID), model.OwnerTypeClient)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.ListDividendPayoutsResponse{
		Data: mapPayouts(payouts),
	})
}

// GetActuaryDividendPayouts godoc
// @Summary List dividend payouts for an actuary
// @Description Returns dividend payout history for a specific actuary (employee). Tax is 0 for actuaries.
// @Tags dividends
// @Produce json
// @Param actId path int true "Actuary ID"
// @Success 200 {object} dto.ListDividendPayoutsResponse
// @Failure 400 {object} errors.AppError
// @Failure 401 {object} errors.AppError
// @Failure 403 {object} errors.AppError
// @Security BearerAuth
// @Router /api/actuary/{actId}/dividends [get]
func (h *DividendHandler) GetActuaryDividendPayouts(c *gin.Context) {
	ctx := c.Request.Context()

	actIDStr := c.Param("actId")
	actID, err := strconv.ParseUint(actIDStr, 10, 64)
	if err != nil {
		_ = c.Error(errors.BadRequestErr("invalid actId"))
		return
	}

	payouts, err := h.dividendService.GetPayoutsForUser(ctx, uint(actID), model.OwnerTypeActuary)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.ListDividendPayoutsResponse{
		Data: mapPayouts(payouts),
	})
}

// mapPayouts converts model slice to DTO slice.
func mapPayouts(payouts []model.DividendPayout) []dto.DividendPayoutResponse {
	result := make([]dto.DividendPayoutResponse, 0, len(payouts))
	for _, p := range payouts {
		ticker := ""
		if p.Stock.Asset.Ticker != "" {
			ticker = p.Stock.Asset.Ticker
		}
		result = append(result, dto.DividendPayoutResponse{
			DividendPayoutID: p.DividendPayoutID,
			UserID:           p.UserID,
			OwnerType:        string(p.OwnerType),
			Stock:            ticker,
			Quantity:         p.Quantity,
			GrossAmount:      p.GrossAmount,
			TaxAmount:        p.TaxAmount,
			NetAmount:        p.NetAmount,
			CurrencyCode:     p.CurrencyCode,
			AccountNumber:    p.AccountNumber,
			PaymentDate:      p.PaymentDate,
		})
	}
	return result
}

// authIdentityID extracts the authenticated user's ID from the token context.
// Used for self-access checks where needed.
func authIdentityID(c *gin.Context) (uint, bool) {
	claims := auth.GetAuthFromContext(c.Request.Context())
	if claims == nil {
		return 0, false
	}
	return uint(claims.IdentityID), true
}

func (h *DividendHandler) TriggerDividends(c *gin.Context) {
	if err := h.dividendService.ProcessDividends(c.Request.Context()); err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "dividends processed"})
}
