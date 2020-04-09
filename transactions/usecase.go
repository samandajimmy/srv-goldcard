package transactions

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// UseCase represent the transactions usecases
type UseCase interface {
	PostBRIPendingTransactions(c echo.Context, pbpt models.PayloadBRIPendingTransactions) models.ResponseErrors
	GetCardBalance(c echo.Context, pl models.PayloadAccNumber) (models.BRICardBalance, error)
	GetTransactionsHistory(c echo.Context, plListTrx models.PayloadListTrx) (interface{}, models.ResponseErrors)
	CheckAccountByAccountNumber(c echo.Context, pl interface{}) (models.Account, error)
	CheckAccountByBrixkey(c echo.Context, pl interface{}) (models.Account, error)
	UpdateAndGetCardBalance(c echo.Context, acc models.Account) (models.Card, error)
	PostPaymentTransaction(c echo.Context, pl models.PayloadPaymentTransactions) models.ResponseErrors
	PaymentInquiry(c echo.Context, ppi models.PlPaymentInquiry) (string, models.ResponseErrors)
	PostPaymentTrxCore(c echo.Context, pl models.PlPaymentTrxCore) models.ResponseErrors
	DecreasedSTL(c echo.Context, pl models.PayloadCoreDecreasedSTL) models.ResponseErrors
}
