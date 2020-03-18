package transactions

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// Repository represent the transactions Repository
type Repository interface {
	GetAccountByBrixKey(c echo.Context, brixkey string) (models.Account, error)
	GetAllTransactionsHistory(c echo.Context, pt models.PayloadHistoryTransactions) (models.ResponseHistoryTransactions, error)
	GetPgTransactionsHistory(c echo.Context, pt models.PayloadHistoryTransactions) (models.ResponseHistoryTransactions, error)
	PostTransactions(c echo.Context, trx models.Transaction) error
	GetAccountByAccountNumber(c echo.Context, acc *models.Account) error
	UpdateCardBalance(c echo.Context, card models.Card) error
	PostPayment(c echo.Context, trx models.Transaction, bill models.Billing) error
}

// RestRepository represent the rest transactions repository contract
type RestRepository interface {
	GetBRICardInformation(c echo.Context, acc models.Account) (models.BRICardBalance, error)
}
