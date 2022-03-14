package usecase

import (
	"srv-goldcard/internal/app/domain/billing"
	"srv-goldcard/internal/app/domain/registration"
	"srv-goldcard/internal/app/domain/transaction"
	"srv-goldcard/internal/app/model"
	"srv-goldcard/internal/pkg/logger"

	"github.com/labstack/echo"
)

type billingsUseCase struct {
	bRepo    billing.Repository
	brRepo   billing.RestRepository
	rrRepo   registration.RestRepository
	tUseCase transaction.UseCase
}

// billingsUseCase represent billings Use Case
func BillingsUseCase(bRepo billing.Repository, brRepo billing.RestRepository, rrRepo registration.RestRepository, tUseCase transaction.UseCase) billing.UseCase {
	return &billingsUseCase{bRepo, brRepo, rrRepo, tUseCase}
}

func (billUS *billingsUseCase) GetBillingStatement(c echo.Context, pl model.PayloadAccNumber) (model.BillingStatement, error) {
	// check account
	acc, err := billUS.tUseCase.CheckAccountByAccountNumber(c, pl)

	if err != nil {
		return model.BillingStatement{}, err
	}

	// Get goldcard account billing
	bill, err := billUS.brRepo.GetBillingsStatement(c, acc)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return model.BillingStatement{}, err
	}

	return bill, nil
}

func (billUS *billingsUseCase) PostBRIPegadaianBillings(c echo.Context, pbpb model.PayloadBRIPegadaianBillings) model.ResponseErrors {
	var errors model.ResponseErrors
	var pgdBill model.PegadaianBilling
	err := pgdBill.MappingPegadaianBilling(c, pbpb)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		errors.SetTitle(model.ErrMappingData.Error())

		return errors
	}

	err = billUS.bRepo.PostPegadaianBillings(c, pgdBill)

	if err != nil {
		errors.SetTitle(model.ErrInsertPegadaianBillings.Error())

		return errors
	}

	return errors
}
