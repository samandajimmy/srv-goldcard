package process_handler

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// UseCase represent the process handler usecases
type UseCase interface {
	IsProcessedAppExisted(c echo.Context, acc models.Account) (bool, error)
	UpsertAppProcess(c echo.Context, acc *models.Account, errStr string) error
	UpdateErrorStatus(c echo.Context, acc models.Account) error
}
