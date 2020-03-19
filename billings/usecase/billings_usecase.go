package usecase

import (
	"encoding/base64"
	"fmt"
	"gade/srv-goldcard/billings"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/registrations"
	"gade/srv-goldcard/transactions"
	"strconv"

	"github.com/labstack/echo"
)

type billingsUseCase struct {
	bRepo    billings.Repository
	rrRepo   registrations.RestRepository
	tUseCase transactions.UseCase
}

// billingsUseCase represent billings Use Case
func BillingsUseCase(bRepo billings.Repository, rrRepo registrations.RestRepository, tUseCase transactions.UseCase) billings.UseCase {
	return &billingsUseCase{bRepo, rrRepo, tUseCase}
}

func (billUS *billingsUseCase) GetBillingStatement(c echo.Context, pl models.PayloadAccNumber) (models.BillingStatement, error) {
	// check account
	acc, err := billUS.tUseCase.CheckAccountByAccountNumber(c, pl)

	if err != nil {
		return models.BillingStatement{}, err
	}

	// Get goldcard account billing
	bill := models.Billing{Account: acc}
	err = billUS.bRepo.GetBillingInquiry(c, &bill)

	if err != nil {
		return models.BillingStatement{}, models.ErrNoBilling
	}

	return models.BillingStatement{
		BillingAmount:     bill.Amount,
		BillingPrintDate:  bill.BillingDate.Format("2006-01-02"),
		BillingDueDate:    bill.BillingDueDate.Format("2006-01-02"),
		BillingMinPayment: int64(bill.MinimumPayment),
	}, nil
}

func (billUS *billingsUseCase) PostBRIPegadaianBillings(c echo.Context, pbpb models.PayloadBRIPegadaianBillings) models.ResponseErrors {
	var errors models.ResponseErrors
	var pgdBill models.PegadaianBilling
	err := pgdBill.MappingPegadaianBilling(c, pbpb)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		errors.SetTitle(models.ErrMappingData.Error())

		return errors
	}

	err = billUS.bRepo.PostPegadaianBillings(c, pgdBill)

	if err != nil {
		errors.SetTitle(models.ErrInsertPegadaianBillings.Error())

		return errors
	}

	return errors
}

func (billUS *billingsUseCase) ValidateBase64(c echo.Context, data string) error {
	_, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (billings *billingsUseCase) PaymentInquiry(c echo.Context, ppi models.PayloadPaymentInquiry) models.ResponseErrors {
	var errors models.ResponseErrors

	// Get Account by Account Number
	acc, err := billings.tUseCase.CheckAccountByAccountNumber(c, ppi)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		errors.SetTitle(models.ErrGetAccByAccountNumber.Error())
		return errors
	}

	fmt.Println("---------------------------------------------")
	fmt.Println(acc)
	fmt.Println("---------------------------------------------")

	// get billings by account
	bill := models.Billing{Account: acc}
	err = billings.bRepo.GetBillingInquiry(c, &bill)
	if err != nil {
		logger.Make(c, nil).Debug(err)

		errors.SetTitleCode("11", models.ErrNoBilling.Error(), "")
		return errors
	}

	// check over payment
	if bill.DebtAmount < ppi.PaymentAmount {
		errors.SetTitleCode("22", models.ErrOverPayment.Error(), strconv.FormatInt(bill.DebtAmount, 10))
		return errors
	}

	// check payment less than 10% remaining payment
	if bill.DebtAmount == bill.Amount && ppi.PaymentAmount < bill.DebtAmount/10 {
		errors.SetTitleCode("22", models.ErrMinimumPayment.Error(), strconv.FormatInt(bill.DebtAmount, 10))
		return errors
	}

	return errors
}
