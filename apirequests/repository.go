package apirequests

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// Repository represent the apirequest's repository contract
type Repository interface {
	InserAPIRequest(echo.Context, models.APIRequest) error
}
