package registrations

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// Repository represent the registration's repository contract
type Repository interface {
	PostAddress(c echo.Context, acc models.Account) error
	PostSavingAccount(echo.Context, models.Account) error
	CreateApplication(echo.Context, models.Applications, models.Account, models.PersonalInformation) error
	GetBankIDByCode(c echo.Context, bankCode string) (int64, error)
	GetAccountByAppNumber(c echo.Context, acc *models.Account) error
	GetAllRegData(c echo.Context, appNumber string) (models.PayloadBriRegister, error)
	UpdateAllRegistrationData(c echo.Context, acc models.Account) error
	GetEmergencyContactIDByType(c echo.Context, typeDef string) (int64, error)
	GetZipcode(c echo.Context, addrData models.AddressData) (string, error)
	GetCityFromZipcode(c echo.Context, zipcode string) (models.AddressData, error)
	UpdateCardLimit(c echo.Context, acc models.Account, fnAfter func() error) error
	UpdateBrixkeyID(c echo.Context, acc models.Account) error
	UpdateAppDocID(c echo.Context, acc models.Applications) error
	GetAppStatus(c echo.Context, app models.Applications) (models.AppStatus, error)
	UpdateAppStatus(c echo.Context, app models.Applications) error
	UpdateApplication(c echo.Context, app models.Applications, col []string) error
	UpsertAppDocument(c echo.Context, app models.Document) error
	PostOccupation(echo.Context, models.Account) error
	GetCoreServiceStatus(c echo.Context) error
	UpdateCoreOpen(c echo.Context, acc *models.Account) error
	GetDocumentByApplicationId(appId int64, docType string) ([]models.Document, error)
	GetSignatoryNameParam(c echo.Context) (string, error)
	GetSignatoryNipParam(c echo.Context) (string, error)
	DeactiveAccount(c echo.Context, acc models.Account) error
	UpdateAppStatusTimeout(echo.Context, models.Applications) error
	ForceUpdateAppStatusTimeout() error
	GetAppOngoing() ([]models.Account, error)
	ForceDeliverAccount(c echo.Context, acc models.Account) error
	ResetAppStatusToCardProcessed(appsId int64) error
	GetAppByCIF(cif string) (models.Applications, error)
}

// RestRepository represent the rest registrations repository contract
type RestRepository interface {
	GetCurrentGoldSTL(c echo.Context) (int64, error)
	OpenGoldcard(c echo.Context, acc models.Account, isRecalculate bool) error
	CloseGoldcard(c echo.Context, acc models.Account) error
	SendNotification(c echo.Context, notif models.PdsNotification, notifType string) error
}
