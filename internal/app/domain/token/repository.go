package token

import (
	"srv-goldcard/internal/app/model"

	"github.com/labstack/echo"
)

// Repository represent the Account Token's repository contract
type Repository interface {
	Create(c echo.Context, accToken *model.AccountToken) error
	GetByUsername(c echo.Context, accToken *model.AccountToken) error
	UpdateToken(c echo.Context, accToken *model.AccountToken) error
	UpdateAllAccountTokenExpiry() error
}
