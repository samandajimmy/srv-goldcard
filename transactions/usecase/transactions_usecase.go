package usecase

import (
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/registrations"
	"gade/srv-goldcard/transactions"
	"reflect"

	"github.com/google/uuid"
	"github.com/labstack/echo"
)

type transactionsUseCase struct {
	trRepo transactions.Repository
	rrRepo registrations.RestRepository
}

// TransactionsUseCase represent Transactions Use Case
func TransactionsUseCase(trRepo transactions.Repository, rrRepo registrations.RestRepository) transactions.UseCase {
	return &transactionsUseCase{trRepo, rrRepo}
}

func (trx *transactionsUseCase) PostBRIPendingTransactions(c echo.Context, pl models.PayloadBRIPendingTransactions) models.ResponseErrors {
	var errors models.ResponseErrors
	trans, err := trx.checkAccount(c, pl)

	if err != nil {
		errors.SetTitle(models.ErrGetAccByBrixkey.Error())
		return errors
	}
	// Generate ref transactions pegadaian
	refTrxPg, _ := uuid.NewRandom()

	// Get curr STL
	currStl, err := trx.rrRepo.GetCurrentGoldSTL(c)

	if err != nil {
		errors.SetTitle(models.ErrGetCurrSTL.Error())
		return errors
	}

	err = trans.MappingTransactions(c, pl, trans, refTrxPg.String(), currStl)

	if err != nil {
		errors.SetTitle(models.ErrMappingData.Error())
		return errors
	}

	err = trx.trRepo.PostBRIPendingTransactions(c, trans)

	if err != nil {
		errors.SetTitle(models.ErrInsertTransactions.Error())
		return errors
	}

	return errors
}

func (trx *transactionsUseCase) checkAccount(c echo.Context, pl interface{}) (models.Transaction, error) {
	r := reflect.ValueOf(pl)
	BrixKey := r.FieldByName("BrixKey")

	if BrixKey.IsZero() {
		return models.Transaction{}, nil
	}

	trans := models.Transaction{Account: models.Account{BrixKey: BrixKey.String()}}
	err := trx.trRepo.GetAccountByBrixKey(c, &trans)

	if err != nil {
		return models.Transaction{}, models.ErrGetAccByBrixkey
	}

	return trans, nil
}
