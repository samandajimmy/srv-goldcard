package registrations

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// Repository represent the campaigntrx's repository contract
type Repository interface {
	PostAddress(echo.Context, *models.Registrations) error
}
