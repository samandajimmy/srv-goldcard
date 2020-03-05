package repository

import (
	"database/sql"
	gcdb "gade/srv-goldcard/database"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/transactions"
	"math"
	"strconv"
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

func (PSQLTrx *psqlTransactionsRepository) GetAllTransactionsHistory(c echo.Context, pt models.PayloadHistoryTransactions) (models.ResponseHistoryTransactions, error) {
	trx := models.ResponseHistoryTransactions{}

	_, err := PSQLTrx.DBpg.Query(&trx.ListHistoryTransactions, `SELECT t.ref_trx, t.nominal, t.trx_date, t.description FROM transactions t 
		LEFT JOIN accounts a ON a.id = t.account_id WHERE a.account_number = ? ORDER BY t.created_at`,
		pt.AccountNumber)

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return trx, err
	}

	trx.IsLastPage = true
	return trx, nil

}

func (PSQLTrx *psqlTransactionsRepository) GetAccountByBrixKey(c echo.Context, trx *models.Transaction) error {
	newAcc := models.Account{}
	err := PSQLTrx.DBpg.Model(&newAcc).Relation("Application").Relation("PersonalInformation").Relation("Card").
		Where("brixkey = ?", trx.Account.BrixKey).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return err
	}

	newTrx := models.Transaction{Account: newAcc}

	if err == pg.ErrNoRows {
		return models.ErrGetAccByBrixkey
	}

	*trx = newTrx
	return nil
}

func (PSQLTrx *psqlTransactionsRepository) PostTransactions(c echo.Context, trx models.Transaction) error {
	trx.CreatedAt = time.Now()
	var nilFilters []string
	stmts := []*gcdb.PipelineStmt{
		gcdb.NewPipelineStmt(`INSERT INTO transactions (account_id, ref_trx_pgdn, transaction_id, nominal, gold_nominal,
			type, status, balance, gold_balance, description, compare_id, created_at, trx_date)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9 ,$10, $11, $12, $13) RETURNING account_id;`,
			nilFilters, trx.AccountId, trx.RefTrxPgdn, trx.TransactionID, trx.Nominal, trx.GoldNominal, trx.Type, trx.Status, trx.Balance,
			trx.GoldBalance, trx.Description, trx.CompareID, time.Now(), trx.TrxDate),

		gcdb.NewPipelineStmt(`UPDATE cards c 
			set balance = $1, gold_balance = $2, updated_at = $3
			FROM accounts a 
			WHERE a.card_id = c.id
			AND a.id = `+strconv.Itoa(int(trx.AccountId))+` RETURNING c.id;`,
			nilFilters, trx.Balance, trx.GoldBalance, time.Now()),
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

func (PSQLTrx *psqlTransactionsRepository) GetPgTransactionsHistory(c echo.Context, pt models.PayloadHistoryTransactions) (models.ResponseHistoryTransactions, error) {
	trx := models.ResponseHistoryTransactions{}
	offset := (pt.Pagination.Page - 1) * pt.Pagination.Limit

	_, err := PSQLTrx.DBpg.Query(&trx.ListHistoryTransactions, `SELECT t.ref_trx, t.nominal, t.trx_date, t.description FROM transactions t 
		LEFT JOIN accounts a ON a.id = t.account_id  WHERE a.account_number = ? ORDER BY t.created_at LIMIT ? OFFSET ?`,
		pt.AccountNumber, pt.Pagination.Limit, offset)

	if err != nil && err != pg.ErrNoRows {
		return trx, err
	}

	// get total data transactions
	total, err := PSQLTrx.getTotalTransactions(c, pt)

	if err != nil && err != pg.ErrNoRows {
		return trx, err
	}

	// flag isLastPage
	totalPage := total / float64(pt.Pagination.Limit)
	if float64(pt.Pagination.Page) == math.Ceil(totalPage) {
		trx.IsLastPage = true

		return trx, nil
	}

	return trx, nil
}

func (PSQLTrx *psqlTransactionsRepository) getTotalTransactions(c echo.Context, pt models.PayloadHistoryTransactions) (float64, error) {
	var count float64

	_, err := PSQLTrx.DBpg.Query(&count, `SELECT count(t.id) FROM transactions t 
		LEFT JOIN accounts a ON a.id = t.account_id WHERE a.account_number = ?`,
		pt.AccountNumber)

	if err != nil && err != pg.ErrNoRows {
		return count, err
	}

	return count, err
}

func (PSQLTrx *psqlTransactionsRepository) GetAccountByAccountNumber(c echo.Context, acc *models.Account) error {
	newAcc := models.Account{}
	err := PSQLTrx.DBpg.Model(&newAcc).Relation("Application").Relation("PersonalInformation").
		Relation("Card").
		Where("account_number = ?", acc.AccountNumber).Select()

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
