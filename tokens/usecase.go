package tokens

import (
	"github.com/labstack/echo"

	"gade/srv-goldcard/models"
)

// UseCase represent the token's usecases
type UseCase interface {
	CreateToken(c echo.Context, accToken *models.AccountToken) error
	GetToken(c echo.Context, username string, password string) (*models.AccountToken, error)
	RefreshToken(c echo.Context, username string, password string) (*models.AccountToken, error)
}
