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
	trxRepo transactions.Repository
	rrRepo  registrations.RestRepository
}

// TransactionsUseCase represent Transactions Use Case
func TransactionsUseCase(trxRepo transactions.Repository, rrRepo registrations.RestRepository) transactions.UseCase {
	return &transactionsUseCase{trxRepo, rrRepo}
}

func (trxUS *transactionsUseCase) PostBRIPendingTransactions(c echo.Context, pl models.PayloadBRIPendingTransactions) models.ResponseErrors {
	var errors models.ResponseErrors
	trx, err := trxUS.checkAccount(c, pl)

	if err != nil {
		errors.SetTitle(models.ErrGetAccByBrixkey.Error())
		return errors
	}
	// Generate ref transactions pegadaian
	refTrxPgdn, _ := uuid.NewRandom()

	// Get curr STL
	currStl, err := trxUS.rrRepo.GetCurrentGoldSTL(c)

	if err != nil {
		errors.SetTitle(models.ErrGetCurrSTL.Error())
		return errors
	}

	err = trx.MappingTransactions(c, pl, trx, refTrxPgdn.String(), currStl)

	if err != nil {
		errors.SetTitle(models.ErrMappingData.Error())
		return errors
	}

	err = trxUS.trxRepo.PostTransactions(c, trx)

	if err != nil {
		errors.SetTitle(models.ErrInsertTransactions.Error())
		return errors
	}

	return errors
}

func (trxUS *transactionsUseCase) GetTransactionsHistory(c echo.Context, pht models.PayloadHistoryTransactions) (interface{}, models.ResponseErrors) {
	var errors models.ResponseErrors
	if pht.Pagination.Limit != 0 {
		return trxUS.GetPgTransactionsHistory(c, pht)
	}

	result, err := trxUS.trxRepo.GetAllTransactionsHistory(c, pht)

	if err != nil {
		errors.SetTitle(models.ErrGetHistoryTransactions.Error())
		return result, errors
	}

	return result, errors
}

func (trxUS *transactionsUseCase) checkAccount(c echo.Context, pl interface{}) (models.Transaction, error) {
	r := reflect.ValueOf(pl)
	BrixKey := r.FieldByName("BrixKey")

	// Get Account by BrixKey
	trx := models.Transaction{Account: models.Account{BrixKey: BrixKey.String()}}
	err := trxUS.trxRepo.GetAccountByBrixKey(c, &trx)

	if err != nil {
		return models.Transaction{}, models.ErrGetAccByBrixkey
	}

	return trx, nil
}

func (trxUS *transactionsUseCase) GetPgTransactionsHistory(c echo.Context, pht models.PayloadHistoryTransactions) (interface{}, models.ResponseErrors) {
	var errors models.ResponseErrors
	result, err := trxUS.trxRepo.GetPgTransactionsHistory(c, pht)

	if err != nil {
		errors.SetTitle(models.ErrGetHistoryTransactions.Error())
		return result, errors
	}

	return result, errors
}
