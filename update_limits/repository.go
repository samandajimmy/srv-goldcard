package update_limits

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// Repository represent the transactions Repository
type Repository interface {
	GetEmailByKey(c echo.Context) (string, error)
}

// RestRepository represent the rest billings repository contract
type RestRepository interface {
	// gc to  bri
	PostBRIGtePayment(c echo.Context, bill models.Billing) error
}
