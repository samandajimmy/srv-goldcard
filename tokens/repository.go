package tokens

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// Repository represent the Account Token's repository contract
type Repository interface {
	Create(c echo.Context, accToken *models.AccountToken) error
	GetByUsername(c echo.Context, accToken *models.AccountToken) error
	UpdateToken(c echo.Context, accToken *models.AccountToken) error
}
