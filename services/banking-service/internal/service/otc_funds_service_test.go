package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/banking-service/internal/dto"
	"github.com/RAF-SI-2025/Banka-4-Backend/services/banking-service/internal/model"
)

type fakeOtcReservationRepo struct {
	reservation *model.OtcFundsReservation
	createErr   error
	findErr     error
	saveErr     error
}

func (f *fakeOtcReservationRepo) Create(_ context.Context, r *model.OtcFundsReservation) error {
	if f.createErr != nil {
		return f.createErr
	}
	f.reservation = r
	return nil
}

func (f *fakeOtcReservationRepo) FindByExecutionID(_ context.Context, _ string) (*model.OtcFundsReservation, error) {
	if f.findErr != nil {
		return nil, f.findErr
	}
	return f.reservation, nil
}

func (f *fakeOtcReservationRepo) Save(_ context.Context, r *model.OtcFundsReservation) error {
	if f.saveErr != nil {
		return f.saveErr
	}
	f.reservation = r
	return nil
}

type fakeOtcAccountRepo struct {
	accounts     map[string]*model.Account
	updateBalErr error
}

func (f *fakeOtcAccountRepo) Create(_ context.Context, _ *model.Account) error { return nil }
func (f *fakeOtcAccountRepo) AccountNumberExists(_ context.Context, _ string) (bool, error) {
	return false, nil
}
func (f *fakeOtcAccountRepo) GetByAccountNumber(_ context.Context, _ string) (*model.Account, error) {
	return nil, nil
}
func (f *fakeOtcAccountRepo) Update(_ context.Context, _ *model.Account) error { return nil }
func (f *fakeOtcAccountRepo) FindAllByClientID(_ context.Context, _ uint) ([]model.Account, error) {
	return nil, nil
}
func (f *fakeOtcAccountRepo) FindByAccountNumberAndClientID(_ context.Context, _ string, _ uint) (*model.Account, error) {
	return nil, nil
}
func (f *fakeOtcAccountRepo) NameExistsForClient(_ context.Context, _ uint, _ string, _ string) (bool, error) {
	return false, nil
}
func (f *fakeOtcAccountRepo) UpdateName(_ context.Context, _ string, _ string) error { return nil }
func (f *fakeOtcAccountRepo) UpdateLimits(_ context.Context, _ string, _ float64, _ float64) error {
	return nil
}
func (f *fakeOtcAccountRepo) FindAll(_ context.Context, _ *dto.ListAccountsQuery) ([]*model.Account, int64, error) {
	return nil, 0, nil
}
func (f *fakeOtcAccountRepo) FindByClientID(_ context.Context, _ uint) ([]model.Account, error) {
	return nil, nil
}
func (f *fakeOtcAccountRepo) FindByAccountType(_ context.Context, _ model.AccountType) (*model.Account, error) {
	return nil, nil
}
func (f *fakeOtcAccountRepo) FindByAccountNumber(_ context.Context, accountNumber string) (*model.Account, error) {
	if f.accounts != nil {
		return f.accounts[accountNumber], nil
	}
	return nil, nil
}
func (f *fakeOtcAccountRepo) UpdateBalance(_ context.Context, _ *model.Account) error {
	return f.updateBalErr
}

func newOtcFundsService(
	reservationRepo *fakeOtcReservationRepo,
	accountRepo *fakeOtcAccountRepo,
) *OtcFundsService {
	if reservationRepo == nil {
		reservationRepo = &fakeOtcReservationRepo{}
	}
	if accountRepo == nil {
		accountRepo = &fakeOtcAccountRepo{}
	}
	txManager := &fakeBankingTxManager{}
	return NewOtcFundsService(accountRepo, reservationRepo, txManager, nil)
}

func TestGetByExecutionID_Success(t *testing.T) {
	t.Parallel()

	expected := &model.OtcFundsReservation{
		ExecutionID:         "exec-1",
		BuyerAccountNumber:  "buyer-acc",
		SellerAccountNumber: "seller-acc",
		TradeAmount:         1000,
		TradeCurrencyCode:   model.RSD,
		Status:              model.OtcFundsReservationStatusReserved,
	}
	repo := &fakeOtcReservationRepo{reservation: expected}
	svc := newOtcFundsService(repo, nil)

	result, err := svc.GetByExecutionID(context.Background(), "exec-1")

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "exec-1", result.ExecutionID)
	require.Equal(t, model.OtcFundsReservationStatusReserved, result.Status)
}

func TestGetByExecutionID_RepoError(t *testing.T) {
	t.Parallel()

	repo := &fakeOtcReservationRepo{findErr: fmt.Errorf("db connection lost")}
	svc := newOtcFundsService(repo, nil)

	result, err := svc.GetByExecutionID(context.Background(), "exec-1")

	require.Error(t, err)
	require.Nil(t, result)
}

func TestGetByExecutionID_NotFound(t *testing.T) {
	t.Parallel()

	repo := &fakeOtcReservationRepo{reservation: nil}
	svc := newOtcFundsService(repo, nil)

	result, err := svc.GetByExecutionID(context.Background(), "missing-exec")

	require.Error(t, err)
	require.Nil(t, result)
}

func TestGetByExecutionID_TrimsWhitespace(t *testing.T) {
	t.Parallel()

	expected := &model.OtcFundsReservation{
		ExecutionID: "exec-1",
		Status:      model.OtcFundsReservationStatusReserved,
	}
	repo := &fakeOtcReservationRepo{reservation: expected}
	svc := newOtcFundsService(repo, nil)

	result, err := svc.GetByExecutionID(context.Background(), "  exec-1  ")

	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestRelease_EmptyExecutionID(t *testing.T) {
	t.Parallel()

	svc := newOtcFundsService(nil, nil)
	result, err := svc.Release(context.Background(), "")

	require.Error(t, err)
	require.Nil(t, result)
}

func TestRelease_ReservationNotFound(t *testing.T) {
	t.Parallel()

	repo := &fakeOtcReservationRepo{reservation: nil}
	svc := newOtcFundsService(repo, nil)

	result, err := svc.Release(context.Background(), "exec-1")

	require.Error(t, err)
	require.Nil(t, result)
}

func TestRelease_FindByExecutionIDError(t *testing.T) {
	t.Parallel()

	repo := &fakeOtcReservationRepo{findErr: fmt.Errorf("db error")}
	svc := newOtcFundsService(repo, nil)

	result, err := svc.Release(context.Background(), "exec-1")

	require.Error(t, err)
	require.Nil(t, result)
}

func TestRelease_AlreadyReleased_Idempotent(t *testing.T) {
	t.Parallel()

	reservation := &model.OtcFundsReservation{
		ExecutionID:         "exec-1",
		BuyerAccountNumber:  "buyer-acc",
		SellerAccountNumber: "seller-acc",
		Status:              model.OtcFundsReservationStatusReleased,
	}
	repo := &fakeOtcReservationRepo{reservation: reservation}
	svc := newOtcFundsService(repo, nil)

	result, err := svc.Release(context.Background(), "exec-1")

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, model.OtcFundsReservationStatusReleased, result.Status)
}

func TestRelease_CommittedStatus_Error(t *testing.T) {
	t.Parallel()

	reservation := &model.OtcFundsReservation{
		ExecutionID: "exec-1",
		Status:      model.OtcFundsReservationStatusCommitted,
	}
	repo := &fakeOtcReservationRepo{reservation: reservation}
	svc := newOtcFundsService(repo, nil)

	result, err := svc.Release(context.Background(), "exec-1")

	require.Error(t, err)
	require.Nil(t, result)
}

func TestRelease_RefundedStatus_Error(t *testing.T) {
	t.Parallel()

	reservation := &model.OtcFundsReservation{
		ExecutionID: "exec-1",
		Status:      model.OtcFundsReservationStatusRefunded,
	}
	repo := &fakeOtcReservationRepo{reservation: reservation}
	svc := newOtcFundsService(repo, nil)

	result, err := svc.Release(context.Background(), "exec-1")

	require.Error(t, err)
	require.Nil(t, result)
}

func TestRelease_UnknownStatus_Error(t *testing.T) {
	t.Parallel()

	reservation := &model.OtcFundsReservation{
		ExecutionID: "exec-1",
		Status:      model.OtcFundsReservationStatus("UNKNOWN"),
	}
	repo := &fakeOtcReservationRepo{reservation: reservation}
	svc := newOtcFundsService(repo, nil)

	result, err := svc.Release(context.Background(), "exec-1")

	require.Error(t, err)
	require.Nil(t, result)
}

func TestRelease_Reserved_BuyerAccountNotFound(t *testing.T) {
	t.Parallel()

	reservation := &model.OtcFundsReservation{
		ExecutionID:        "exec-1",
		BuyerAccountNumber: "missing-buyer",
		Status:             model.OtcFundsReservationStatusReserved,
		SourceAmount:       500,
	}
	resRepo := &fakeOtcReservationRepo{reservation: reservation}
	accRepo := &fakeOtcAccountRepo{accounts: map[string]*model.Account{}}
	svc := newOtcFundsService(resRepo, accRepo)

	result, err := svc.Release(context.Background(), "exec-1")

	require.Error(t, err)
	require.Nil(t, result)
}

func TestRelease_Reserved_UpdateBalanceError(t *testing.T) {
	t.Parallel()

	reservation := &model.OtcFundsReservation{
		ExecutionID:        "exec-1",
		BuyerAccountNumber: "buyer-acc",
		Status:             model.OtcFundsReservationStatusReserved,
		SourceAmount:       500,
	}
	resRepo := &fakeOtcReservationRepo{reservation: reservation}
	accRepo := &fakeOtcAccountRepo{
		accounts: map[string]*model.Account{
			"buyer-acc": {
				AccountNumber:    "buyer-acc",
				AvailableBalance: 10000,
				Currency:         model.Currency{Code: model.RSD},
			},
		},
		updateBalErr: fmt.Errorf("db write error"),
	}
	svc := newOtcFundsService(resRepo, accRepo)

	result, err := svc.Release(context.Background(), "exec-1")

	require.Error(t, err)
	require.Nil(t, result)
}

func TestRelease_Reserved_Success(t *testing.T) {
	t.Parallel()

	reservation := &model.OtcFundsReservation{
		ExecutionID:        "exec-1",
		BuyerAccountNumber: "buyer-acc",
		Status:             model.OtcFundsReservationStatusReserved,
		SourceAmount:       500,
	}
	resRepo := &fakeOtcReservationRepo{reservation: reservation}
	accRepo := &fakeOtcAccountRepo{
		accounts: map[string]*model.Account{
			"buyer-acc": {
				AccountNumber:    "buyer-acc",
				AvailableBalance: 10000,
				Currency:         model.Currency{Code: model.RSD},
			},
		},
	}
	svc := newOtcFundsService(resRepo, accRepo)

	result, err := svc.Release(context.Background(), "exec-1")

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, model.OtcFundsReservationStatusReleased, result.Status)
}

func TestRelease_SaveError(t *testing.T) {
	t.Parallel()

	reservation := &model.OtcFundsReservation{
		ExecutionID:        "exec-1",
		BuyerAccountNumber: "buyer-acc",
		Status:             model.OtcFundsReservationStatusReserved,
		SourceAmount:       500,
	}
	resRepo := &fakeOtcReservationRepo{
		reservation: reservation,
		saveErr:     fmt.Errorf("save failed"),
	}
	accRepo := &fakeOtcAccountRepo{
		accounts: map[string]*model.Account{
			"buyer-acc": {
				AccountNumber:    "buyer-acc",
				AvailableBalance: 10000,
				Currency:         model.Currency{Code: model.RSD},
			},
		},
	}
	svc := newOtcFundsService(resRepo, accRepo)

	result, err := svc.Release(context.Background(), "exec-1")

	require.Error(t, err)
	require.Nil(t, result)
}

func newOtcFundsServiceWithExchange(
	reservationRepo *fakeOtcReservationRepo,
	accountRepo *fakeOtcAccountRepo,
	exchangeRateRepo *fakeExchangeRateRepo,
) *OtcFundsService {
	if reservationRepo == nil {
		reservationRepo = &fakeOtcReservationRepo{}
	}
	if accountRepo == nil {
		accountRepo = &fakeOtcAccountRepo{}
	}
	txManager := &fakeBankingTxManager{}
	var exchSvc *ExchangeService
	if exchangeRateRepo != nil {
		exchSvc = &ExchangeService{repo: exchangeRateRepo}
	}
	return NewOtcFundsService(accountRepo, reservationRepo, txManager, exchSvc)
}

func TestConvertTradeAmount_SameCurrency_ReturnsAmount(t *testing.T) {
	t.Parallel()

	svc := newOtcFundsService(nil, nil)
	result, err := svc.convertTradeAmount(context.Background(), 1000, model.RSD, model.RSD)

	require.NoError(t, err)
	require.Equal(t, 1000.0, result)
}

func TestConvertTradeAmount_DifferentCurrency_Converts(t *testing.T) {
	t.Parallel()

	exchangeRateRepo := &fakeExchangeRateRepo{
		rates: []model.ExchangeRate{
			{CurrencyCode: model.EUR, BuyRate: 117.0, SellRate: 118.0},
			{CurrencyCode: model.USD, BuyRate: 105.0, SellRate: 106.0},
		},
	}
	svc := newOtcFundsServiceWithExchange(nil, nil, exchangeRateRepo)

	result, err := svc.convertTradeAmount(context.Background(), 100, model.EUR, model.USD)

	require.NoError(t, err)
	require.Greater(t, result, 0.0)
}

func TestConvertTradeAmount_ExchangeServiceError_ReturnsError(t *testing.T) {
	t.Parallel()

	exchangeRateRepo := &fakeExchangeRateRepo{
		err: fmt.Errorf("db error"),
	}
	svc := newOtcFundsServiceWithExchange(nil, nil, exchangeRateRepo)

	result, err := svc.convertTradeAmount(context.Background(), 100, model.EUR, model.USD)

	require.Error(t, err)
	require.Equal(t, 0.0, result)
}

func TestReserve_EmptyExecutionID(t *testing.T) {
	t.Parallel()

	svc := newOtcFundsService(nil, nil)
	result, err := svc.Reserve(context.Background(), "", "buyer", "seller", 100, model.RSD)

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "execution id is required")
}

func TestReserve_EmptyBuyerAccount(t *testing.T) {
	t.Parallel()

	svc := newOtcFundsService(nil, nil)
	result, err := svc.Reserve(context.Background(), "exec-1", "", "seller", 100, model.RSD)

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "buyer and seller account numbers are required")
}

func TestReserve_EmptySellerAccount(t *testing.T) {
	t.Parallel()

	svc := newOtcFundsService(nil, nil)
	result, err := svc.Reserve(context.Background(), "exec-1", "buyer", "", 100, model.RSD)

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "buyer and seller account numbers are required")
}

func TestReserve_SameBuyerAndSeller(t *testing.T) {
	t.Parallel()

	svc := newOtcFundsService(nil, nil)
	result, err := svc.Reserve(context.Background(), "exec-1", "same-acc", "same-acc", 100, model.RSD)

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "buyer and seller accounts must be different")
}

func TestReserve_ZeroAmount(t *testing.T) {
	t.Parallel()

	svc := newOtcFundsService(nil, nil)
	result, err := svc.Reserve(context.Background(), "exec-1", "buyer", "seller", 0, model.RSD)

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "amount must be greater than zero")
}

func TestReserve_NegativeAmount(t *testing.T) {
	t.Parallel()

	svc := newOtcFundsService(nil, nil)
	result, err := svc.Reserve(context.Background(), "exec-1", "buyer", "seller", -100, model.RSD)

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "amount must be greater than zero")
}

func TestReserve_UnsupportedCurrency(t *testing.T) {
	t.Parallel()

	svc := newOtcFundsService(nil, nil)
	result, err := svc.Reserve(context.Background(), "exec-1", "buyer", "seller", 100, model.CurrencyCode("XYZ"))

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "unsupported trade currency")
}

func TestReserve_FindByExecutionIDError(t *testing.T) {
	t.Parallel()

	resRepo := &fakeOtcReservationRepo{findErr: fmt.Errorf("db error")}
	svc := newOtcFundsService(resRepo, nil)

	result, err := svc.Reserve(context.Background(), "exec-1", "buyer", "seller", 100, model.RSD)

	require.Error(t, err)
	require.Nil(t, result)
}

func TestReserve_IdempotentWithSameParams(t *testing.T) {
	t.Parallel()

	existing := &model.OtcFundsReservation{
		ExecutionID:         "exec-1",
		BuyerAccountNumber:  "buyer",
		SellerAccountNumber: "seller",
		TradeAmount:         100,
		TradeCurrencyCode:   model.RSD,
		Status:              model.OtcFundsReservationStatusReserved,
	}
	resRepo := &fakeOtcReservationRepo{reservation: existing}
	svc := newOtcFundsService(resRepo, nil)

	result, err := svc.Reserve(context.Background(), "exec-1", "buyer", "seller", 100, model.RSD)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, model.OtcFundsReservationStatusReserved, result.Status)
}

func TestReserve_ConflictWithDifferentParams(t *testing.T) {
	t.Parallel()

	existing := &model.OtcFundsReservation{
		ExecutionID:         "exec-1",
		BuyerAccountNumber:  "buyer",
		SellerAccountNumber: "seller",
		TradeAmount:         200,
		TradeCurrencyCode:   model.RSD,
	}
	resRepo := &fakeOtcReservationRepo{reservation: existing}
	svc := newOtcFundsService(resRepo, nil)

	result, err := svc.Reserve(context.Background(), "exec-1", "buyer", "seller", 100, model.RSD)

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "execution id already exists with different reservation parameters")
}

func TestReserve_BuyerAccountNotFound(t *testing.T) {
	t.Parallel()

	resRepo := &fakeOtcReservationRepo{reservation: nil}
	accRepo := &fakeOtcAccountRepo{accounts: map[string]*model.Account{}}
	svc := newOtcFundsService(resRepo, accRepo)

	result, err := svc.Reserve(context.Background(), "exec-1", "buyer", "seller", 100, model.RSD)

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "buyer account not found")
}

func TestReserve_SellerAccountNotFound(t *testing.T) {
	t.Parallel()

	resRepo := &fakeOtcReservationRepo{reservation: nil}
	accRepo := &fakeOtcAccountRepo{
		accounts: map[string]*model.Account{
			"buyer": {
				AccountNumber:    "buyer",
				AvailableBalance: 10000,
				Balance:          10000,
				Currency:         model.Currency{Code: model.RSD},
			},
		},
	}
	svc := newOtcFundsService(resRepo, accRepo)

	result, err := svc.Reserve(context.Background(), "exec-1", "buyer", "seller", 100, model.RSD)

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "seller account not found")
}

func TestReserve_InsufficientBuyerFunds(t *testing.T) {
	t.Parallel()

	resRepo := &fakeOtcReservationRepo{reservation: nil}
	accRepo := &fakeOtcAccountRepo{
		accounts: map[string]*model.Account{
			"buyer": {
				AccountNumber:    "buyer",
				AvailableBalance: 50,
				Balance:          50,
				Currency:         model.Currency{Code: model.RSD},
			},
			"seller": {
				AccountNumber:    "seller",
				AvailableBalance: 10000,
				Balance:          10000,
				Currency:         model.Currency{Code: model.RSD},
			},
		},
	}
	svc := newOtcFundsService(resRepo, accRepo)

	result, err := svc.Reserve(context.Background(), "exec-1", "buyer", "seller", 100, model.RSD)

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "insufficient buyer funds")
}

func TestReserve_UpdateBalanceError(t *testing.T) {
	t.Parallel()

	resRepo := &fakeOtcReservationRepo{reservation: nil}
	accRepo := &fakeOtcAccountRepo{
		accounts: map[string]*model.Account{
			"buyer": {
				AccountNumber:    "buyer",
				AvailableBalance: 10000,
				Balance:          10000,
				Currency:         model.Currency{Code: model.RSD},
			},
			"seller": {
				AccountNumber:    "seller",
				AvailableBalance: 10000,
				Balance:          10000,
				Currency:         model.Currency{Code: model.RSD},
			},
		},
		updateBalErr: fmt.Errorf("db write error"),
	}
	svc := newOtcFundsService(resRepo, accRepo)

	result, err := svc.Reserve(context.Background(), "exec-1", "buyer", "seller", 100, model.RSD)

	require.Error(t, err)
	require.Nil(t, result)
}

func TestReserve_CreateReservationError(t *testing.T) {
	t.Parallel()

	resRepo := &fakeOtcReservationRepo{reservation: nil, createErr: fmt.Errorf("create failed")}
	accRepo := &fakeOtcAccountRepo{
		accounts: map[string]*model.Account{
			"buyer": {
				AccountNumber:    "buyer",
				AvailableBalance: 10000,
				Balance:          10000,
				Currency:         model.Currency{Code: model.RSD},
			},
			"seller": {
				AccountNumber:    "seller",
				AvailableBalance: 10000,
				Balance:          10000,
				Currency:         model.Currency{Code: model.RSD},
			},
		},
	}
	svc := newOtcFundsService(resRepo, accRepo)

	result, err := svc.Reserve(context.Background(), "exec-1", "buyer", "seller", 100, model.RSD)

	require.Error(t, err)
	require.Nil(t, result)
}

func TestReserve_Success(t *testing.T) {
	t.Parallel()

	resRepo := &fakeOtcReservationRepo{reservation: nil}
	accRepo := &fakeOtcAccountRepo{
		accounts: map[string]*model.Account{
			"buyer": {
				AccountNumber:    "buyer",
				AvailableBalance: 10000,
				Balance:          10000,
				Currency:         model.Currency{Code: model.RSD},
			},
			"seller": {
				AccountNumber:    "seller",
				AvailableBalance: 10000,
				Balance:          10000,
				Currency:         model.Currency{Code: model.RSD},
			},
		},
	}
	svc := newOtcFundsService(resRepo, accRepo)

	result, err := svc.Reserve(context.Background(), "exec-1", "buyer", "seller", 100, model.RSD)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, model.OtcFundsReservationStatusReserved, result.Status)
	require.Equal(t, 100.0, result.TradeAmount)
}

func TestReserve_ConvertTradeAmount_SourceError(t *testing.T) {
	t.Parallel()

	resRepo := &fakeOtcReservationRepo{reservation: nil}
	accRepo := &fakeOtcAccountRepo{
		accounts: map[string]*model.Account{
			"buyer": {
				AccountNumber:    "buyer",
				AvailableBalance: 10000,
				Balance:          10000,
				Currency:         model.Currency{Code: model.EUR},
			},
			"seller": {
				AccountNumber:    "seller",
				AvailableBalance: 10000,
				Balance:          10000,
				Currency:         model.Currency{Code: model.RSD},
			},
		},
	}
	// ExchangeService with empty rates will fail conversion
	exchangeRateRepo := &fakeExchangeRateRepo{rates: []model.ExchangeRate{}}
	svc := newOtcFundsServiceWithExchange(resRepo, accRepo, exchangeRateRepo)

	result, err := svc.Reserve(context.Background(), "exec-1", "buyer", "seller", 100, model.RSD)

	require.Error(t, err)
	require.Nil(t, result)
}

func TestCommit_EmptyExecutionID(t *testing.T) {
	t.Parallel()

	svc := newOtcFundsService(nil, nil)
	result, err := svc.Commit(context.Background(), "")

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "execution id is required")
}

func TestCommit_ReservationNotFound(t *testing.T) {
	t.Parallel()

	repo := &fakeOtcReservationRepo{reservation: nil}
	svc := newOtcFundsService(repo, nil)

	result, err := svc.Commit(context.Background(), "exec-1")

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "OTC funds reservation not found")
}

func TestCommit_FindByExecutionIDError(t *testing.T) {
	t.Parallel()

	repo := &fakeOtcReservationRepo{findErr: fmt.Errorf("db error")}
	svc := newOtcFundsService(repo, nil)

	result, err := svc.Commit(context.Background(), "exec-1")

	require.Error(t, err)
	require.Nil(t, result)
}

func TestCommit_AlreadyCommitted_Idempotent(t *testing.T) {
	t.Parallel()

	reservation := &model.OtcFundsReservation{
		ExecutionID: "exec-1",
		Status:      model.OtcFundsReservationStatusCommitted,
	}
	repo := &fakeOtcReservationRepo{reservation: reservation}
	svc := newOtcFundsService(repo, nil)

	result, err := svc.Commit(context.Background(), "exec-1")

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, model.OtcFundsReservationStatusCommitted, result.Status)
}

func TestCommit_ReleasedStatus_Error(t *testing.T) {
	t.Parallel()

	reservation := &model.OtcFundsReservation{
		ExecutionID: "exec-1",
		Status:      model.OtcFundsReservationStatusReleased,
	}
	repo := &fakeOtcReservationRepo{reservation: reservation}
	svc := newOtcFundsService(repo, nil)

	result, err := svc.Commit(context.Background(), "exec-1")

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "cannot commit released OTC funds")
}

func TestCommit_RefundedStatus_Error(t *testing.T) {
	t.Parallel()

	reservation := &model.OtcFundsReservation{
		ExecutionID: "exec-1",
		Status:      model.OtcFundsReservationStatusRefunded,
	}
	repo := &fakeOtcReservationRepo{reservation: reservation}
	svc := newOtcFundsService(repo, nil)

	result, err := svc.Commit(context.Background(), "exec-1")

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "cannot commit refunded OTC funds")
}

func TestCommit_UnknownStatus_Error(t *testing.T) {
	t.Parallel()

	reservation := &model.OtcFundsReservation{
		ExecutionID: "exec-1",
		Status:      model.OtcFundsReservationStatus("UNKNOWN"),
	}
	repo := &fakeOtcReservationRepo{reservation: reservation}
	svc := newOtcFundsService(repo, nil)

	result, err := svc.Commit(context.Background(), "exec-1")

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "cannot commit OTC funds in current status")
}

func TestCommit_BuyerAccountNotFound(t *testing.T) {
	t.Parallel()

	reservation := &model.OtcFundsReservation{
		ExecutionID:         "exec-1",
		BuyerAccountNumber:  "missing-buyer",
		SellerAccountNumber: "seller",
		Status:              model.OtcFundsReservationStatusReserved,
		SourceAmount:        500,
	}
	resRepo := &fakeOtcReservationRepo{reservation: reservation}
	accRepo := &fakeOtcAccountRepo{accounts: map[string]*model.Account{}}
	svc := newOtcFundsService(resRepo, accRepo)

	result, err := svc.Commit(context.Background(), "exec-1")

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "buyer account not found")
}

func TestCommit_SellerAccountNotFound(t *testing.T) {
	t.Parallel()

	reservation := &model.OtcFundsReservation{
		ExecutionID:         "exec-1",
		BuyerAccountNumber:  "buyer",
		SellerAccountNumber: "missing-seller",
		Status:              model.OtcFundsReservationStatusReserved,
		SourceAmount:        500,
	}
	resRepo := &fakeOtcReservationRepo{reservation: reservation}
	accRepo := &fakeOtcAccountRepo{
		accounts: map[string]*model.Account{
			"buyer": {
				AccountNumber:    "buyer",
				Balance:          10000,
				AvailableBalance: 5000,
			},
		},
	}
	svc := newOtcFundsService(resRepo, accRepo)

	result, err := svc.Commit(context.Background(), "exec-1")

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "seller account not found")
}

func TestCommit_ReservedFundsInconsistent(t *testing.T) {
	t.Parallel()

	reservation := &model.OtcFundsReservation{
		ExecutionID:         "exec-1",
		BuyerAccountNumber:  "buyer",
		SellerAccountNumber: "seller",
		Status:              model.OtcFundsReservationStatusReserved,
		SourceAmount:        5000,
	}
	resRepo := &fakeOtcReservationRepo{reservation: reservation}
	accRepo := &fakeOtcAccountRepo{
		accounts: map[string]*model.Account{
			"buyer": {
				AccountNumber:    "buyer",
				Balance:          10000,
				AvailableBalance: 9000, // reserved = 10000 - 9000 = 1000 < 5000
			},
			"seller": {
				AccountNumber:    "seller",
				Balance:          10000,
				AvailableBalance: 10000,
			},
		},
	}
	svc := newOtcFundsService(resRepo, accRepo)

	result, err := svc.Commit(context.Background(), "exec-1")

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "reserved buyer funds are inconsistent")
}

func TestCommit_BuyerBalanceBelowReserved(t *testing.T) {
	t.Parallel()

	reservation := &model.OtcFundsReservation{
		ExecutionID:         "exec-1",
		BuyerAccountNumber:  "buyer",
		SellerAccountNumber: "seller",
		Status:              model.OtcFundsReservationStatusReserved,
		SourceAmount:        5000,
	}
	resRepo := &fakeOtcReservationRepo{reservation: reservation}
	accRepo := &fakeOtcAccountRepo{
		accounts: map[string]*model.Account{
			"buyer": {
				AccountNumber:    "buyer",
				Balance:          3000,  // balance < sourceAmount
				AvailableBalance: -2000, // reserved = 3000 - (-2000) = 5000 >= 5000 passes first check
			},
			"seller": {
				AccountNumber:    "seller",
				Balance:          10000,
				AvailableBalance: 10000,
			},
		},
	}
	svc := newOtcFundsService(resRepo, accRepo)

	result, err := svc.Commit(context.Background(), "exec-1")

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "buyer balance is below reserved OTC funds")
}

func TestCommit_UpdateBuyerBalanceError(t *testing.T) {
	t.Parallel()

	reservation := &model.OtcFundsReservation{
		ExecutionID:             "exec-1",
		BuyerAccountNumber:      "buyer",
		SellerAccountNumber:     "seller",
		Status:                  model.OtcFundsReservationStatusReserved,
		SourceAmount:            500,
		DestinationAmount:       500,
		SourceCurrencyCode:      model.RSD,
		DestinationCurrencyCode: model.RSD,
	}
	resRepo := &fakeOtcReservationRepo{reservation: reservation}
	accRepo := &fakeOtcAccountRepo{
		accounts: map[string]*model.Account{
			"buyer": {
				AccountNumber:    "buyer",
				Balance:          10000,
				AvailableBalance: 5000, // reserved = 5000 >= 500
			},
			"seller": {
				AccountNumber:    "seller",
				Balance:          10000,
				AvailableBalance: 10000,
			},
		},
		updateBalErr: fmt.Errorf("db write error"),
	}
	svc := newOtcFundsService(resRepo, accRepo)

	result, err := svc.Commit(context.Background(), "exec-1")

	require.Error(t, err)
	require.Nil(t, result)
}

func TestCommit_Success(t *testing.T) {
	t.Parallel()

	reservation := &model.OtcFundsReservation{
		ExecutionID:             "exec-1",
		BuyerAccountNumber:      "buyer",
		SellerAccountNumber:     "seller",
		Status:                  model.OtcFundsReservationStatusReserved,
		SourceAmount:            500,
		DestinationAmount:       500,
		SourceCurrencyCode:      model.RSD,
		DestinationCurrencyCode: model.RSD,
	}
	resRepo := &fakeOtcReservationRepo{reservation: reservation}
	accRepo := &fakeOtcAccountRepo{
		accounts: map[string]*model.Account{
			"buyer": {
				AccountNumber:    "buyer",
				Balance:          10000,
				AvailableBalance: 5000,
			},
			"seller": {
				AccountNumber:    "seller",
				Balance:          10000,
				AvailableBalance: 10000,
			},
		},
	}
	svc := newOtcFundsService(resRepo, accRepo)

	result, err := svc.Commit(context.Background(), "exec-1")

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, model.OtcFundsReservationStatusCommitted, result.Status)
}

func TestCommit_SaveError(t *testing.T) {
	t.Parallel()

	reservation := &model.OtcFundsReservation{
		ExecutionID:         "exec-1",
		BuyerAccountNumber:  "buyer",
		SellerAccountNumber: "seller",
		Status:              model.OtcFundsReservationStatusReserved,
		SourceAmount:        500,
		DestinationAmount:   500,
	}
	resRepo := &fakeOtcReservationRepo{
		reservation: reservation,
		saveErr:     fmt.Errorf("save failed"),
	}
	accRepo := &fakeOtcAccountRepo{
		accounts: map[string]*model.Account{
			"buyer": {
				AccountNumber:    "buyer",
				Balance:          10000,
				AvailableBalance: 5000,
			},
			"seller": {
				AccountNumber:    "seller",
				Balance:          10000,
				AvailableBalance: 10000,
			},
		},
	}
	svc := newOtcFundsService(resRepo, accRepo)

	result, err := svc.Commit(context.Background(), "exec-1")

	require.Error(t, err)
	require.Nil(t, result)
}

func TestRefund_EmptyExecutionID(t *testing.T) {
	t.Parallel()

	svc := newOtcFundsService(nil, nil)
	result, err := svc.Refund(context.Background(), "")

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "execution id is required")
}

func TestRefund_ReservationNotFound(t *testing.T) {
	t.Parallel()

	repo := &fakeOtcReservationRepo{reservation: nil}
	svc := newOtcFundsService(repo, nil)

	result, err := svc.Refund(context.Background(), "exec-1")

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "OTC funds reservation not found")
}

func TestRefund_FindByExecutionIDError(t *testing.T) {
	t.Parallel()

	repo := &fakeOtcReservationRepo{findErr: fmt.Errorf("db error")}
	svc := newOtcFundsService(repo, nil)

	result, err := svc.Refund(context.Background(), "exec-1")

	require.Error(t, err)
	require.Nil(t, result)
}

func TestRefund_AlreadyRefunded_Idempotent(t *testing.T) {
	t.Parallel()

	reservation := &model.OtcFundsReservation{
		ExecutionID: "exec-1",
		Status:      model.OtcFundsReservationStatusRefunded,
	}
	repo := &fakeOtcReservationRepo{reservation: reservation}
	svc := newOtcFundsService(repo, nil)

	result, err := svc.Refund(context.Background(), "exec-1")

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, model.OtcFundsReservationStatusRefunded, result.Status)
}

func TestRefund_ReleasedStatus_Error(t *testing.T) {
	t.Parallel()

	reservation := &model.OtcFundsReservation{
		ExecutionID: "exec-1",
		Status:      model.OtcFundsReservationStatusReleased,
	}
	repo := &fakeOtcReservationRepo{reservation: reservation}
	svc := newOtcFundsService(repo, nil)

	result, err := svc.Refund(context.Background(), "exec-1")

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "cannot refund released OTC funds")
}

func TestRefund_ReservedStatus_Error(t *testing.T) {
	t.Parallel()

	reservation := &model.OtcFundsReservation{
		ExecutionID: "exec-1",
		Status:      model.OtcFundsReservationStatusReserved,
	}
	repo := &fakeOtcReservationRepo{reservation: reservation}
	svc := newOtcFundsService(repo, nil)

	result, err := svc.Refund(context.Background(), "exec-1")

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "cannot refund uncommitted OTC funds")
}

func TestRefund_UnknownStatus_Error(t *testing.T) {
	t.Parallel()

	reservation := &model.OtcFundsReservation{
		ExecutionID: "exec-1",
		Status:      model.OtcFundsReservationStatus("UNKNOWN"),
	}
	repo := &fakeOtcReservationRepo{reservation: reservation}
	svc := newOtcFundsService(repo, nil)

	result, err := svc.Refund(context.Background(), "exec-1")

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "cannot refund OTC funds in current status")
}

func TestRefund_BuyerAccountNotFound(t *testing.T) {
	t.Parallel()

	reservation := &model.OtcFundsReservation{
		ExecutionID:         "exec-1",
		BuyerAccountNumber:  "missing-buyer",
		SellerAccountNumber: "seller",
		Status:              model.OtcFundsReservationStatusCommitted,
		SourceAmount:        500,
		DestinationAmount:   500,
	}
	resRepo := &fakeOtcReservationRepo{reservation: reservation}
	accRepo := &fakeOtcAccountRepo{accounts: map[string]*model.Account{}}
	svc := newOtcFundsService(resRepo, accRepo)

	result, err := svc.Refund(context.Background(), "exec-1")

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "buyer account not found")
}

func TestRefund_SellerAccountNotFound(t *testing.T) {
	t.Parallel()

	reservation := &model.OtcFundsReservation{
		ExecutionID:         "exec-1",
		BuyerAccountNumber:  "buyer",
		SellerAccountNumber: "missing-seller",
		Status:              model.OtcFundsReservationStatusCommitted,
		SourceAmount:        500,
		DestinationAmount:   500,
	}
	resRepo := &fakeOtcReservationRepo{reservation: reservation}
	accRepo := &fakeOtcAccountRepo{
		accounts: map[string]*model.Account{
			"buyer": {
				AccountNumber:    "buyer",
				Balance:          10000,
				AvailableBalance: 10000,
			},
		},
	}
	svc := newOtcFundsService(resRepo, accRepo)

	result, err := svc.Refund(context.Background(), "exec-1")

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "seller account not found")
}

func TestRefund_InsufficientSellerFunds_Balance(t *testing.T) {
	t.Parallel()

	reservation := &model.OtcFundsReservation{
		ExecutionID:         "exec-1",
		BuyerAccountNumber:  "buyer",
		SellerAccountNumber: "seller",
		Status:              model.OtcFundsReservationStatusCommitted,
		SourceAmount:        500,
		DestinationAmount:   500,
	}
	resRepo := &fakeOtcReservationRepo{reservation: reservation}
	accRepo := &fakeOtcAccountRepo{
		accounts: map[string]*model.Account{
			"buyer": {
				AccountNumber:    "buyer",
				Balance:          10000,
				AvailableBalance: 10000,
			},
			"seller": {
				AccountNumber:    "seller",
				Balance:          100, // insufficient
				AvailableBalance: 100,
			},
		},
	}
	svc := newOtcFundsService(resRepo, accRepo)

	result, err := svc.Refund(context.Background(), "exec-1")

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "insufficient seller funds for refund")
}

func TestRefund_InsufficientSellerFunds_AvailableBalance(t *testing.T) {
	t.Parallel()

	reservation := &model.OtcFundsReservation{
		ExecutionID:         "exec-1",
		BuyerAccountNumber:  "buyer",
		SellerAccountNumber: "seller",
		Status:              model.OtcFundsReservationStatusCommitted,
		SourceAmount:        500,
		DestinationAmount:   500,
	}
	resRepo := &fakeOtcReservationRepo{reservation: reservation}
	accRepo := &fakeOtcAccountRepo{
		accounts: map[string]*model.Account{
			"buyer": {
				AccountNumber:    "buyer",
				Balance:          10000,
				AvailableBalance: 10000,
			},
			"seller": {
				AccountNumber:    "seller",
				Balance:          10000,
				AvailableBalance: 100, // insufficient available
			},
		},
	}
	svc := newOtcFundsService(resRepo, accRepo)

	result, err := svc.Refund(context.Background(), "exec-1")

	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "insufficient seller funds for refund")
}

func TestRefund_UpdateSellerBalanceError(t *testing.T) {
	t.Parallel()

	reservation := &model.OtcFundsReservation{
		ExecutionID:         "exec-1",
		BuyerAccountNumber:  "buyer",
		SellerAccountNumber: "seller",
		Status:              model.OtcFundsReservationStatusCommitted,
		SourceAmount:        500,
		DestinationAmount:   500,
	}
	resRepo := &fakeOtcReservationRepo{reservation: reservation}
	accRepo := &fakeOtcAccountRepo{
		accounts: map[string]*model.Account{
			"buyer": {
				AccountNumber:    "buyer",
				Balance:          10000,
				AvailableBalance: 10000,
			},
			"seller": {
				AccountNumber:    "seller",
				Balance:          10000,
				AvailableBalance: 10000,
			},
		},
		updateBalErr: fmt.Errorf("db write error"),
	}
	svc := newOtcFundsService(resRepo, accRepo)

	result, err := svc.Refund(context.Background(), "exec-1")

	require.Error(t, err)
	require.Nil(t, result)
}

func TestRefund_Success(t *testing.T) {
	t.Parallel()

	reservation := &model.OtcFundsReservation{
		ExecutionID:         "exec-1",
		BuyerAccountNumber:  "buyer",
		SellerAccountNumber: "seller",
		Status:              model.OtcFundsReservationStatusCommitted,
		SourceAmount:        500,
		DestinationAmount:   500,
	}
	resRepo := &fakeOtcReservationRepo{reservation: reservation}
	accRepo := &fakeOtcAccountRepo{
		accounts: map[string]*model.Account{
			"buyer": {
				AccountNumber:    "buyer",
				Balance:          10000,
				AvailableBalance: 10000,
			},
			"seller": {
				AccountNumber:    "seller",
				Balance:          10000,
				AvailableBalance: 10000,
			},
		},
	}
	svc := newOtcFundsService(resRepo, accRepo)

	result, err := svc.Refund(context.Background(), "exec-1")

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, model.OtcFundsReservationStatusRefunded, result.Status)
}

func TestRefund_SaveError(t *testing.T) {
	t.Parallel()

	reservation := &model.OtcFundsReservation{
		ExecutionID:         "exec-1",
		BuyerAccountNumber:  "buyer",
		SellerAccountNumber: "seller",
		Status:              model.OtcFundsReservationStatusCommitted,
		SourceAmount:        500,
		DestinationAmount:   500,
	}
	resRepo := &fakeOtcReservationRepo{
		reservation: reservation,
		saveErr:     fmt.Errorf("save failed"),
	}
	accRepo := &fakeOtcAccountRepo{
		accounts: map[string]*model.Account{
			"buyer": {
				AccountNumber:    "buyer",
				Balance:          10000,
				AvailableBalance: 10000,
			},
			"seller": {
				AccountNumber:    "seller",
				Balance:          10000,
				AvailableBalance: 10000,
			},
		},
	}
	svc := newOtcFundsService(resRepo, accRepo)

	result, err := svc.Refund(context.Background(), "exec-1")

	require.Error(t, err)
	require.Nil(t, result)
}
