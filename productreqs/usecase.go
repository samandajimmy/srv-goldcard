package productreqs

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// UseCase represent the product requirements usecases
type UseCase interface {
	ProductRequirements(echo.Context) (models.Requirements, error)
	InsertPublicHolidayDate(echo.Context, models.PayloadInsertPublicHoliday) (models.PublicHolidayDate, error)
	GetPublicHolidayDate(c echo.Context) (models.PublicHolidayDate, error)
}
