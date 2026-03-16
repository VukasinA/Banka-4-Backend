package validator

import (
	"banking-service/internal/model"

	"github.com/go-playground/validator/v10"
)

func validateForeignCurrency(fl validator.FieldLevel) bool {
	return model.AllowedForeignCurrencies[model.CurrencyCode(fl.Field().String())]
}
