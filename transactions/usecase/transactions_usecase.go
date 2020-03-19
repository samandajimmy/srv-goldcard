package usecase

import (
	"gade/srv-goldcard/billings"
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
	billRepo billings.Repository
	trxrRepo transactions.RestRepository
	rrRepo   registrations.RestRepository
}

// TransactionsUseCase represent Transactions Use Case
func TransactionsUseCase(trxRepo transactions.Repository, billRepo billings.Repository,
	trxrRepo transactions.RestRepository, rrRepo registrations.RestRepository) transactions.UseCase {
	return &transactionsUseCase{trxRepo, billRepo, trxrRepo, rrRepo}
}

func (trxUS *transactionsUseCase) PostBRIPendingTransactions(c echo.Context, pl models.PayloadBRIPendingTransactions) models.ResponseErrors {
	var errors models.ResponseErrors
	var notif models.PdsNotification
	acc, err := trxUS.CheckAccountByBrixkey(c, pl)

	if err != nil {
		errors.SetTitle(models.ErrGetAccByBrixkey.Error())
		return errors
	}

	// init account trx
	trx := models.Transaction{AccountId: acc.ID, Account: acc}
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

func (trxUS *transactionsUseCase) PostPaymentTransaction(c echo.Context, pl models.PayloadPaymentTransactions) models.ResponseErrors {
	// TODO: do we need push notif or email?
	// TODO: do we need push to core?
	// TODO: we need to add function update billing status when bill is fully paid?
	// TODO: do we need to concurrent to optim response time?
	var errors models.ResponseErrors
	acc, err := trxUS.CheckAccountByBrixkey(c, pl)

	if err != nil {
		errors.SetTitle(models.ErrGetAccByBrixkey.Error())
		return errors
	}

	// init account trx
	trx := models.Transaction{AccountId: acc.ID, Account: acc}
	// Generate ref transactions pegadaian
	refTrxPgdn, _ := uuid.NewRandom()
	// Get curr STL
	currStl, err := trxUS.rrRepo.GetCurrentGoldSTL(c)

	if err != nil {
		errors.SetTitle(models.ErrGetCurrSTL.Error())
		return errors
	}

	// populate payment transaction for insert
	err = trx.MappingPaymentTransaction(c, pl, trx, refTrxPgdn.String(), currStl)

	if err != nil {
		errors.SetTitle(models.ErrMappingData.Error())
		return errors
	}

	// init current billing
	bill := models.Billing{
		Account: acc,
	}

	// get last published billing to map payment transaction into billings
	err = trxUS.billRepo.GetBillingInquiry(c, &bill)

	if err != nil {
		errors.SetTitle(err.Error())
		return errors
	}

	// prepare billing debt
	// get debt amount
	bill.DebtAmount = bill.DebtAmount - trx.Nominal
	// get debt gold amount
	bill.DebtGold = bill.Account.Card.ConvertMoneyToGold(bill.DebtAmount, currStl)
	// set debt stl
	bill.DebtSTL = currStl

	// insert payment transaction
	err = trxUS.trxRepo.PostPayment(c, trx, bill)

	if err != nil {
		errors.SetTitle(models.ErrInsertPaymentTransactions.Error())
		return errors
	}

	// update card balance to BRI after success receive billing payment
	go func() {
		_, err = trxUS.UpdateAndGetCardBalance(c, trx.Account)

		if err != nil {
			logger.Make(c, nil).Debug(err)
		}
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

func (trxUS *transactionsUseCase) CheckAccountByBrixkey(c echo.Context, pl interface{}) (models.Account, error) {
	r := reflect.ValueOf(pl)
	BrixKey := r.FieldByName("BrixKey")

	// Get trx Account by BrixKey
	acc, err := trxUS.trxRepo.GetAccountByBrixKey(c, BrixKey.String())

	if err != nil {
		return models.Account{}, models.ErrGetAccByBrixkey
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
