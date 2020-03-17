package usecase

import (
	"encoding/base64"
	"gade/srv-goldcard/billings"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/transactions"

	"github.com/labstack/echo"
)

type billingsUseCase struct {
	bRepo    billings.Repository
	tUseCase transactions.UseCase
}

// billingsUseCase represent billings Use Case
func BillingsUseCase(bRepo billings.Repository, tUseCase transactions.UseCase) billings.UseCase {
	return &billingsUseCase{bRepo, tUseCase}
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

	// validate base64 file payload
	if err := billUS.ValidateBase64(c, pbpb.FileBase64); err != nil {
		errors.SetTitle(models.ErrValidateBase64.Error())

		return errors
	}

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
