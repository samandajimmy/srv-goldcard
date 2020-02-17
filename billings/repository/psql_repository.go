package repository

import (
	"database/sql"
	"gade/srv-goldcard/billings"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"strconv"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
)

type psqlBillings struct {
	Conn *sql.DB
	DBpg *pg.DB
}

// NewPsqlBillings will create an object that represent the billings.Repository interface
func NewPsqlBillingsRepository(Conn *sql.DB, dbpg *pg.DB) billings.Repository {
	return &psqlBillings{Conn, dbpg}
}

func (PSQLBill *psqlBillings) GetBilling(c echo.Context, bill *models.Billing) error {
	// Populate billing date to filter which billing to get
	time := time.Now()

	// if day < 2 then it shown previous month billing
	if time.Day() < 2 {
		time = time.AddDate(0, -1, 0)
	}

	yyyy := time.Format("2006")
	MM := time.Format("01")
	dd, err := PSQLBill.GetBillingPrintDateParam(c)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return models.ErrGetParameter
	}

	HHmmss := "00:00:00"
	billingDate := yyyy + "-" + MM + "-" + dd + " " + HHmmss

	// Get last billing published
	newBill := models.Billing{}
	err = PSQLBill.DBpg.Model(&newBill).Relation("Account").
		Where("account.account_number = ?", bill.Account.AccountNumber).
		Where("billing_date >= ?", billingDate).
		Order("created_at DESC").
		Limit(1).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return err
	}

	if err == pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return models.ErrNoBilling
	}

	*bill = newBill
	return nil
}

func (PSQLBill *psqlBillings) GetAccountByAccountNumber(c echo.Context, bill *models.Billing) error {
	newAcc := models.Account{}
	err := PSQLBill.DBpg.Model(&newAcc).
		Where("account_number = ?", bill.Account.AccountNumber).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return err
	}

	newBill := models.Billing{Account: newAcc}

	if err == pg.ErrNoRows {
		return models.ErrAppNumberNotFound
	}

	*bill = newBill
	return nil
}

func (PSQLBill *psqlBillings) GetMinPaymentParam(c echo.Context) (float64, error) {
	newPrm := models.Parameter{}
	err := PSQLBill.DBpg.Model(&newPrm).
		Where("key = ?", "BILLING_MIN_PAYMENT").Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return 0, err
	}

	if err == pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return 0, models.ErrGetParameter
	}

	minPayConst, err := strconv.ParseFloat(newPrm.Value, 64)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return 0, models.ErrParseParameter
	}

	return minPayConst, nil
}

func (PSQLBill *psqlBillings) GetDueDateParam(c echo.Context) (string, error) {
	newPrm := models.Parameter{}
	err := PSQLBill.DBpg.Model(&newPrm).
		Where("key = ?", "BILLING_DUE_DATE").Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return "", err
	}

	if err == pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return "", models.ErrGetParameter
	}

	return newPrm.Value, nil
}

func (PSQLBill *psqlBillings) GetBillingPrintDateParam(c echo.Context) (string, error) {
	newPrm := models.Parameter{}
	err := PSQLBill.DBpg.Model(&newPrm).
		Where("key = ?", "BILLING_PRINT_DATE").Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return "", err
	}

	if err == pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return "", models.ErrGetParameter
	}

	return newPrm.Value, nil
}
