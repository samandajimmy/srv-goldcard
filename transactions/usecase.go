package transactions

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// UseCase represent the registrations usecases
type UseCase interface {
	PostBRIPendingTransactions(c echo.Context, pbpt models.PayloadBRIPendingTransactions) models.ResponseErrors
}
