package activations

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// Repository represent the activations repository contract
type Repository interface {
	PostActivations(echo.Context, models.Account) error
	GetAccountByAppNumber(c echo.Context, acc *models.Account) error
	UpdateGoldLimit(echo.Context, models.Card) error
}

// RestRepository represent the rest activations repository contract
type RestRepository interface {
	GetDetailGoldUser(c echo.Context, accNumber string) (map[string]string, error)
	ActivationsToCore(c echo.Context, acc models.Account) error
	OpenRecalculateToCore(c echo.Context, acc models.Account) error
	ActivationsToBRI(c echo.Context, acc models.Account, pa models.PayloadActivations) error
}
