package productreq

import (
	"srv-goldcard/internal/app/model"

	"github.com/labstack/echo"
)

// UseCase represent the product requirements usecases
type UseCase interface {
	ProductRequirements(echo.Context) (model.Requirements, error)
	InsertPublicHolidayDate(echo.Context, model.PayloadInsertPublicHoliday) (model.PublicHolidayDate, error)
	GetPublicHolidayDate(c echo.Context) (model.PublicHolidayDate, error)
}
