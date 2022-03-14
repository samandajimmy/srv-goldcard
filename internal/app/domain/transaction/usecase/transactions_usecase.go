package usecase

import (
	"encoding/json"
	"reflect"
	"srv-goldcard/internal/app/domain/billing"
	"srv-goldcard/internal/app/domain/registration"
	"srv-goldcard/internal/app/domain/transaction"
	"srv-goldcard/internal/app/model"
	"srv-goldcard/internal/pkg/logger"

	"github.com/labstack/echo"
	"github.com/leekchan/accounting"

	"time"
)

type transactionsUseCase struct {
	trxRepo  transaction.Repository
	billRepo billing.Repository
	trxrRepo transaction.RestRepository
	rRepo    registration.Repository
	rrRepo   registration.RestRepository
}

// TransactionsUseCase represent Transactions Use Case
func TransactionsUseCase(trxRepo transaction.Repository, billRepo billing.Repository,
	trxrRepo transaction.RestRepository, rRepo registration.Repository, rrRepo registration.RestRepository) transaction.UseCase {
	return &transactionsUseCase{trxRepo, billRepo, trxrRepo, rRepo, rrRepo}
}

func (trxUS *transactionsUseCase) PostBRIPendingTransactions(c echo.Context, pl model.PayloadBRIPendingTransactions) model.ResponseErrors {
	var errors model.ResponseErrors
	var notif model.PdsNotification
	acc, err := trxUS.CheckAccountByBrixkey(c, pl)

	if err != nil {
		errors.SetTitle(model.ErrGetAccByBrixkey.Error())
		return errors
	}

	// init account trx
	trx := model.Transaction{AccountId: acc.ID, Account: acc}
	// Get curr STL
	trx.CurrStl, err = trxUS.rrRepo.GetCurrentGoldSTL(c)

	if err != nil {
		errors.SetTitle(model.ErrGetCurrSTL.Error())
		return errors
	}

	// mapping all trx data needed for notification
	err = trx.MappingTrx(pl, model.TypeTrxCredit, true)

	if err != nil {
		errors.SetTitle(model.ErrMappingData.Error())
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

func (trxUS *transactionsUseCase) PostPaymentTransaction(c echo.Context, pl model.PayloadPaymentTransactions) model.ResponseErrors {
	// TODO: do we need push notif or email?
	// TODO: do we need push to core?
	// TODO: we need to add function update billing status when bill is fully paid?
	// TODO: do we need to concurrent to optim response time?
	var errors model.ResponseErrors
	acc, err := trxUS.CheckAccountByBrixkey(c, pl)

	if err != nil {
		errors.SetTitle(model.ErrGetAccByBrixkey.Error())
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
		errors.SetTitle(model.ErrInsertPaymentTransactions.Error())
		return errors
	}

	// post payment transaction to core
	err = trxUS.trxrRepo.PostPaymentTransactionToCore(c, bill)
	if err != nil {
		errors.SetTitleCode("22", model.ErrPostPaymentTransactionToCore.Error(), "")
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

func (trxUS *transactionsUseCase) PostPaymentTrxCore(c echo.Context, pl model.PlPaymentTrxCore) model.ResponseErrors {
	var errors model.ResponseErrors
	var notif model.PdsNotification
	acc, err := trxUS.CheckAccountByAccountNumber(c, pl)

	if err != nil {
		errors.SetTitle(model.ErrGetAccByAccountNumber.Error())
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
	pl.Source = model.SourceCore
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
		errors.SetTitle(model.ErrInsertPaymentTransactions.Error())
		return errors
	}

	// post payment to bri
	err = trxUS.trxrRepo.PostPaymentBRI(c, acc, pl.PaymentAmount)

	if err != nil {
		errors.SetTitle(model.ErrPostPaymentBRI.Error())
		return errors
	}

	// post payment to core when post to bri succeded
	err = trxUS.trxrRepo.PostPaymentCoreNotif(c, acc, pl)

	if err != nil {
		errors.SetTitle(model.ErrPostPaymentCoreNotif.Error())
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

func (trxUS *transactionsUseCase) GetTransactionsHistory(c echo.Context, plListTrx model.PayloadListTrx) (interface{}, model.ResponseErrors) {
	var errors model.ResponseErrors
	var result, pendingArr, postedArr []model.ListTrx
	now := time.Now()
	nowDate := now.Format(model.DateFormat)
	yesterdayDate := now.AddDate(0, 0, -1).Format(model.DateFormat)

	// Get Account by Account Number
	acc, err := trxUS.CheckAccountByAccountNumber(c, plListTrx)

	if err != nil {
		errors.SetTitle(model.ErrGetAccByAccountNumber.Error())

		return model.ResponseListTrx{}, errors
	}

	BRIPending, err := trxUS.trxrRepo.GetBRIPendingTrx(c, acc, yesterdayDate, nowDate)

	if err != nil {
		errors.SetTitle(model.ErrGetHistoryTransactions.Error())
		return model.ResponseListTrx{}, errors
	}

	for _, singleBRIPendigTrx := range BRIPending.TransactionData {
		pendingArr = append(pendingArr, model.ListTrx{
			RefTrx:      "-",
			Nominal:     singleBRIPendigTrx.BillAmount,
			TrxDate:     time.Unix(singleBRIPendigTrx.TransactionDate/1000, 0).Format(model.DateTimeFormatZone),
			Description: singleBRIPendigTrx.Description})
	}

	pendingArr = model.ReverseArray(pendingArr)

	BRIPosted, err := trxUS.trxrRepo.GetBRIPostedTrx(c, acc.BrixKey)

	if err != nil {
		errors.SetTitle(model.ErrGetHistoryTransactions.Error())
		return model.ResponseListTrx{}, errors
	}

	for _, singleBRIPostedTrx := range BRIPosted.ListOfTransactions {
		postedArr = append(postedArr, model.ListTrx{
			RefTrx:      singleBRIPostedTrx.TrxReff,
			Nominal:     singleBRIPostedTrx.TrxAmount,
			TrxDate:     time.Unix(singleBRIPostedTrx.EffectiveDate/1000, 0).Format(model.DateTimeFormatZone),
			Description: singleBRIPostedTrx.TrxDesc})
	}

	postedArr = model.ReverseArray(postedArr)

	result = append(result, pendingArr...)
	result = append(result, postedArr...)

	return result, errors
}

func (trxUS *transactionsUseCase) CheckAccountByAccountNumber(c echo.Context, pl interface{}) (model.Account, error) {
	r := reflect.ValueOf(pl)
	accNumber := r.FieldByName("AccountNumber")

	// Get Account by Account Number
	acc := model.Account{AccountNumber: accNumber.String()}
	err := trxUS.trxRepo.GetAccountByAccountNumber(c, &acc)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return model.Account{}, model.ErrGetAccByAccountNumber
	}

	return acc, nil
}

func (trxUS *transactionsUseCase) CheckAccountByBrixkey(c echo.Context, pl interface{}) (model.Account, error) {
	r := reflect.ValueOf(pl)
	BrixKey := r.FieldByName("BrixKey")

	// Get trx Account by BrixKey
	acc, err := trxUS.trxRepo.GetAccountByBrixKey(c, BrixKey.String())

	if err != nil {
		return model.Account{}, model.ErrGetAccByBrixkey
	}

	return acc, nil
}

func (trxUS *transactionsUseCase) GetCardBalance(c echo.Context, pl model.PayloadAccNumber) (model.BRICardBalance, error) {
	var briCardBal model.BRICardBalance
	// check account number
	acc, err := trxUS.CheckAccountByAccountNumber(c, pl)

	if err != nil {
		return briCardBal, err
	}

	// update and get card balance by account
	card, err := trxUS.UpdateAndGetCardBalance(c, acc)

	if err != nil {
		return briCardBal, model.ErrGetCardBalance
	}

	// mapping + return bri card balance
	return model.BRICardBalance{
		AvailableCredit: card.Balance,
		CreditLimit:     float64(card.CardLimit),
	}, nil
}

func (trxUS *transactionsUseCase) UpdateAndGetCardBalance(c echo.Context, acc model.Account) (model.Card, error) {
	// define channel buffer
	errPromise := make(chan error)
	briCardBal := make(chan model.BRICardBalance)
	currStl := make(chan int64)

	go func() {
		// get balance from bank
		cardBal, err := trxUS.trxrRepo.GetBRICardInformation(c, acc)

		if err != nil {
			briCardBal <- model.BRICardBalance{}
			errPromise <- model.ErrGetCardBalance
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
			errPromise <- model.ErrGetCurrSTL
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
		return model.Card{}, err
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
		return model.Card{}, model.ErrUpdateCardBalance
	}

	return acc.Card, nil
}

func (trxUS *transactionsUseCase) PaymentInquiry(c echo.Context, pl model.PlPaymentInquiry) (map[string]interface{}, model.ResponseErrors) {
	var errors model.ResponseErrors
	var ac = accounting.Accounting{Symbol: "Rp ", Thousand: "."}
	response := map[string]interface{}{}
	// Get Account by Account Number
	acc, err := trxUS.CheckAccountByAccountNumber(c, pl)

	if err != nil {
		errors.SetTitle(model.ErrGetAccByAccountNumber.Error())
		return response, errors
	}

	// init current billing
	bill := model.Billing{Account: acc}
	// get last published billing to map payment transaction into billings
	err = trxUS.billRepo.GetBillingInquiry(c, &bill)

	if err != nil {
		errors.SetTitleCode("11", model.ErrNoBilling.Error(), "")
		return response, errors
	}

	// payment inquiry to core
	respInquiry, err := trxUS.trxrRepo.CorePaymentInquiry(c, pl, acc)

	if err != nil {
		errors.SetTitleCode("11", model.ErrNoBilling.Error(), "")
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
		if bill.Amount < int64(model.BillFiftyThousands) && pl.PaymentAmount != bill.Amount {
			errors.SetTitleCode("22", model.DynamicErr(model.ErrExactMatchPaymentAmount, []interface{}{ac.FormatMoney(bill.Amount)}).Error(), "")
			return response, errors
		}

		// if equal or greater than 50.000 AND less than 500.000, minimum payment amount equal to 50.000
		// if equal or greater than 500.000, minimum payment amount equal to 10% of bill amount
		if pl.PaymentAmount < int64(bill.MinimumPayment) {
			errors.SetTitleCode("22", model.DynamicErr(model.ErrMinPaymentAmount, []interface{}{ac.FormatMoney(bill.MinimumPayment)}).Error(), "")
			return response, errors
		}
	}

	// check over payment
	if bill.DebtAmount < pl.PaymentAmount {
		errors.SetTitleCode("22", model.DynamicErr(model.ErrOverPayment, []interface{}{ac.FormatMoney(bill.DebtAmount)}).Error(), "")
		return response, errors
	}

	// prepare the payment inquiry data
	paymentInq := model.PaymentInquiry{
		AccountId:    acc.ID,
		BillingId:    bill.ID,
		RefTrx:       refTrx,
		Nominal:      pl.PaymentAmount,
		CoreResponse: respJSON,
	}

	// insert payment inquiry
	err = trxUS.trxRepo.PostPaymentInquiry(c, paymentInq)

	if err != nil {
		errors.SetTitle(model.ErrInsertPaymentTransactions.Error())
		return response, errors
	}

	return respInquiry, errors
}

func (trxUS *transactionsUseCase) prepareTrxAndBill(c echo.Context, acc model.Account, pl interface{}) (model.Transaction,
	model.Billing, model.ResponseErrors) {
	var errors = model.ResponseErrors{Code: "00"}
	var err error
	// init account trx
	trx := model.Transaction{AccountId: acc.ID, Account: acc}
	// init current billing
	bill := model.Billing{Account: acc}
	// Get curr STL
	trx.CurrStl, err = trxUS.rrRepo.GetCurrentGoldSTL(c)

	if err != nil {
		errors.SetTitle(model.ErrGetCurrSTL.Error())
		return model.Transaction{}, model.Billing{}, errors
	}

	// populate payment transaction for insert
	err = trx.MappingTrx(pl, model.TypeTrxDebit, false)

	if err != nil {
		errors.SetTitle(model.ErrMappingData.Error())
		return model.Transaction{}, model.Billing{}, errors
	}

	// get last published billing to map payment transaction into billings
	err = trxUS.billRepo.GetBillingInquiry(c, &bill)

	if err != nil {
		errors.SetTitleCode("11", model.ErrNoBilling.Error(), "")
		return model.Transaction{}, model.Billing{}, errors
	}

	return trx, bill, errors
}

func (trxUS *transactionsUseCase) payTheBill(c echo.Context, bill *model.Billing, trx model.Transaction) error {
	// prepare billing debt
	// get debt amount
	bill.DebtAmount = bill.DebtAmount - trx.Nominal

	// change bill status to paid
	if bill.DebtAmount <= 0 && bill.Status != model.BillTrxPaid {
		bill.Status = model.BillTrxPaid
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
		logger.Make(c, nil).Debug(model.ErrSetVar)

		return []byte{}, model.ErrSetVar
	}

	// convert respInquiry to json
	respJSON, err := json.Marshal(respInquiry)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return []byte{}, err
	}

	return respJSON, nil
}

func (trxUS *transactionsUseCase) CheckAccountByCIF(c echo.Context, pl model.PayloadCIF) (model.Account, error) {
	// Get Account by Account Number
	acc := model.Account{CIF: pl.CIF}
	err := trxUS.trxRepo.GetAccountByCIF(c, &acc)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return model.Account{}, model.ErrGetAccByCIF
	}

	return acc, nil
}
