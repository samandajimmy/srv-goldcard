package usecase

import (
	"math"
	"reflect"
	"time"

	"gade/srv-goldcard/billings"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

type billingsUseCase struct {
	bRepo billings.Repository
}

// billingsUseCase represent billings Use Case
func BillingsUseCase(bRepo billings.Repository) billings.UseCase {
	return &billingsUseCase{bRepo}
}

func (billUS *billingsUseCase) GetBillingStatement(c echo.Context, pl models.PayloadAccNumber) (models.BillingStatement, error) {
	var billStmt models.BillingStatement
	bill, err := billUS.checkAccountByAccountNumber(c, pl)

	if err != nil {
		return billStmt, err
	}

	err = bill.MappingAccountNumberToBilling(c, pl)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return billStmt, err
	}

	// Get goldcard account billing
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

	time := time.Now()
	yyyyMM := time.Format("2006-01")

	billStmt.BillingAmount = bill.Amount
	billStmt.BillingDueDate = yyyyMM + "-" + dueDateParam
	billStmt.BillingMinPayment = int64(math.Ceil(float64(bill.Amount) * minPayParam))
	billStmt.BillingPrintDate = bill.BillingDate.Format("2006-01-02")

	return billStmt, nil
}

func (billUS *billingsUseCase) checkAccountByAccountNumber(c echo.Context, pl interface{}) (models.Billing, error) {
	r := reflect.ValueOf(pl)
	accNumber := r.FieldByName("AccountNumber")

	// Get Account by Account Number
	bill := models.Billing{Account: models.Account{AccountNumber: accNumber.String()}}
	err := billUS.bRepo.GetAccountByAccountNumber(c, &bill)

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return bill, models.ErrGetAccByAccountNumber
	}

	return bill, nil
}

// ProductRequirements represent to get all product requirements
func (billings *billingsUseCase) ProductRequirements(c echo.Context) (models.Requirements, error) {
	return models.RequirementsValue, nil
}
