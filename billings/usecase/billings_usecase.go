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
	bill, err := billUS.brRepo.GetBillingsStatement(c, acc)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return models.BillingStatement{}, err
	}

	return bill, nil
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
