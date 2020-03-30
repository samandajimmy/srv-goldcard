package registrations

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// Repository represent the registration's repository contract
type Repository interface {
	PostAddress(echo.Context, models.Account) error
	PostSavingAccount(echo.Context, models.Account) error
	CreateApplication(echo.Context, models.Applications, models.Account, models.PersonalInformation) error
	GetBankIDByCode(c echo.Context, bankCode string) (int64, error)
	GetAccountByAppNumber(c echo.Context, acc *models.Account) error
	GetAllRegData(c echo.Context, appNumber string) (models.PayloadBriRegister, error)
	UpdateAllRegistrationData(c echo.Context, acc models.Account) error
	GetEmergencyContactIDByType(c echo.Context, typeDef string) (int64, error)
	GetZipcode(c echo.Context, addrData models.AddressData) (string, error)
	GetCityFromZipcode(c echo.Context, acc models.Account) (string, string, error)
	UpdateCardLimit(c echo.Context, acc models.Account) error
	UpdateBrixkeyID(c echo.Context, acc models.Account) error
	UpdateAppDocID(c echo.Context, acc models.Applications) error
	GetAppStatus(c echo.Context, app models.Applications) (models.AppStatus, error)
	UpdateAppStatus(c echo.Context, app models.Applications) error
	UpdateApplication(c echo.Context, app models.Applications, col []string) error
	UpsertAppDocument(c echo.Context, app models.Document) error
	PostOccupation(echo.Context, models.Account) error
	GetCoreServiceStatus(c echo.Context) error
	UpdateAppError(c echo.Context, appNumber, processID string, errStatus bool) error
}

// RestRepository represent the rest registrations repository contract
type RestRepository interface {
	GetCurrentGoldSTL(c echo.Context) (int64, error)
	OpenGoldcard(c echo.Context, acc models.Account, isRecalculate bool) error
	SendNotification(c echo.Context, notif models.PdsNotification, notifType string) error
}
