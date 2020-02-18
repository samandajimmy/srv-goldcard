package transactions

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// Repository represent the transactions Repository
type Repository interface {
	GetAccountByBrixKey(c echo.Context, acc *models.Transaction) error
	GetAllTransactionsHistory(c echo.Context, pt models.PayloadHistoryTransactions) (models.ResponseHistoryTransactions, error)
	GetPgTransactionsHistory(c echo.Context, pt models.PayloadHistoryTransactions) (models.ResponseHistoryTransactions, error)
	PostTransactions(c echo.Context, trx models.Transaction) error
	GetAccountByAccountNumber(c echo.Context, acc *models.Account) error
}

// RestRepository represent the rest transactions repository contract
type RestRepository interface {
	GetBRICardInformation(c echo.Context, acc models.Account) (map[string]interface{}, error)
}
