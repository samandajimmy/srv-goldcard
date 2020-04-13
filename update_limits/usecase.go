package update_limits

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// UseCase represent the registrations usecases
type UseCase interface {
	DecreasedSTL(c echo.Context, pl models.PayloadCoreDecreasedSTL) models.ResponseErrors
}
