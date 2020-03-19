package registrations

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// UseCase represent the registrations usecases
type UseCase interface {
	PostAddress(echo.Context, models.PayloadAddress) error
	PostSavingAccount(echo.Context, models.PayloadSavingAccount) error
	PostPersonalInfo(echo.Context, models.PayloadPersonalInformation) error
	PostRegistration(echo.Context, models.PayloadRegistration) (models.RespRegistration, error)
	PostCardLimit(c echo.Context, pl models.PayloadCardLimit) error
	FinalRegistration(c echo.Context, pl models.PayloadAppNumber) error
	GetAppStatus(c echo.Context, pl models.PayloadAppNumber) (models.AppStatus, error)
	PostOccupation(echo.Context, models.PayloadOccupation) error
	CheckApplication(c echo.Context, pl interface{}) (models.Account, error)
}
