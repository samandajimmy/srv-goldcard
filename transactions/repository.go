package transactions

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// Repository represent the transactions Repository
type Repository interface {
	GetAccountByBrixKey(c echo.Context, brixkey string) (models.Account, error)
	PostTransactions(c echo.Context, trx models.Transaction) error
	GetAccountByAccountNumber(c echo.Context, acc *models.Account) error
	UpdateCardBalance(c echo.Context, card models.Card) error
	PostPayment(c echo.Context, trx models.Transaction, bill models.Billing) error
	PostPaymentInquiry(c echo.Context, paymentInq models.PaymentInquiry) error
	GetPayInquiryByRefTrx(c echo.Context, acc models.Account, refTrx string) (models.PaymentInquiry, error)
	UpdatePayInquiryStatusPaid(c echo.Context, pay models.PaymentInquiry) error
	GetAllActiveAccount(c echo.Context) ([]models.Account, error)
	GetPaymentInquiryNotificationData(c echo.Context, pi models.PaymentInquiry) (models.PaymentInquiryNotificationData, error)
}

// RestRepository represent the rest transactions repository contract
type RestRepository interface {
	GetBRICardInformation(c echo.Context, acc models.Account) (models.BRICardBalance, error)
	CorePaymentInquiry(c echo.Context, pl models.PlPaymentInquiry, acc models.Account) (map[string]interface{}, error)
	PostPaymentTransactionToCore(c echo.Context, bill models.Billing) error
	PostPaymentBRI(c echo.Context, acc models.Account, amount int64) error
	PostPaymentCoreNotif(c echo.Context, acc models.Account, pl models.PlPaymentTrxCore) error
	GetBRIPendingTrx(c echo.Context, acc models.Account, startDate string, endDate string) (models.RespBRIPendingTrxData, error)
	GetBRIPostedTrx(c echo.Context, briXkey string) (models.RespBRIPostedTransaction, error)
}
