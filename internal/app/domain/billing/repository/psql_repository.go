package repository

import (
	"database/sql"
	"srv-goldcard/internal/app/domain/billing"
	"srv-goldcard/internal/app/model"
	"srv-goldcard/internal/pkg/logger"
	"strconv"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
)

type psqlBillings struct {
	Conn *sql.DB
	DBpg *pg.DB
}

// NewPsqlBillings will create an object that represent the billing.Repository interface
func NewPsqlBillingsRepository(Conn *sql.DB, dbpg *pg.DB) billing.Repository {
	return &psqlBillings{Conn, dbpg}
}

func (PSQLBill *psqlBillings) GetBillingInquiry(c echo.Context, bill *model.Billing) error {
	// Get last billing published
	newBill := *bill
	err := PSQLBill.DBpg.Model(&newBill).Relation("Account").
		Where("account_number = ? AND billing.status = ?", bill.Account.AccountNumber, model.BillUnpaid).
		Order("billing_date DESC").Limit(1).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return err
	}

	if err == pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return model.ErrNoBilling
	}

	*bill = newBill
	return nil
}

func (PSQLBill *psqlBillings) GetMinPaymentParam(c echo.Context) (float64, error) {
	newPrm := model.Parameter{}
	err := PSQLBill.DBpg.Model(&newPrm).
		Where("key = ?", "BILLING_MIN_PAYMENT").Limit(1).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return 0, err
	}

	if err == pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return 0, model.ErrGetParameter
	}

	minPayConst, err := strconv.ParseFloat(newPrm.Value, 64)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return 0, model.ErrParseParameter
	}

	return minPayConst, nil
}

func (PSQLBill *psqlBillings) GetDueDateParam(c echo.Context) (int, error) {
	newPrm := model.Parameter{}
	err := PSQLBill.DBpg.Model(&newPrm).
		Where("key = ?", "BILLING_INTERVAL_DUE_DATE").Limit(1).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return 0, err
	}

	if err == pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return 0, model.ErrGetParameter
	}

	dueDateParam, err := strconv.Atoi(newPrm.Value)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return 0, model.ErrParseParameter
	}

	return dueDateParam, nil
}

func (PSQLBill *psqlBillings) GetBillingPrintDateParam(c echo.Context) (string, error) {
	newPrm := model.Parameter{}
	err := PSQLBill.DBpg.Model(&newPrm).
		Where("key = ?", "BILLING_PRINT_DATE").Limit(1).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return "", err
	}

	if err == pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return "", model.ErrGetParameter
	}

	return newPrm.Value, nil
}

func (PSQLBill *psqlBillings) PostPegadaianBillings(c echo.Context, pgdBil model.PegadaianBilling) error {
	pgdBil.CreatedAt = time.Now()

	query := `INSERT INTO pegadaian_billings (ref_id, billing_date, file_name, file_base64, file_extension, created_at)
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`
	stmt, err := PSQLBill.Conn.Prepare(query)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	_, err = stmt.Exec(pgdBil.RefID, pgdBil.BillingDate, pgdBil.FileName, pgdBil.FileBase64, pgdBil.FileExtension, pgdBil.CreatedAt)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}
