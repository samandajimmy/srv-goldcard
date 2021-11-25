package usecase

import (
	"encoding/json"
	"gade/srv-goldcard/billings"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/registrations"
	"gade/srv-goldcard/transactions"
	"reflect"

	"github.com/labstack/echo"
	"github.com/leekchan/accounting"

	"time"
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

	// mapping all trx data needed for notification
	err = trx.MappingTrx(pl, models.TypeTrxCredit, true)

	if err != nil {
		errors.SetTitle(models.ErrMappingData.Error())
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
	err = trxUS.trxrRepo.PostPaymentTransactionToCore(c, bill)
	if err != nil {
		errors.SetTitleCode("22", models.ErrPostPaymentTransactionToCore.Error(), "")
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
	var notif models.PdsNotification
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

	// post payment to bri
	err = trxUS.trxrRepo.PostPaymentBRI(c, acc, pl.PaymentAmount)

	if err != nil {
		errors.SetTitle(models.ErrPostPaymentBRI.Error())
		return errors
	}

	// post payment to core when post to bri succeded
	err = trxUS.trxrRepo.PostPaymentCoreNotif(c, acc, pl)

	if err != nil {
		errors.SetTitle(models.ErrPostPaymentCoreNotif.Error())
		return errors
	}

	// update payment inquiry status to paid
	err = trxUS.trxRepo.UpdatePayInquiryStatusPaid(c, payment)

	if err != nil {
		errors.SetTitle(err.Error())
		return errors
	}

	// Get Payment Inquiry Notification data
	pind, err := trxUS.trxRepo.GetPaymentInquiryNotificationData(c, payment)

	if err != nil {
		errors.SetTitle(err.Error())
		return errors
	}

	// Send notification to user in pds and email
	go func() {
		notif.GcPayment(trx, bill, pind)
		_ = trxUS.rrRepo.SendNotification(c, notif, "email")
		_ = trxUS.rrRepo.SendNotification(c, notif, "mobile")
	}()

	return errors
}

func (trxUS *transactionsUseCase) GetTransactionsHistory(c echo.Context, plListTrx models.PayloadListTrx) (interface{}, models.ResponseErrors) {
	var errors models.ResponseErrors
	var result, pendingArr, postedArr []models.ListTrx
	now := time.Now()
	nowDate := now.Format(models.DateFormat)
	yesterdayDate := now.AddDate(0, 0, -1).Format(models.DateFormat)

	// Get Account by Account Number
	acc, err := trxUS.CheckAccountByAccountNumber(c, plListTrx)

	if err != nil {
		errors.SetTitle(models.ErrGetAccByAccountNumber.Error())

		return models.ResponseListTrx{}, errors
	}

	BRIPending, err := trxUS.trxrRepo.GetBRIPendingTrx(c, acc, yesterdayDate, nowDate)

	if err != nil {
		errors.SetTitle(models.ErrGetHistoryTransactions.Error())
		return models.ResponseListTrx{}, errors
	}

	for _, singleBRIPendigTrx := range BRIPending.TransactionData {
		pendingArr = append(pendingArr, models.ListTrx{
			RefTrx:      "-",
			Nominal:     singleBRIPendigTrx.BillAmount,
			TrxDate:     time.Unix(singleBRIPendigTrx.TransactionDate/1000, 0).Format(models.DateTimeFormatZone),
			Description: singleBRIPendigTrx.Description})
	}

	pendingArr = models.ReverseArray(pendingArr)

	BRIPosted, err := trxUS.trxrRepo.GetBRIPostedTrx(c, acc.BrixKey)

	if err != nil {
		errors.SetTitle(models.ErrGetHistoryTransactions.Error())
		return models.ResponseListTrx{}, errors
	}

	for _, singleBRIPostedTrx := range BRIPosted.ListOfTransactions {
		postedArr = append(postedArr, models.ListTrx{
			RefTrx:      singleBRIPostedTrx.TrxReff,
			Nominal:     singleBRIPostedTrx.TrxAmount,
			TrxDate:     time.Unix(singleBRIPostedTrx.EffectiveDate/1000, 0).Format(models.DateTimeFormatZone),
			Description: singleBRIPostedTrx.TrxDesc})
	}

	postedArr = models.ReverseArray(postedArr)

	result = append(result, pendingArr...)
	result = append(result, postedArr...)

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
	goldBalance := acc.Card.SetGoldLimit(cardBal.AvailableCredit, stl)
	// define previous card limit and balance
	acc.Card.MappingPrevCardData(cardBal)
	// define new card balances
	acc.Card.Balance = cardBal.AvailableCredit
	acc.Card.GoldBalance = goldBalance
	acc.Card.StlBalance = stl
	// set card limit value from response to BRI
	acc.Card.CardLimit = int64(cardBal.CreditLimit)

	// update card balances
	err = trxUS.trxRepo.UpdateCardBalance(c, acc.Card)

	if err != nil {
		return models.Card{}, models.ErrUpdateCardBalance
	}

	return acc.Card, nil
}

func (trxUS *transactionsUseCase) PaymentInquiry(c echo.Context, pl models.PlPaymentInquiry) (map[string]interface{}, models.ResponseErrors) {
	var errors models.ResponseErrors
	var ac = accounting.Accounting{Symbol: "Rp ", Thousand: "."}
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
	respInquiry, err := trxUS.trxrRepo.CorePaymentInquiry(c, pl, acc)

	if err != nil {
		errors.SetTitleCode("11", models.ErrNoBilling.Error(), "")
		return response, errors
	}

	// convert respInquiry to json
	respJSON, err := trxUS.mappingCoreInquiry(c, respInquiry)
	// get refTrx
	refTrx := respInquiry["reffSwitching"].(string)

	if err != nil {
		errors.SetTitle(err.Error())

		return response, errors
	}

	// minimum payment amount validation on first payment of billing
	if bill.Amount == bill.DebtAmount {
		// if less than 50.000, payment amount must be equal to bill amount
		if bill.Amount < int64(models.BillFiftyThousands) && pl.PaymentAmount != bill.Amount {
			errors.SetTitleCode("22", models.DynamicErr(models.ErrExactMatchPaymentAmount, []interface{}{ac.FormatMoney(bill.Amount)}).Error(), "")
			return response, errors
		}

		// if equal or greater than 50.000 AND less than 500.000, minimum payment amount equal to 50.000
		// if equal or greater than 500.000, minimum payment amount equal to 10% of bill amount
		if pl.PaymentAmount < int64(bill.MinimumPayment) {
			errors.SetTitleCode("22", models.DynamicErr(models.ErrMinPaymentAmount, []interface{}{ac.FormatMoney(bill.MinimumPayment)}).Error(), "")
			return response, errors
		}
	}

	// check over payment
	if bill.DebtAmount < pl.PaymentAmount {
		errors.SetTitleCode("22", models.DynamicErr(models.ErrOverPayment, []interface{}{ac.FormatMoney(bill.DebtAmount)}).Error(), "")
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

func (trxUS *transactionsUseCase) mappingCoreInquiry(c echo.Context, respInquiry map[string]interface{}) ([]byte, error) {
	// check reffSwitching variable
	if _, ok := respInquiry["reffSwitching"].(string); !ok {
		logger.Make(c, nil).Debug(models.ErrSetVar)

		return []byte{}, models.ErrSetVar
	}

	// convert respInquiry to json
	respJSON, err := json.Marshal(respInquiry)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return []byte{}, err
	}

	return respJSON, nil
}

func (trxUS *transactionsUseCase) CheckAccountByCIF(c echo.Context, pl interface{}) (models.Account, error) {
	r := reflect.ValueOf(pl)
	cif := r.FieldByName("CIF")

	// Get Account by Account Number
	acc := models.Account{CIF: cif.String()}
	err := trxUS.trxRepo.GetAccountByCIF(c, &acc)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return models.Account{}, models.ErrGetAccByCIF
	}

	return acc, nil
}
