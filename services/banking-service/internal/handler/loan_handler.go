package handler

import (
	"net/http"
	"strconv"

	"banking-service/internal/dto"
	"banking-service/internal/service"
	"common/pkg/errors"

	"github.com/gin-gonic/gin"
)

type LoanHandler struct {
	loanService *service.LoanService
}

func NewLoanHandler(loanService *service.LoanService) *LoanHandler {
	return &LoanHandler{loanService: loanService}
}

// SubmitLoanRequest godocExpand commentComment on line R23
// @Summary      Podnošenje zahteva za kredit
// @Description  Klijent podnosi zahtev za kredit. Vrši se validacija perioda otplate i valute, i računa se mesečna rata na osnovu marže banke.
// @Tags         loans
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateLoanRequest true "Podaci za zahtev kredita"
// @Success      201  {object}  dto.CreateLoanResponse
// @Failure      400  {object}  errors.AppError "Nevalidni podaci, valuta se ne poklapa ili los period otplate"
// @Failure      401  {object}  errors.AppError "Korisnik nije ulogovan"
// @Failure      403  {object}  errors.AppError "Račun ne pripada klijentu"
// @Failure      404  {object}  errors.AppError "Kredit nije pronađen"
// @Failure      500  {object}  errors.AppError "Greška na serveru"
// @Router       /api/loans/request [post]
// @Security     BearerAuth
func (h *LoanHandler) SubmitLoanRequest(c *gin.Context) {
	var req dto.CreateLoanRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequestErr(err.Error()))
		return
	}

	clientID, err := strconv.ParseUint(c.Param("client_id"), 10, 64)
	if err != nil {
		c.Error(errors.BadRequestErr("invalid client id"))
		return
	}

	resp, err := h.loanService.SubmitLoanRequest(c.Request.Context(), &req, uint(clientID))
	if err != nil {
		c.Error(err)
		return
	}

	// Uspešan odgovor
	c.JSON(http.StatusCreated, resp)
}

// GetLoans godoc
// @Summary      Pregled svih kredita klijenta
// @Description  Vraća listu kredita. Podržava sortiranje po iznosu.
// @Tags         loans
// @Produce      json
// @Param        sort query string false "Sortiraj po iznosu: 'asc' ili 'desc'"
// @Success      200  {array}   dto.LoanResponse
// @Router       /api/loans [get]
// @Security     BearerAuth
// @Failure      400  {object}  errors.AppError "Nevalidni podaci, valuta se ne poklapa ili los period otplate"
// @Failure      401  {object}  errors.AppError "Korisnik nije ulogovan"
// @Failure      403  {object}  errors.AppError "Račun ne pripada klijentu"
// @Failure      404  {object}  errors.AppError "Kredit nije pronađen"
// @Failure      500  {object}  errors.AppError "Greška na serveru"
func (h *LoanHandler) GetLoans(c *gin.Context) {
	clientID, err := strconv.ParseUint(c.Param("client_id"), 10, 64)
	if err != nil {
		c.Error(errors.BadRequestErr("invalid client id"))
		return
	}

	sortParam := c.Query("sort")
	sortByAmountDesc := sortParam == "desc"

	loans, err := h.loanService.GetClientLoans(c.Request.Context(), uint(clientID), sortByAmountDesc)
	if err != nil {
		c.Error(errors.InternalErr(err))
		return
	}

	c.JSON(http.StatusOK, loans)
}

// GetLoanByID godoc
// @Summary      Detalji kredita
// @Description  Vraća detaljne informacije o kreditu uključujući plan otplate (rate).
// @Tags         loans
// @Produce      json
// @Param        id   path      int  true  "ID kredita"
// @Success      200  {object}  dto.LoanDetailsResponse
// @Router       /api/loans/{id} [get]
// @Security     BearerAuth
// @Failure      400  {object}  errors.AppError "Nevalidni podaci, valuta se ne poklapa ili los period otplate"
// @Failure      401  {object}  errors.AppError "Korisnik nije ulogovan"
// @Failure      403  {object}  errors.AppError "Račun ne pripada klijentu"
// @Failure      404  {object}  errors.AppError "Kredit nije pronađen"
// @Failure      500  {object}  errors.AppError "Greška na serveru"
func (h *LoanHandler) GetLoanByID(c *gin.Context) {
	clientID, err := strconv.ParseUint(c.Param("client_id"), 10, 64)
	if err != nil {
		c.Error(errors.BadRequestErr("invalid client id"))
		return
	}

	loanID, err := strconv.ParseUint(c.Param("loan_id"), 10, 64)
	if err != nil {
		c.Error(errors.BadRequestErr("invalid loan id"))
		return
	}

	details, err := h.loanService.GetLoanDetails(c.Request.Context(), uint(clientID), uint(loanID))
	if err != nil {
		c.Error(errors.NotFoundErr(err.Error()))
		return
	}

	c.JSON(http.StatusOK, details)
}

func (h *LoanHandler) ListLoanRequests(c *gin.Context) {
	var query dto.ListLoanRequestsQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		c.Error(errors.BadRequestErr(err.Error()))
		return
	}

	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 10
	}

	requests, total, err := h.loanService.GetLoanRequests(c.Request.Context(), &query)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      requests,
		"total":     total,
		"page":      query.Page,
		"page_size": query.PageSize,
	})
}
func (h *LoanHandler) ApproveLoanRequest(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.BadRequestErr("invalid loan request id"))
		return
	}

	if err := h.loanService.ApproveLoanRequest(c.Request.Context(), uint(id)); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Loan request approved successfully"})
}

func (h *LoanHandler) RejectLoanRequest(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(errors.BadRequestErr("invalid loan request id"))
		return
	}

	if err := h.loanService.RejectLoanRequest(c.Request.Context(), uint(id)); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Loan request rejected successfully"})
}
