package billings

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// UseCase represent the product requirements usecases
type UseCase interface {
	GetBillingStatement(c echo.Context, pl models.PayloadAccNumber) (models.BillingStatement, error)
	PostBRIPegadaianBillings(c echo.Context, pbpb models.PayloadBRIPegadaianBillings) models.ResponseErrors
	PaymentInquiry(c echo.Context, ppi models.PayloadPaymentInquiry) models.ResponseErrors
}
