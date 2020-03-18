package usecase

import (
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
	var notif models.PdsNotification
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

	// mapping all trx data needed
	err = trx.MappingTransactions(c, pl, trx, refTrxPgdn.String(), currStl)

	if err != nil {
		errors.SetTitle(models.ErrMappingData.Error())
		return errors
	}

	// store trx to db
	err = trxUS.trxRepo.PostTransactions(c, trx)

	if err != nil {
		errors.SetTitle(models.ErrInsertTransactions.Error())
		return errors
	}

	// update card balance by account
	go func() {
		_, err := trxUS.UpdateAndGetCardBalance(c, trx.Account)

		if err != nil {
			logger.Make(c, nil).Debug(err)
		}
	}()

	// Send notification to user in pds
	go func() {
		notif.GcTransaction(trx)
		_ = trxUS.rrRepo.SendNotification(c, notif, "mobile")
	}()

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

	// Get trx Account by BrixKey
	trx, err := trxUS.trxRepo.GetTrxAccountByBrixKey(c, BrixKey.String())

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
	// check account number
	acc, err := trxUS.CheckAccountByAccountNumber(c, pl)

	if err != nil {
		return briCardBal, err
	}

	// update and get card balance by account
	card, err := trxUS.UpdateAndGetCardBalance(c, acc)

	if err != nil {
		return briCardBal, models.ErrGetCardBalance
	}

	// mapping + return bri card balance
	return models.BRICardBalance{
		CurrentBalance: float64(card.Balance),
		CreditLimit:    float64(card.CardLimit),
	}, nil
}

func (trxUS *transactionsUseCase) UpdateAndGetCardBalance(c echo.Context, acc models.Account) (models.Card, error) {
	// define channel buffer
	errPromise := make(chan error)
	briCardBal := make(chan models.BRICardBalance)
	currStl := make(chan int64)

	go func() {
		// get balance from bank
		cardBal, err := trxUS.trxrRepo.GetBRICardInformation(c, acc)

		if err != nil {
			briCardBal <- models.BRICardBalance{}
			errPromise <- models.ErrGetCardBalance
			return
		}

		briCardBal <- cardBal
		errPromise <- nil
	}()

	go func() {
		// get current gold price
		stl, err := trxUS.rrRepo.GetCurrentGoldSTL(c)

		if err != nil {
			currStl <- 0
			errPromise <- models.ErrGetCurrSTL
			return
		}

		currStl <- stl
		errPromise <- nil
	}()

	// get current stl, card balance and error promises
	stl := <-currStl
	cardBal := <-briCardBal
	err := <-errPromise

	// check error promises
	if err != nil {
		return models.Card{}, err
	}

	// get gold balance
	goldBalance := acc.Card.ConvertMoneyToGold(int64(cardBal.CurrentBalance), stl)

	// define new card balances
	acc.Card.Balance = int64(cardBal.CurrentBalance)
	acc.Card.GoldBalance = goldBalance
	acc.Card.StlBalance = stl

	// update card balances
	err = trxUS.trxRepo.UpdateCardBalance(c, acc.Card)

	if err != nil {
		return models.Card{}, models.ErrUpdateCardBalance
	}

	return acc.Card, nil
}
