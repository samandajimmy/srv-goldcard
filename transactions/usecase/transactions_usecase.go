package usecase

import (
	"encoding/json"
	"gade/srv-goldcard/billings"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/registrations"
	"gade/srv-goldcard/transactions"
	"reflect"
	"strconv"

	"github.com/labstack/echo"
)

type transactionsUseCase struct {
	trxRepo  transactions.Repository
	billRepo billings.Repository
	trxrRepo transactions.RestRepository
	rRepo    registrations.Repository
	rrRepo   registrations.RestRepository
}

// TransactionsUseCase represent Transactions Use Case
func TransactionsUseCase(trxRepo transactions.Repository, billRepo billings.Repository,
	trxrRepo transactions.RestRepository, rRepo registrations.Repository, rrRepo registrations.RestRepository) transactions.UseCase {
	return &transactionsUseCase{trxRepo, billRepo, trxrRepo, rRepo, rrRepo}
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
	// Get curr STL
	trx.CurrStl, err = trxUS.rrRepo.GetCurrentGoldSTL(c)

	if err != nil {
		errors.SetTitle(models.ErrGetCurrSTL.Error())
		return errors
	}

	// mapping all trx data needed
	err = trx.MappingTrx(pl, models.TypeTrxCredit, true)

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

	// prepare account trx and account billing
	trx, bill, errors := trxUS.prepareTrxAndBill(c, acc, pl)

	if errors.Code != "00" {
		errors.SetTitle(err.Error())
		return errors
	}

	// prepare billing debt
	_ = trxUS.payTheBill(c, &bill, trx)
	// insert payment transaction
	err = trxUS.trxRepo.PostPayment(c, trx, bill)

	if err != nil {
		errors.SetTitle(models.ErrInsertPaymentTransactions.Error())
		return errors
	}

	// post payment transaction to core
	err = trxUS.trxrRepo.PostPaymentTransactionToCore(c, trx, acc)
	if err != nil {
		errors.SetTitle(models.ErrPostPaymentTransactionToCore.Error())
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

func (trxUS *transactionsUseCase) PostPaymentTrxCore(c echo.Context, pl models.PlPaymentTrxCore) models.ResponseErrors {
	var errors models.ResponseErrors
	acc, err := trxUS.CheckAccountByAccountNumber(c, pl)

	if err != nil {
		errors.SetTitle(models.ErrGetAccByAccountNumber.Error())
		return errors
	}

	// get payment trx data with ref_trx
	payment, err := trxUS.trxRepo.GetPayInquiryByRefTrx(c, acc, pl.RefTrx)

	if err != nil {
		errors.SetTitle(err.Error())
		return errors
	}

	// init payment amount and source
	pl.PaymentAmount = payment.Nominal
	pl.Source = models.SourceCore
	// prepare account trx and account billing
	trx, bill, errors := trxUS.prepareTrxAndBill(c, acc, pl)

	if errors.Title != "" {
		errors.SetTitle(err.Error())
		return errors
	}

	bill = payment.Billing
	// prepare billing debt
	_ = trxUS.payTheBill(c, &bill, trx)
	// insert payment transaction
	err = trxUS.trxRepo.PostPayment(c, trx, bill)

	if err != nil {
		errors.SetTitle(models.ErrInsertPaymentTransactions.Error())
		return errors
	}

	// update payment inquiry status to paid
	err = trxUS.trxRepo.UpdatePayInquiryStatusPaid(c, payment)

	if err != nil {
		errors.SetTitle(err.Error())
		return errors
	}

	return errors
}

func (trxUS *transactionsUseCase) GetTransactionsHistory(c echo.Context, plListTrx models.PayloadListTrx) (interface{}, models.ResponseErrors) {
	var errors models.ResponseErrors

	// Get Account by Account Number
	acc, err := trxUS.CheckAccountByAccountNumber(c, plListTrx)

	if err != nil {
		errors.SetTitle(models.ErrGetAccByAccountNumber.Error())

		return models.ResponseListTrx{}, errors
	}

	result, err := trxUS.trxRepo.GetPgTransactionsHistory(c, acc, plListTrx)

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
		AvailableCredit: card.Balance,
		CreditLimit:     float64(card.CardLimit),
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
	goldBalance := acc.Card.ConvertMoneyToGold(cardBal.AvailableCredit, stl)

	// define new card balances
	acc.Card.Balance = cardBal.AvailableCredit
	acc.Card.GoldBalance = goldBalance
	acc.Card.StlBalance = stl

	// update card balances
	err = trxUS.trxRepo.UpdateCardBalance(c, acc.Card)

	if err != nil {
		return models.Card{}, models.ErrUpdateCardBalance
	}

	return acc.Card, nil
}

func (trxUS *transactionsUseCase) PaymentInquiry(c echo.Context, pl models.PlPaymentInquiry) (map[string]interface{}, models.ResponseErrors) {
	var errors models.ResponseErrors
	response := map[string]interface{}{}
	// Get Account by Account Number
	acc, err := trxUS.CheckAccountByAccountNumber(c, pl)

	if err != nil {
		errors.SetTitle(models.ErrGetAccByAccountNumber.Error())
		return response, errors
	}

	// init current billing
	bill := models.Billing{Account: acc}
	// get last published billing to map payment transaction into billings
	err = trxUS.billRepo.GetBillingInquiry(c, &bill)

	if err != nil {
		errors.SetTitleCode("11", models.ErrNoBilling.Error(), "")
		return response, errors
	}

	// payment inquiry to core
	respInquiry, err := trxUS.trxrRepo.CorePaymentInquiry(c, pl)

	if err != nil {
		errors.SetTitleCode("11", models.ErrNoBilling.Error(), "")
		return response, errors
	}

	if _, ok := respInquiry["reffSwitching"].(string); !ok {
		errors.SetTitle(models.ErrSetVar.Error())

		return response, errors
	}

	// get refTrx
	refTrx := respInquiry["reffSwitching"].(string)
	// convert respInquiry to json
	respJSON, err := json.Marshal(respInquiry)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		errors.SetTitle(err.Error())

		return response, errors
	}

	// check over payment
	if bill.DebtAmount < pl.PaymentAmount {
		errors.SetTitleCode("22", models.ErrOverPayment.Error(), strconv.FormatInt(bill.DebtAmount, 10))
		return response, errors
	}

	// check payment less than 10% remaining payment but only at the first payment
	if bill.DebtAmount == bill.Amount && pl.PaymentAmount <= int64(bill.MinimumPayment) {
		errors.SetTitleCode("22", models.ErrMinimumPayment.Error(), strconv.FormatInt(bill.DebtAmount, 10))
		return response, errors
	}

	// prepare the payment inquiry data
	paymentInq := models.PaymentInquiry{
		AccountId:    acc.ID,
		BillingId:    bill.ID,
		RefTrx:       refTrx,
		Nominal:      pl.PaymentAmount,
		CoreResponse: respJSON,
	}

	// insert payment inquiry
	err = trxUS.trxRepo.PostPaymentInquiry(c, paymentInq)

	if err != nil {
		errors.SetTitle(models.ErrInsertPaymentTransactions.Error())
		return response, errors
	}

	return respInquiry, errors
}

func (trxUS *transactionsUseCase) prepareTrxAndBill(c echo.Context, acc models.Account, pl interface{}) (models.Transaction,
	models.Billing, models.ResponseErrors) {
	var errors = models.ResponseErrors{Code: "00"}
	var err error
	// init account trx
	trx := models.Transaction{AccountId: acc.ID, Account: acc}
	// init current billing
	bill := models.Billing{Account: acc}
	// Get curr STL
	trx.CurrStl, err = trxUS.rrRepo.GetCurrentGoldSTL(c)

	if err != nil {
		errors.SetTitle(models.ErrGetCurrSTL.Error())
		return models.Transaction{}, models.Billing{}, errors
	}

	// populate payment transaction for insert
	err = trx.MappingTrx(pl, models.TypeTrxDebit, false)

	if err != nil {
		errors.SetTitle(models.ErrMappingData.Error())
		return models.Transaction{}, models.Billing{}, errors
	}

	// get last published billing to map payment transaction into billings
	err = trxUS.billRepo.GetBillingInquiry(c, &bill)

	if err != nil {
		errors.SetTitleCode("11", models.ErrNoBilling.Error(), "")
		return models.Transaction{}, models.Billing{}, errors
	}

	return trx, bill, errors
}

func (trxUS *transactionsUseCase) payTheBill(c echo.Context, bill *models.Billing, trx models.Transaction) error {
	// prepare billing debt
	// get debt amount
	bill.DebtAmount = bill.DebtAmount - trx.Nominal

	// change bill status to paid
	if bill.DebtAmount <= 0 && bill.Status != models.BillTrxPaid {
		bill.Status = models.BillTrxPaid
	}

	// get debt gold amount
	bill.DebtGold = bill.Account.Card.ConvertMoneyToGold(bill.DebtAmount, trx.CurrStl)
	// set debt stl
	bill.DebtSTL = trx.CurrStl

	return nil
}

// DecreasedSTL is a func to recalculate gold card rupiah limit when occurs stl decreased equal or more than 5%
func (trxUS *transactionsUseCase) DecreasedSTL(c echo.Context, pcds models.PayloadCoreDecreasedSTL) models.ResponseErrors {
	var errors models.ResponseErrors
	var notif models.PdsNotification
	var oldCard models.Card

	// check if payload decreased five percent is false then return
	if pcds.DecreasedFivePercent != "true" {
		return errors
	}

	// Get CurrentStl from Core payload
	currStl := pcds.STL

	// Get All Active Account
	allAccs, err := trxUS.trxRepo.GetAllActiveAccount(c)

	if err != nil {
		errors.SetTitle(models.ErrGetAccByAccountNumber.Error())
		return errors
	}

	for _, acc := range allAccs {
		notif = models.PdsNotification{}
		oldCard = acc.Card

		// set card limit
		err = acc.Card.SetCardLimit(currStl)

		if err != nil {
			logger.Make(c, nil).Debug(err)
			continue
		}

		// update card limit in db
		refId, err := trxUS.rRepo.UpdateCardLimit(c, acc, true)

		if err != nil {
			continue
		}

		// Send notification to user in pds
		notif.GcDecreasedSTL(acc, oldCard, refId)
		_ = trxUS.rrRepo.SendNotification(c, notif, "mobile")
	}

	return errors
}
