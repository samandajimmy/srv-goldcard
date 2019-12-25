package applications

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// Repository represent the application's repository contract
type Repository interface {
	PostSavingAccount(echo.Context, *models.Applications) error
}
