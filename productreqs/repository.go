package productreqs

import (
	"github.com/labstack/echo"
)

// Repository represent the product requirements's repository contract
type Repository interface {
	GetSertPublicHolidayDate(c echo.Context, phds []string) (string, error)
}
