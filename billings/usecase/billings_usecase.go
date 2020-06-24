package usecase

import (
	"gade/srv-goldcard/billings"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/registrations"
	"gade/srv-goldcard/transactions"

	"github.com/labstack/echo"
)

type billingsUseCase struct {
	bRepo    billings.Repository
	brRepo   billings.RestRepository
	rrRepo   registrations.RestRepository
	tUseCase transactions.UseCase
}

// billingsUseCase represent billings Use Case
func BillingsUseCase(bRepo billings.Repository, brRepo billings.RestRepository, rrRepo registrations.RestRepository, tUseCase transactions.UseCase) billings.UseCase {
	return &billingsUseCase{bRepo, brRepo, rrRepo, tUseCase}
}

func (billUS *billingsUseCase) GetBillingStatement(c echo.Context, pl models.PayloadAccNumber) (models.BillingStatement, error) {
	// check account
	acc, err := billUS.tUseCase.CheckAccountByAccountNumber(c, pl)

	if err != nil {
		return models.BillingStatement{}, err
	}

	// Get goldcard account billing
	bill, respCode := billUS.brRepo.GetBillingsStatement(c, acc)

	logger.Make(c, nil).Debug(bill)

	if respCode == "5X" {
		return models.BillingStatement{}, models.ErrNoBilling
	}

	response := bill["listOfStatements"].(map[string]interface{})["statementHeader"].(map[string]interface{})

	return models.BillingStatement{
		BillingAmount:     int64(response["totalPayment"].(float64)),
		BillingPrintDate:  response["statementDate"].(string),
		BillingDueDate:    response["paymentDueDate"].(string),
		BillingMinPayment: int64(response["totalPayment"].(float64)) * 10 / 100,
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
