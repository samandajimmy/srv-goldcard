package transactions

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// UseCase represent the transactions usecases
type UseCase interface {
	PostBRIPendingTransactions(c echo.Context, pbpt models.PayloadBRIPendingTransactions) models.ResponseErrors
	GetCardBalance(c echo.Context, pl models.PayloadAccNumber) (models.BRICardBalance, error)
}
