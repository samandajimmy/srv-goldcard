package transactions

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// UseCase represent the transactions usecases
type Repository interface {
	GetAccountByBrixKey(c echo.Context, acc *models.Transaction) error
	GetAllTransactionsHistory(c echo.Context, pt models.PayloadHistoryTransactions) (models.ResponseHistoryTransactions, error)
	GetPgTransactionsHistory(c echo.Context, pt models.PayloadHistoryTransactions) (models.ResponseHistoryTransactions, error)
	PostTransactions(c echo.Context, trx models.Transaction) error
}
