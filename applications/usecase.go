package applications

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// UseCase represent the applications usecases
type UseCase interface {
	PostSavingAccount(echo.Context, *models.Applications) error
}
