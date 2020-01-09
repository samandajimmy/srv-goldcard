package registrations

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// UseCase represent the rewards usecases
type UseCase interface {
	PostAddress(echo.Context, models.PayloadAddress) error
	PostSavingAccount(echo.Context, models.PayloadSavingAccount) error
	PostPersonalInfo(echo.Context, models.PayloadPersonalInformation) error
	PostRegistration(echo.Context, models.PayloadRegistration) (string, error)
	PostCardLimit(c echo.Context, pl models.PayloadCardLimit) error
}
