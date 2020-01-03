package registrations

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// Repository represent the campaigntrx's repository contract
type Repository interface {
	PostAddress(echo.Context, *models.Registrations) error
	GetAddress(echo.Context, string) (string, error)
	PostSavingAccount(echo.Context, *models.Applications) error
	CreateApplication(echo.Context, models.Applications, models.Account, models.PersonalInformation) error
	GetBankIDByCode(c echo.Context, bankCode string) (int64, error)
}
