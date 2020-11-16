package process_handler

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// UseCase represent the process handler usecases
type UseCase interface {
	PostProcessHandler(c echo.Context, ps models.ProcessStatus) error
	UpdateErrorStatus(c echo.Context, acc models.Account) error
}
