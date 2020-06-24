package billings

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// Repository represent the billings repository contract
type Repository interface {
	GetBillingInquiry(c echo.Context, bill *models.Billing) error
	GetMinPaymentParam(c echo.Context) (float64, error)
	GetDueDateParam(c echo.Context) (int, error)
	PostPegadaianBillings(c echo.Context, pgdBil models.PegadaianBilling) error
}

// RestRepository represent the rest transactions repository contract
type RestRepository interface {
	GetBillingsStatement(c echo.Context, acc models.Account) (map[string]interface{}, string)
}
