package repository

import (
	"database/sql"
	gcdb "gade/srv-goldcard/database"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/transactions"
	"math"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
)

type psqlTransactionsRepository struct {
	Conn *sql.DB
	DBpg *pg.DB
}

// NewPsqlTransactionsRepository will create an object that represent the transactions.Repository interface
func NewPsqlTransactionsRepository(Conn *sql.DB, DBpg *pg.DB) transactions.Repository {
	return &psqlTransactionsRepository{Conn, DBpg}
}

func (PSQLTrx *psqlTransactionsRepository) GetAccountByBrixKey(c echo.Context, brixkey string) (models.Account, error) {
	acc := models.Account{}
	err := PSQLTrx.DBpg.Model(&acc).Relation("Application").Relation("PersonalInformation").
		Relation("Card").Where("brixkey = ?", brixkey).Limit(1).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return acc, err
	}

	if err == pg.ErrNoRows {
		return acc, models.ErrGetAccByBrixkey
	}

	return acc, nil
}

func (PSQLTrx *psqlTransactionsRepository) PostTransactions(c echo.Context, trx models.Transaction) error {
	trx.CreatedAt = models.NowDbpg()
	err := PSQLTrx.DBpg.Insert(&trx)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (PSQLTrx *psqlTransactionsRepository) GetPgTransactionsHistory(c echo.Context, acc models.Account, plListTrx models.PayloadListTrx) (models.ResponseListTrx, error) {
	trx := models.ResponseListTrx{}
	pagination := plListTrx.Pagination
	offset := (pagination.Page - 1) * pagination.Limit

	// count all data
	totalCount, err := PSQLTrx.DBpg.Model(&trx.ListTrx).
		Where("account_id = ?", acc.ID).
		Count()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return trx, err
	}

	if pagination.Limit == 0 {
		pagination.Limit = int64(totalCount)
		trx.IsLastPage = true
	}

	// get the transactions
	err = PSQLTrx.DBpg.Model(&trx.ListTrx).
		Where("account_id = ?", acc.ID).
		Limit(int(pagination.Limit)).Offset(int(offset)).
		Order("trx_date desc", "id asc").
		Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return trx, err
	}

	// flag isLastPage
	totalPage := float64(totalCount) / float64(pagination.Limit)

	if float64(pagination.Page) == math.Ceil(totalPage) {
		trx.IsLastPage = true
	}

	return trx, nil
}

func (PSQLTrx *psqlTransactionsRepository) GetAccountByAccountNumber(c echo.Context, acc *models.Account) error {
	newAcc := models.Account{}
	err := PSQLTrx.DBpg.Model(&newAcc).Relation("Application").Relation("PersonalInformation").
		Relation("Card").Relation("Occupation").Relation("Correspondence").Relation("EmergencyContact").
		Where("account_number = ? AND account.status = ?", acc.AccountNumber, models.AccStatusActive).
		Limit(1).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return err
	}

	if err == pg.ErrNoRows {
		return models.ErrAppNumberNotFound
	}

	*acc = newAcc
	return nil
}

func (PSQLTrx *psqlTransactionsRepository) UpdateCardBalance(c echo.Context, card models.Card) error {
	card.UpdatedAt = models.NowDbpg()
	col := []string{"balance", "gold_balance", "stl_balance", "updated_at"}
	_, err := PSQLTrx.DBpg.Model(&card).Column(col...).WherePK().Update()

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (PSQLTrx *psqlTransactionsRepository) UpdatePayInquiryStatusPaid(c echo.Context, pay models.PaymentInquiry) error {
	pay.UpdatedAt = models.NowDbpg()
	pay.Status = models.BillTrxPaid
	col := []string{"status", "updated_at"}
	_, err := PSQLTrx.DBpg.Model(&pay).Column(col...).WherePK().Update()

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (PSQLTrx *psqlTransactionsRepository) GetPayInquiryByRefTrx(c echo.Context, acc models.Account, refTrx string) (models.PaymentInquiry, error) {
	payment := models.PaymentInquiry{}
	err := PSQLTrx.DBpg.Model(&payment).Relation("Billing").
		Where("payment_inquiry.account_id = ? AND ref_trx = ?", acc.ID, refTrx).
		Where("payment_inquiry.status = ?", models.BillTrxUnpaid).
		Limit(1).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return payment, err
	}

	if err == pg.ErrNoRows {
		return payment, models.ErrGetPaymentTransaction
	}

	return payment, nil
}

func (PSQLTrx *psqlTransactionsRepository) PostPaymentInquiry(c echo.Context, paymentInq models.PaymentInquiry) error {
	paymentInq.InquiryDate = models.NowDbpg()
	paymentInq.CreatedAt = models.NowDbpg()
	err := PSQLTrx.DBpg.Insert(&paymentInq)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (PSQLTrx *psqlTransactionsRepository) PostPayment(c echo.Context, trx models.Transaction, bill models.Billing) error {
	var nilFilters []string
	billPayment := trx.BillingPayments[0]

	stmts := []*gcdb.PipelineStmt{
		// insert payment trx
		gcdb.NewPipelineStmt(`INSERT INTO transactions (account_id, ref_trx_pgdn, ref_trx, nominal,
			gold_nominal, type, status, balance, gold_balance, methods, trx_date, description,
			compare_id, transaction_id, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) RETURNING id;`,
			[]string{"trxID"}, trx.AccountId, trx.RefTrxPgdn, trx.RefTrx, trx.Nominal, trx.GoldBalance,
			trx.Type, trx.Status, trx.Balance, trx.GoldBalance, trx.Methods, trx.TrxDate, trx.Description,
			trx.CompareID, trx.TransactionID, time.Now()),
		// insert billing payment
		gcdb.NewPipelineStmt(`INSERT INTO billing_payments (trx_id, bill_id, source, created_at)
			VALUES ({trxID}, $1, $2, $3) RETURNING id;`,
			nilFilters, bill.ID, billPayment.Source, time.Now()),
		// update billings record
		gcdb.NewPipelineStmt(`UPDATE billings set debt_amount = $1,
			debt_gold = $2, debt_stl = $3, updated_at = $4 WHERE id = $5;`,
			nilFilters, bill.DebtAmount, bill.DebtGold, bill.DebtSTL, time.Now(), bill.ID),
	}

	err := gcdb.WithTransaction(PSQLTrx.Conn, func(tx gcdb.Transaction) error {
		return gcdb.RunPipelineQueryRow(tx, stmts...)
	})

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (PSQLTrx *psqlTransactionsRepository) GetAllActiveAccount(c echo.Context) ([]models.Account, error) {
	var listActiveAccount []models.Account
	err := PSQLTrx.DBpg.Model(&listActiveAccount).Relation("Card").Relation("PersonalInformation").
		Where("account.status = ?", models.AccStatusActive).Select()

	if err != nil || (listActiveAccount == nil) {
		logger.Make(nil, nil).Debug(err)

		return listActiveAccount, err
	}

	return listActiveAccount, nil
}

func (PSQLTrx *psqlTransactionsRepository) GetPaymentInquiryNotificationData(c echo.Context, pi models.PaymentInquiry) (models.PaymentInquiryNotificationData, error) {
	var pind models.PaymentInquiryNotificationData

	nilPind := models.PaymentInquiryNotificationData{}
	query := `SELECT id, core_response->>'reffSwitching' as reff_switching, core_response->>'administrasi' as administration
	FROM payment_inquiries where id = ? limit 1;`

	_, err := PSQLTrx.DBpg.Query(&pind, query, pi.ID)

	if err != nil || (pind == nilPind) {
		logger.Make(nil, nil).Debug(err)

		return nilPind, err
	}

	return pind, nil
}
