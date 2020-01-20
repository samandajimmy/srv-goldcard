package registrations

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// Repository represent the campaigntrx's repository contract
type Repository interface {
	PostAddress(echo.Context, models.Account) error
	PostSavingAccount(echo.Context, models.Account) error
	CreateApplication(echo.Context, models.Applications, models.Account, models.PersonalInformation) error
	GetBankIDByCode(c echo.Context, bankCode string) (int64, error)
	GetAccountByAppNumber(c echo.Context, appNumber string) (models.Account, error)
	GetAllRegData(c echo.Context, appNumber string) (models.PayloadPersonalInformation, error)
	UpdateAllRegistrationData(c echo.Context, acc models.Account) error
	GetEmergencyContactIDByType(c echo.Context, typeDef string) (int64, error)
	GetZipcode(c echo.Context, addrData models.AddressData) (string, error)
	UpdateCardLimit(c echo.Context, acc models.Account) error
	UpdateBrixkeyID(c echo.Context, acc models.Account) error
	UpdateAppDocID(c echo.Context, acc models.Applications) error
	GetAppByID(c echo.Context, appID int64) (models.Applications, error)
	UpdateGetAppStatus(c echo.Context, app models.Applications) (models.AppStatus, error)
	PostOccupation(echo.Context, models.Account) error
}
