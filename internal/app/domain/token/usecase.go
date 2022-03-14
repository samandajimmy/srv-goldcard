package token

import (
	"srv-goldcard/internal/app/model"

	"github.com/labstack/echo"
)

// UseCase represent the token's usecases
type UseCase interface {
	CreateToken(c echo.Context, accToken *model.AccountToken) error
	GetToken(c echo.Context, username string, password string) (*model.AccountToken, error)
	RefreshToken(c echo.Context, username string, password string) (*model.AccountToken, error)
	RefreshAllToken() error
}
