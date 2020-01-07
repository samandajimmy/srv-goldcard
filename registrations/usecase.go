package registrations

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// UseCase represent the rewards usecases
type UseCase interface {
	PostAddress(echo.Context, *models.Registrations) error
	GetAddress(echo.Context, string) (map[string]interface{}, error)
	PostSavingAccount(echo.Context, *models.Applications) error
	PostPersonalInfo(echo.Context, models.PayloadPersonalInformation) error
	PostRegistration(echo.Context, models.PayloadRegistration) (string, error)
}
