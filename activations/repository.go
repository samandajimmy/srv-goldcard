package activations

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// Repository represent the activations repository contract
type Repository interface {
	PostActivations(echo.Context, models.Account) error
}
