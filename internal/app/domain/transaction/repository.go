package transaction

import (
	"srv-goldcard/internal/app/model"

	"github.com/labstack/echo"
)

// Repository represent the transactions Repository
type Repository interface {
	GetAccountByBrixKey(c echo.Context, brixkey string) (model.Account, error)
	PostTransactions(c echo.Context, trx model.Transaction) error
	GetAccountByAccountNumber(c echo.Context, acc *model.Account) error
	UpdateCardBalance(c echo.Context, card model.Card) error
	PostPayment(c echo.Context, trx model.Transaction, bill model.Billing) error
	PostPaymentInquiry(c echo.Context, paymentInq model.PaymentInquiry) error
	GetPayInquiryByRefTrx(c echo.Context, acc model.Account, refTrx string) (model.PaymentInquiry, error)
	UpdatePayInquiryStatusPaid(c echo.Context, pay model.PaymentInquiry) error
	GetAllActiveAccount(c echo.Context) ([]model.Account, error)
	GetPaymentInquiryNotificationData(c echo.Context, pi model.PaymentInquiry) (model.PaymentInquiryNotificationData, error)
	GetAccountByCIF(c echo.Context, acc *model.Account) error
}

// RestRepository represent the rest transactions repository contract
type RestRepository interface {
	GetBRIAppStatus(c echo.Context, brixkey string) (model.BRIAppStatus, error)
	GetBRICardInformation(c echo.Context, acc model.Account) (model.BRICardBalance, error)
	CorePaymentInquiry(c echo.Context, pl model.PlPaymentInquiry, acc model.Account) (map[string]interface{}, error)
	PostPaymentTransactionToCore(c echo.Context, bill model.Billing) error
	PostPaymentBRI(c echo.Context, acc model.Account, amount int64) error
	PostPaymentCoreNotif(c echo.Context, acc model.Account, pl model.PlPaymentTrxCore) error
	GetBRIPendingTrx(c echo.Context, acc model.Account, startDate string, endDate string) (model.RespBRIPendingTrxData, error)
	GetBRIPostedTrx(c echo.Context, briXkey string) (model.RespBRIPostedTransaction, error)
}
