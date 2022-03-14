package billing

import (
	"srv-goldcard/internal/app/model"

	"github.com/labstack/echo"
)

// UseCase represent the product requirements usecases
type UseCase interface {
	GetBillingStatement(c echo.Context, pl model.PayloadAccNumber) (model.BillingStatement, error)
	PostBRIPegadaianBillings(c echo.Context, pbpb model.PayloadBRIPegadaianBillings) model.ResponseErrors
}
