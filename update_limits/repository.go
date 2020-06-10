package update_limits

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// Repository represent the transactions Repository
type Repository interface {
	GetEmailByKey(c echo.Context) (string, error)
	GetLastLimitUpdate(accId int64) (models.LimitUpdate, error)
	GetAccountBySavingAccount(c echo.Context, savingAcc string) (models.Account, error)
}

// Repository represent the update limits Rest Repository
type RestRepository interface {
	CorePostUpdateLimit(c echo.Context, savingAccNum string, card models.Card, cif string) error
	BRIPostUpdateLimit(c echo.Context, acc models.Account, doc models.Document) error
	CorePostInquiryUpdateLimit(c echo.Context, cif string, savingAccNum string, nominalLimit int64) string
}
