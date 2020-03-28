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
	trx.CreatedAt = time.Now()
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
		Order("trx_date asc", "id asc").
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

	logger.MakeStructToJSON(trx)

	return trx, nil
}

func (PSQLTrx *psqlTransactionsRepository) GetAccountByAccountNumber(c echo.Context, acc *models.Account) error {
	newAcc := models.Account{}
	err := PSQLTrx.DBpg.Model(&newAcc).Relation("Application").Relation("PersonalInformation").
		Relation("Card").
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
	card.UpdatedAt = time.Now()
	col := []string{"balance", "gold_balance", "stl_balance", "updated_at"}
	_, err := PSQLTrx.DBpg.Model(&card).Column(col...).WherePK().Update()

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (PSQLTrx *psqlTransactionsRepository) PostPayment(c echo.Context, trx models.Transaction, bill models.Billing) error {
	var nilFilters []string

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
			nilFilters, bill.ID, trx.Source, time.Now()),
		// update billings record
		gcdb.NewPipelineStmt(`UPDATE billings set debt_amount = $1, debt_gold = $2, debt_stl = $3,
			updated_at = $4 WHERE id = $5;`,
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
