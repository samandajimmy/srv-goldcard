package registration

import (
	"srv-goldcard/internal/app/model"

	"github.com/labstack/echo"
)

// Repository represent the registration's repository contract
type Repository interface {
	PostAddress(c echo.Context, acc model.Account) error
	PostSavingAccount(echo.Context, model.Account) error
	CreateApplication(echo.Context, model.Applications, model.Account, model.PersonalInformation) error
	GetBankIDByCode(c echo.Context, bankCode string) (int64, error)
	GetAccountByAppNumber(c echo.Context, acc *model.Account) error
	GetAllRegData(c echo.Context, appNumber string) (model.PayloadBriRegister, error)
	UpdateAllRegistrationData(c echo.Context, acc model.Account) error
	GetEmergencyContactIDByType(c echo.Context, typeDef string) (int64, error)
	GetZipcode(c echo.Context, addrData model.AddressData) (string, error)
	GetCityFromZipcode(c echo.Context, zipcode string) (model.AddressData, error)
	UpdateCardLimit(c echo.Context, acc model.Account, fnAfter func() error) error
	UpdateBrixkeyID(c echo.Context, acc model.Account) error
	UpdateAppDocID(c echo.Context, acc model.Applications) error
	GetAppStatus(c echo.Context, app model.Applications) (model.AppStatus, error)
	UpdateAppStatus(c echo.Context, app model.Applications) error
	UpdateApplication(c echo.Context, app model.Applications, col []string) error
	UpsertAppDocument(c echo.Context, app model.Document) error
	PostOccupation(echo.Context, model.Account) error
	GetCoreServiceStatus(c echo.Context) error
	UpdateCoreOpen(c echo.Context, acc *model.Account) error
	GetDocumentByApplicationId(appId int64, docType string) ([]model.Document, error)
	GetSignatoryNameParam(c echo.Context) (string, error)
	GetSignatoryNipParam(c echo.Context) (string, error)
	DeactiveAccount(c echo.Context, acc model.Account) error
	UpdateAppStatusTimeout(echo.Context, model.Applications) error
	ForceUpdateAppStatusTimeout() error
	GetAppOngoing() ([]model.Account, error)
	ForceDeliverAccount(c echo.Context, acc model.Account) error
	ResetAppStatusToCardProcessed(appsId int64) error
	GetAppByCIF(cif string) (model.Applications, error)
}

// RestRepository represent the rest registrations repository contract
type RestRepository interface {
	GetCurrentGoldSTL(c echo.Context) (int64, error)
	OpenGoldcard(c echo.Context, acc model.Account, isRecalculate bool) error
	CloseGoldcard(c echo.Context, acc model.Account) error
	SendNotification(c echo.Context, notif model.PdsNotification, notifType string) error
}
