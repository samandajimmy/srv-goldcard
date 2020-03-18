package repository

import (
	"database/sql"
	"gade/srv-goldcard/billings"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"strconv"

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

func (PSQLBill *psqlBillings) GetBillingInquiry(c echo.Context, bill *models.Billing) error {
	// Get last billing published
	newBill := models.Billing{}
	err := PSQLBill.DBpg.Model(&newBill).Relation("Account").
		Where("account_number = ? AND billing.status = ?", bill.Account.AccountNumber, models.BillUnpaid).
		Order("billing_date DESC").Limit(1).Select()

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

func (PSQLBill *psqlBillings) GetMinPaymentParam(c echo.Context) (float64, error) {
	newPrm := models.Parameter{}
	err := PSQLBill.DBpg.Model(&newPrm).
		Where("key = ?", "BILLING_MIN_PAYMENT").Limit(1).Select()

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
		Where("key = ?", "BILLING_INTERVAL_DUE_DATE").Limit(1).Select()

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
		Where("key = ?", "BILLING_PRINT_DATE").Limit(1).Select()

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

func (PSQLBill *psqlBillings) PostPegadaianBillings(c echo.Context, pgdBil models.PegadaianBilling) error {
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
