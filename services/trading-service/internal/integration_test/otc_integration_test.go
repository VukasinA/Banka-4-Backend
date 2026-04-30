//go:build integration

package integration_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/RAF-SI-2025/Banka-4-Backend/services/trading-service/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// ---------------------------------------------------------------------------
// Seed helpers
// ---------------------------------------------------------------------------

func seedAssetAndStock(t *testing.T, db *gorm.DB, ticker string) (*model.Asset, *model.Stock) {
	t.Helper()
	asset := &model.Asset{Ticker: ticker, Name: "Company " + ticker, AssetType: model.AssetTypeStock}
	require.NoError(t, db.Create(asset).Error)
	stock := &model.Stock{AssetID: asset.AssetID}
	require.NoError(t, db.Create(stock).Error)
	return asset, stock
}

func seedOwnership(t *testing.T, db *gorm.DB, userID uint, assetID uint, amount, public float64) {
	t.Helper()
	require.NoError(t, db.Create(&model.AssetOwnership{
		UserId:       userID,
		OwnerType:    model.OwnerTypeClient,
		AssetID:      assetID,
		Amount:       amount,
		PublicAmount: public,
	}).Error)
}

func clientAuthHeader(t *testing.T, identityID, clientID uint) string {
	return authHeaderForClient(t, identityID, clientID)
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestOTC_CreateOffer_Success(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	router, _ := setupTestRouter(t, db)

	asset, stock := seedAssetAndStock(t, db, uniqueValue(t, "CRTO"))
	seedOwnership(t, db, 20, asset.AssetID, 100, 100)

	body := map[string]any{
		"seller_id": 20, "stock_id": stock.StockID, "amount": 10,
		"price_per_stock": 50.0, "premium": 5.0,
		"settlement_date":      time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339),
		"buyer_account_number": "buyer-account-001",
	}
	rec := performRequest(t, router, http.MethodPost, "/api/otc/offers", body, clientAuthHeader(t, 10, 10))
	requireStatus(t, rec, http.StatusCreated)

	resp := decodeResponse[map[string]any](t, rec)
	assert.Equal(t, float64(10), resp["buyer_id"])
	assert.Equal(t, float64(20), resp["seller_id"])
	assert.Equal(t, float64(10), resp["amount"])
	assert.Equal(t, "ACTIVE", resp["status"])
}

func TestOTC_CreateOffer_SelfOffer_BadRequest(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	router, _ := setupTestRouter(t, db)

	asset, stock := seedAssetAndStock(t, db, uniqueValue(t, "SELF"))
	seedOwnership(t, db, 10, asset.AssetID, 100, 100)

	body := map[string]any{
		"seller_id": 10, "stock_id": stock.StockID, "amount": 5,
		"price_per_stock": 50.0, "premium": 5.0,
		"settlement_date":      time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339),
		"buyer_account_number": "acc",
	}
	rec := performRequest(t, router, http.MethodPost, "/api/otc/offers", body, clientAuthHeader(t, 10, 10))
	requireStatus(t, rec, http.StatusBadRequest)
}

func TestOTC_CreateOffer_Unauthorized(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	router, _ := setupTestRouter(t, db)

	body := map[string]any{
		"seller_id": 20, "stock_id": 1, "amount": 5,
		"price_per_stock": 50.0, "premium": 5.0,
		"settlement_date":      time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		"buyer_account_number": "acc",
	}
	rec := performRequest(t, router, http.MethodPost, "/api/otc/offers", body, "")
	requireStatus(t, rec, http.StatusUnauthorized)
}

func TestOTC_SendCounterOffer_Success(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	router, _ := setupTestRouter(t, db)

	asset, stock := seedAssetAndStock(t, db, uniqueValue(t, "CNTR"))
	seedOwnership(t, db, 20, asset.AssetID, 100, 100)

	offerBody := map[string]any{
		"seller_id": 20, "stock_id": stock.StockID, "amount": 10,
		"price_per_stock": 50.0, "premium": 5.0,
		"settlement_date":      time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339),
		"buyer_account_number": "buyer-acc",
	}
	rec := performRequest(t, router, http.MethodPost, "/api/otc/offers", offerBody, clientAuthHeader(t, 10, 10))
	requireStatus(t, rec, http.StatusCreated)
	offerID := uint(decodeResponse[map[string]any](t, rec)["otc_offer_id"].(float64))

	counterBody := map[string]any{
		"amount": 8, "price_per_stock": 55.0, "premium": 6.0,
		"settlement_date": time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339),
		"account_number":  "seller-acc",
	}
	rec = performRequest(t, router, http.MethodPut, fmt.Sprintf("/api/otc/offers/%d/counter", offerID), counterBody, clientAuthHeader(t, 20, 20))
	requireStatus(t, rec, http.StatusOK)

	updated := decodeResponse[map[string]any](t, rec)
	assert.Equal(t, float64(8), updated["amount"])
	assert.Equal(t, float64(20), updated["modified_by"])
}

func TestOTC_SendCounterOffer_SameUserTwice_BadRequest(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	router, _ := setupTestRouter(t, db)

	asset, stock := seedAssetAndStock(t, db, uniqueValue(t, "DBL"))
	seedOwnership(t, db, 20, asset.AssetID, 100, 100)

	offerBody := map[string]any{
		"seller_id": 20, "stock_id": stock.StockID, "amount": 10,
		"price_per_stock": 50.0, "premium": 5.0,
		"settlement_date":      time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339),
		"buyer_account_number": "buyer-acc",
	}
	rec := performRequest(t, router, http.MethodPost, "/api/otc/offers", offerBody, clientAuthHeader(t, 10, 10))
	requireStatus(t, rec, http.StatusCreated)
	offerID := uint(decodeResponse[map[string]any](t, rec)["otc_offer_id"].(float64))

	counterBody := map[string]any{
		"amount": 9, "price_per_stock": 50.0, "premium": 5.0,
		"settlement_date": time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339),
	}
	rec = performRequest(t, router, http.MethodPut, fmt.Sprintf("/api/otc/offers/%d/counter", offerID), counterBody, clientAuthHeader(t, 10, 10))
	requireStatus(t, rec, http.StatusBadRequest)
}

func TestOTC_AcceptOffer_Success_CreatesContractAndIncreasesReservedAmount(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	router, _ := setupTestRouter(t, db)

	asset, stock := seedAssetAndStock(t, db, uniqueValue(t, "ACPT"))
	seedOwnership(t, db, 20, asset.AssetID, 100, 100)

	offerBody := map[string]any{
		"seller_id": 20, "stock_id": stock.StockID, "amount": 10,
		"price_per_stock": 50.0, "premium": 5.0,
		"settlement_date":      time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339),
		"buyer_account_number": "buyer-acc",
	}
	rec := performRequest(t, router, http.MethodPost, "/api/otc/offers", offerBody, clientAuthHeader(t, 10, 10))
	requireStatus(t, rec, http.StatusCreated)
	offerID := uint(decodeResponse[map[string]any](t, rec)["otc_offer_id"].(float64))

	rec = performRequest(t, router, http.MethodPatch,
		fmt.Sprintf("/api/otc/offers/%d/accept", offerID),
		map[string]any{"account_number": "seller-acc"},
		clientAuthHeader(t, 20, 20))
	requireStatus(t, rec, http.StatusCreated)

	contract := decodeResponse[map[string]any](t, rec)
	require.NotNil(t, contract["otc_option_contract_id"])
	assert.Equal(t, float64(10), contract["amount"])
	assert.Equal(t, float64(50.0), contract["strike_price"])

	// Core assertion: reserved_amount must increase on the seller's row.
	var ownership model.AssetOwnership
	err := db.Where("user_id = ? AND owner_type = ? AND asset_id = ?",
		20, model.OwnerTypeClient, asset.AssetID).First(&ownership).Error
	require.NoError(t, err)
	assert.Equal(t, float64(10), ownership.ReservedAmount)
}

func TestOTC_AcceptOffer_BuyerCannotAcceptOwnOffer(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	router, _ := setupTestRouter(t, db)

	asset, stock := seedAssetAndStock(t, db, uniqueValue(t, "SELF2"))
	seedOwnership(t, db, 20, asset.AssetID, 100, 100)

	offerBody := map[string]any{
		"seller_id": 20, "stock_id": stock.StockID, "amount": 10,
		"price_per_stock": 50.0, "premium": 5.0,
		"settlement_date":      time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339),
		"buyer_account_number": "buyer-acc",
	}
	rec := performRequest(t, router, http.MethodPost, "/api/otc/offers", offerBody, clientAuthHeader(t, 10, 10))
	requireStatus(t, rec, http.StatusCreated)
	offerID := uint(decodeResponse[map[string]any](t, rec)["otc_offer_id"].(float64))

	rec = performRequest(t, router, http.MethodPatch, fmt.Sprintf("/api/otc/offers/%d/accept", offerID), nil, clientAuthHeader(t, 10, 10))
	requireStatus(t, rec, http.StatusBadRequest)
}

func TestOTC_RejectOffer_Success(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	router, _ := setupTestRouter(t, db)

	asset, stock := seedAssetAndStock(t, db, uniqueValue(t, "REJ"))
	seedOwnership(t, db, 20, asset.AssetID, 100, 100)

	offerBody := map[string]any{
		"seller_id": 20, "stock_id": stock.StockID, "amount": 10,
		"price_per_stock": 50.0, "premium": 5.0,
		"settlement_date":      time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339),
		"buyer_account_number": "buyer-acc",
	}
	rec := performRequest(t, router, http.MethodPost, "/api/otc/offers", offerBody, clientAuthHeader(t, 10, 10))
	requireStatus(t, rec, http.StatusCreated)
	offerID := uint(decodeResponse[map[string]any](t, rec)["otc_offer_id"].(float64))

	rec = performRequest(t, router, http.MethodPatch, fmt.Sprintf("/api/otc/offers/%d/reject", offerID), nil, clientAuthHeader(t, 20, 20))
	requireStatus(t, rec, http.StatusOK)
	assert.Equal(t, "REJECTED", decodeResponse[map[string]any](t, rec)["status"])
}

func TestOTC_GetMyActiveOffers_ReturnsOnlyOwnOffers(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	router, _ := setupTestRouter(t, db)

	asset, stock := seedAssetAndStock(t, db, uniqueValue(t, "ACTIVE"))
	seedOwnership(t, db, 20, asset.AssetID, 100, 100)

	for i := 0; i < 2; i++ {
		body := map[string]any{
			"seller_id": 20, "stock_id": stock.StockID, "amount": 5,
			"price_per_stock": 50.0, "premium": 5.0,
			"settlement_date":      time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339),
			"buyer_account_number": "buyer-acc",
		}
		rec := performRequest(t, router, http.MethodPost, "/api/otc/offers", body, clientAuthHeader(t, 10, 10))
		requireStatus(t, rec, http.StatusCreated)
	}

	rec := performRequest(t, router, http.MethodGet, "/api/otc/offers/active", nil, clientAuthHeader(t, 10, 10))
	requireStatus(t, rec, http.StatusOK)
	assert.Len(t, decodeResponse[[]map[string]any](t, rec), 2)

	// unrelated user sees nothing
	rec = performRequest(t, router, http.MethodGet, "/api/otc/offers/active", nil, clientAuthHeader(t, 99, 99))
	requireStatus(t, rec, http.StatusOK)
	assert.Empty(t, decodeResponse[[]map[string]any](t, rec))
}

func TestOTC_GetMyOptionContracts_AfterAccept(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)
	router, _ := setupTestRouter(t, db)

	asset, stock := seedAssetAndStock(t, db, uniqueValue(t, "CTCT"))
	seedOwnership(t, db, 20, asset.AssetID, 100, 100)

	offerBody := map[string]any{
		"seller_id": 20, "stock_id": stock.StockID, "amount": 10,
		"price_per_stock": 50.0, "premium": 5.0,
		"settlement_date":      time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339),
		"buyer_account_number": "buyer-acc",
	}
	rec := performRequest(t, router, http.MethodPost, "/api/otc/offers", offerBody, clientAuthHeader(t, 10, 10))
	requireStatus(t, rec, http.StatusCreated)
	offerID := uint(decodeResponse[map[string]any](t, rec)["otc_offer_id"].(float64))

	rec = performRequest(t, router, http.MethodPatch,
		fmt.Sprintf("/api/otc/offers/%d/accept", offerID),
		map[string]any{"account_number": "seller-acc"},
		clientAuthHeader(t, 20, 20))
	requireStatus(t, rec, http.StatusCreated)

	// both buyer and seller must see the contract
	for _, id := range []uint{10, 20} {
		rec = performRequest(t, router, http.MethodGet, "/api/otc/contracts", nil, clientAuthHeader(t, id, id))
		requireStatus(t, rec, http.StatusOK)
		assert.Len(t, decodeResponse[[]map[string]any](t, rec), 1)
	}
}
