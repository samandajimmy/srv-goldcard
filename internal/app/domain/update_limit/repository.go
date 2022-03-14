package update_limit

import (
	"srv-goldcard/internal/app/model"

	"github.com/labstack/echo"
)

// Repository represent the transactions Repository
type Repository interface {
	GetEmailByKey(c echo.Context) (string, error)
	GetLastLimitUpdate(c echo.Context, accId int64) (model.LimitUpdate, error)
	GetAccountBySavingAccount(c echo.Context, savingAcc string) (model.Account, error)
	InsertUpdateCardLimit(c echo.Context, limitUpdt model.LimitUpdate) error
	GetLimitUpdate(c echo.Context, refId string) (model.LimitUpdate, error)
	UpdateCardLimitData(c echo.Context, limitUpdt model.LimitUpdate) error
	GetsertGtePayment(c echo.Context, pl model.PayloadCoreGtePayment) (model.GtePayment, error)
	UpdateGtePayment(c echo.Context, gtePayment model.GtePayment, cols []string) error
	GetUpdateLimitInquiriesClosedDate(c echo.Context) (string, error)
}

// Repository represent the update limits Rest Repository
type RestRepository interface {
	CorePostUpdateLimit(c echo.Context, savingAccNum string, card model.Card, cif string) error
	BRIPostUpdateLimit(c echo.Context, acc model.Account, doc model.Document) error
	CorePostInquiryUpdateLimit(c echo.Context, cif string, savingAccNum string, nominalLimit int64) string
}
