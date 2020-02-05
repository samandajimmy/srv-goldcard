package transactions

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// UseCase represent the registrations usecases
type Repository interface {
	PostBRIPendingTransactions(c echo.Context, trx models.Transactions) error
}
