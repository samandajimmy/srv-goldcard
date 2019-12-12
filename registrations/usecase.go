package registrations

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// UseCase represent the rewards usecases
type UseCase interface {
	Registrations(echo.Context, *models.Registrations) error
}
