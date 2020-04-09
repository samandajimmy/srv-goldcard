package process_handler

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// UseCase represent the process handler usecases
type UseCase interface {
	PostProcessHandler(c echo.Context, ps models.ProcessStatus) error
	UpdateCounterError(c echo.Context, acc models.Account)
	UpdateErrorStatus(c echo.Context, acc models.Account) error
}
