package registration

import (
	"srv-goldcard/internal/app/model"

	"github.com/labstack/echo"
)

// UseCase represent the registrations usecases
type UseCase interface {
	PostAddress(echo.Context, model.PayloadAddress) error
	GetAddress(echo.Context, model.PayloadAppNumber) (model.RespGetAddress, error)
	PostSavingAccount(echo.Context, model.PayloadSavingAccount) error
	PostPersonalInfo(echo.Context, model.PayloadPersonalInformation) error
	PostRegistration(echo.Context, model.PayloadRegistration) (model.RespRegistration, error)
	PostCardLimit(c echo.Context, pl model.PayloadCardLimit) error
	FinalRegistrationScheduler(c echo.Context, pl model.PayloadAppNumber) error
	FinalRegistrationPdsApi(c echo.Context, pl model.PayloadAppNumber) error
	GetAppStatus(c echo.Context, pl model.PayloadAppNumber) (model.AppStatus, error)
	PostOccupation(echo.Context, model.PayloadOccupation) error
	CheckApplication(c echo.Context, pl interface{}) (model.Account, error)
	ResetRegistration(c echo.Context, pl model.PayloadAppNumber) error
	GenerateSlipTEDocument(c echo.Context, acc *model.Account) error
	UploadAppDoc(c echo.Context, brixkey string, doc model.Document) error
	RefreshAppTimeoutJob()
	ForceDeliver(c echo.Context, pl model.PayloadAppNumber) error
}
