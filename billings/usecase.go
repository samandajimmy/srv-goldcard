package billings

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// UseCase represent the product requirements usecases
type UseCase interface {
	GetBillingStatement(c echo.Context, pl models.PayloadAccNumber) (models.BillingStatement, error)
}
