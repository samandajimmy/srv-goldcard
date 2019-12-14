package productreqs

import (
	"github.com/labstack/echo"
)

// Repository represent the product requirements repository
type Repository interface {
	ProductRequirements(echo.Context) (map[string]interface{}, error)
}
