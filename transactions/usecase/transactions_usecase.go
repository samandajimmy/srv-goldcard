package usecase

import (
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/transactions"

	"github.com/labstack/echo"
)

type transactionsUseCase struct {
	trRepo transactions.Repository
}

// RegistrationsUseCase represent Registrations Use Case
func RegistrationsUseCase(trRepo transactions.Repository) transactions.UseCase {
	return &transactionsUseCase{trRepo}
}

func (tuc *transactionsUseCase) PostBRIPendingTransactions(c echo.Context, pl models.PayloadBRIPendingTransactions) error {
	return nil
}
