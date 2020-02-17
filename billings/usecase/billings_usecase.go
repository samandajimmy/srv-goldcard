package usecase

import (
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
	var billStmt models.BillingStatement
	acc, err := billUS.tUseCase.CheckAccountByAccountNumber(c, pl)

	if err != nil {
		return billStmt, err
	}

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return billStmt, err
	}

	// Get goldcard account billing
	bill := models.Billing{Account: acc}
	err = billUS.bRepo.GetBilling(c, &bill)

	if err != nil {
		return billStmt, err
	}

	// Get minimum payment parameters
	minPayParam, err := billUS.bRepo.GetMinPaymentParam(c)

	if err != nil {
		return billStmt, err
	}

	// Get billing due date parameters
	dueDateParam, err := billUS.bRepo.GetDueDateParam(c)

	if err != nil {
		return billStmt, err
	}

	err = bill.MapBillingStatementResponse(c, dueDateParam, minPayParam, &billStmt)

	if err != nil {
		return billStmt, err
	}

	return billStmt, nil
}

// ProductRequirements represent to get all product requirements
func (billings *billingsUseCase) ProductRequirements(c echo.Context) (models.Requirements, error) {
	return models.RequirementsValue, nil
}
