package process_handler

import (
	"srv-goldcard/internal/app/model"

	"github.com/labstack/echo"
)

// UseCase represent the process handler usecases
type UseCase interface {
	IsProcessedAppExisted(c echo.Context, acc model.Account) (bool, error)
	UpsertAppProcess(c echo.Context, acc *model.Account, errStr string) error
	UpdateErrorStatus(c echo.Context, acc model.Account) error
}
