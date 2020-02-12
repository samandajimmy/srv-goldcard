package transactions

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// UseCase represent the transactions usecases
type Repository interface {
	PostBRIPendingTransactions(c echo.Context, trans models.Transaction) error
	GetAccountByBrixKey(c echo.Context, acc *models.Transaction) error
	GetTransactionsHistory(c echo.Context, pt models.PayloadHistoryTransactions) ([]models.ResponseHistoryTransactions, error)
	PostTransactionsANDCardBalance(c echo.Context, trx models.Transaction) error
}
