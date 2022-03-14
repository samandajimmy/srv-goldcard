package transaction

import (
	"srv-goldcard/internal/app/model"

	"github.com/labstack/echo"
)

// UseCase represent the transactions usecases
type UseCase interface {
	PostBRIPendingTransactions(c echo.Context, pbpt model.PayloadBRIPendingTransactions) model.ResponseErrors
	GetCardBalance(c echo.Context, pl model.PayloadAccNumber) (model.BRICardBalance, error)
	GetTransactionsHistory(c echo.Context, plListTrx model.PayloadListTrx) (interface{}, model.ResponseErrors)
	CheckAccountByAccountNumber(c echo.Context, pl interface{}) (model.Account, error)
	CheckAccountByBrixkey(c echo.Context, pl interface{}) (model.Account, error)
	UpdateAndGetCardBalance(c echo.Context, acc model.Account) (model.Card, error)
	PostPaymentTransaction(c echo.Context, pl model.PayloadPaymentTransactions) model.ResponseErrors
	PaymentInquiry(c echo.Context, ppi model.PlPaymentInquiry) (map[string]interface{}, model.ResponseErrors)
	PostPaymentTrxCore(c echo.Context, pl model.PlPaymentTrxCore) model.ResponseErrors
	CheckAccountByCIF(c echo.Context, pl model.PayloadCIF) (model.Account, error)
}
