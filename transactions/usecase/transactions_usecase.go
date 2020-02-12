package usecase

import (
	"encoding/json"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/registrations"
	"gade/srv-goldcard/transactions"
	"reflect"
	"strconv"

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
	// currStl, err := trxUS.rrRepo.GetCurrentGoldSTL(c)
	currStl, _ := strconv.ParseInt("560782", 10, 64)

	if err != nil {
		errors.SetTitle(models.ErrGetCurrSTL.Error())
		return errors
	}

	err = trx.MappingTransactions(c, pl, trx, refTrxPgdn.String(), currStl)

	if err != nil {
		errors.SetTitle(models.ErrMappingData.Error())
		return errors
	}

	err = trxUS.trxRepo.PostBRIPendingTransactions(c, trx)

	if err != nil {
		errors.SetTitle(models.ErrInsertTransactions.Error())
		return errors
	}

	return errors
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

func (trxUS *transactionsUseCase) checkAccountByAccountNumber(c echo.Context, pl interface{}) (models.Transaction, error) {
	r := reflect.ValueOf(pl)
	accNumber := r.FieldByName("AccountNumber")

	// Get Account by Account Number
	trx := models.Transaction{Account: models.Account{AccountNumber: accNumber.String()}}
	err := trxUS.trxRepo.GetAccountByAccountNumber(c, &trx)

	if err != nil {
		return models.Transaction{}, models.ErrGetAccByAccountNumber
	}

	return trx, nil
}

func (trxUS *transactionsUseCase) GetCardBalance(c echo.Context, pl models.PayloadAccNumber) (models.BRICardBalance, error) {
	var briCBal models.BRICardBalance
	trx, err := trxUS.checkAccountByAccountNumber(c, pl)

	if err != nil {
		return briCBal, err
	}

	err = trx.MappingTransactionsAccount(c, pl)

	if err != nil {
		return briCBal, err
	}

	// Hit BRI endpoint for check card information
	gBriCInfo, err := trxUS.trxrRepo.GetBRICardInformation(c, trx.Account)

	if err != nil {
		return briCBal, err
	}

	jmGbci, _ := json.Marshal(gBriCInfo)
	json.Unmarshal(jmGbci, &briCBal)

	return briCBal, err
}
