package activation

import (
	"srv-goldcard/internal/app/model"

	"github.com/labstack/echo"
)

// Repository represent the activation repository contract
type Repository interface {
	PostActivations(echo.Context, model.Account) error
	UpdateGoldLimit(echo.Context, model.Card) error
	GetStoredGoldPrice(c echo.Context) (int64, error)
}

// RestRepository represent the rest activation repository contract
type RestRepository interface {
	GetDetailGoldUser(c echo.Context, accNumber string) (map[string]interface{}, error)
	ActivationsToCore(c echo.Context, acc *model.Account) error
	ActivationsToBRI(c echo.Context, acc model.Account, pa model.PayloadActivations) error
}
