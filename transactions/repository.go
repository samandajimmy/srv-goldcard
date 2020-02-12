package transactions

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// UseCase represent the transactions usecases
type Repository interface {
	PostBRIPendingTransactions(c echo.Context, trans models.Transaction) error
	GetAccountByBrixKey(c echo.Context, acc *models.Transaction) error
	GetAccountByAccountNumber(c echo.Context, acc *models.Transaction) error
}

type RestRepository interface {
	GetBRICardInformation(c echo.Context, acc models.Account) (map[string]interface{}, error)
}
