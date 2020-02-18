package usecase

import (
	"encoding/json"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/registrations"
	"gade/srv-goldcard/transactions"
	"reflect"

	"github.com/google/uuid"
	"github.com/labstack/echo"
)

type transactionsUseCase struct {
	trxRepo  transactions.Repository
	trxrRepo transactions.RestRepository
	rrRepo   registrations.RestRepository
}

// TransactionsUseCase represent Transactions Use Case
func TransactionsUseCase(trxRepo transactions.Repository, trxrRepo transactions.RestRepository, rrRepo registrations.RestRepository) transactions.UseCase {
	return &transactionsUseCase{trxRepo, trxrRepo, rrRepo}
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

func (trxUS *transactionsUseCase) CheckAccountByAccountNumber(c echo.Context, pl interface{}) (models.Account, error) {
	r := reflect.ValueOf(pl)
	accNumber := r.FieldByName("AccountNumber")

	// Get Account by Account Number
	acc := models.Account{AccountNumber: accNumber.String()}
	err := trxUS.trxRepo.GetAccountByAccountNumber(c, &acc)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return models.Account{}, models.ErrGetAccByAccountNumber
	}

	return acc, nil
}

func (trxUS *transactionsUseCase) GetCardBalance(c echo.Context, pl models.PayloadAccNumber) (models.BRICardBalance, error) {
	var briCardBal models.BRICardBalance
	acc, err := trxUS.CheckAccountByAccountNumber(c, pl)

	if err != nil {
		return briCardBal, err
	}

	// Request BRI endpoint for check card information
	briCardInfo, err := trxUS.trxrRepo.GetBRICardInformation(c, acc)

	if err != nil {
		return briCardBal, models.ErrGetCardBalance
	}

	mrshlCardInfo, err := json.Marshal(briCardInfo)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return briCardBal, models.ErrGetCardBalance
	}

	err = json.Unmarshal(mrshlCardInfo, &briCardBal)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return briCardBal, models.ErrGetCardBalance
	}

	return briCardBal, err
}
