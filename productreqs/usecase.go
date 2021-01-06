package productreqs

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// UseCase represent the product requirements usecases
type UseCase interface {
	ProductRequirements(echo.Context) (models.Requirements, error)
	GetSertPublicHolidayDate(echo.Context, models.PayloadGetSertPublicHoliday) (models.PublicHolidayDate, error)
}
