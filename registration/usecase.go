package registration

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// UseCase represent the rewards usecases
type UseCase interface {
	TestRegistration(echo.Context, *models.Registration) error
}
