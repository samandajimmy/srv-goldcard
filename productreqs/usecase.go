package productreqs

import (
	"github.com/labstack/echo"
)

// UseCase represent the product requirements usecases
type UseCase interface {
	ProductRequirements(echo.Context) (map[string]interface{}, error)
}
