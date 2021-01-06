package productreqs

import (
	"github.com/labstack/echo"
)

// Repository represent the product requirements's repository contract
type Repository interface {
	GetSertPublicHolidayDate(echo.Context, []string) (string, error)
}
