package billing

import (
	"srv-goldcard/internal/app/model"

	"github.com/labstack/echo"
)

// Repository represent the billings repository contract
type Repository interface {
	GetBillingInquiry(c echo.Context, bill *model.Billing) error
	GetMinPaymentParam(c echo.Context) (float64, error)
	GetDueDateParam(c echo.Context) (int, error)
	PostPegadaianBillings(c echo.Context, pgdBil model.PegadaianBilling) error
}

// RestRepository represent the rest transactions repository contract
type RestRepository interface {
	GetBillingsStatement(c echo.Context, acc model.Account) (model.BillingStatement, error)
}
