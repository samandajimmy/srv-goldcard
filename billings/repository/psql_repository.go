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
	month, err := PSQLBill.calculateBillingMonth(c)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	// Get last billing published
	newBill := models.Billing{}
	err = PSQLBill.DBpg.Model(&newBill).Relation("Account").
		Where("account.account_number = ?", bill.Account.AccountNumber).
		Where("EXTRACT(MONTH FROM billing_date) = ?", month).Select()

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

func (PSQLBill *psqlBillings) calculateBillingMonth(c echo.Context) (time.Month, error) {
	// Populate billing date to filter which billing to get
	timeNow := time.Now()

	dd, err := PSQLBill.GetBillingPrintDateParam(c)

	if err != nil {
		return 0, err
	}

	day, err := strconv.Atoi(dd)

	if err != nil {
		return 0, err
	}

	// if day < 2 then it shown previous month billing (only for the 1st)
	if timeNow.Day() < day {
		month := timeNow.AddDate(0, -2, 0).Month()

		return month, nil
	}

	month := timeNow.AddDate(0, -1, 0).Month()

	return month, nil
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

func (PSQLBill *psqlBillings) GetDueDateParam(c echo.Context) (int, error) {
	newPrm := models.Parameter{}
	err := PSQLBill.DBpg.Model(&newPrm).
		Where("key = ?", "BILLING_INTERVAL_DUE_DATE").Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return 0, err
	}

	if err == pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return 0, models.ErrGetParameter
	}

	dueDateParam, err := strconv.Atoi(newPrm.Value)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return 0, models.ErrParseParameter
	}

	return dueDateParam, nil
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
