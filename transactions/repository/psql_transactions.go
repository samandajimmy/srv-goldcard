package repository

import (
	"database/sql"
	gcdb "gade/srv-goldcard/database"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/transactions"
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

func (PSQLTrx *psqlTransactionsRepository) PostBRIPendingTransactions(c echo.Context, trx models.Transaction) error {
	err := PSQLTrx.PostTransactionsANDCardBalance(c, trx)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (PSQLTrx *psqlTransactionsRepository) GetTransactionsHistory(c echo.Context, pt models.PayloadHistoryTransactions) ([]models.ResponseHistoryTransactions, error) {
	trx := []models.ResponseHistoryTransactions{}
	_, err := PSQLTrx.DBpg.Query(&trx, `SELECT t.nominal, t.trx_date, t.status, t.description FROM transactions t 
		LEFT JOIN accounts a ON a.id = t.account_id WHERE a.account_number = ? LIMIT ? OFFSET ?`,
		pt.AccountNumber, pt.Pagination.Limit, pt.Pagination.Offset)

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return trx, err
	}

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

func (PSQLTrx *psqlTransactionsRepository) PostTransactionsANDCardBalance(c echo.Context, trx models.Transaction) error {
	trx.CreatedAt = time.Now()
	var nilFilters []string
	stmts := []*gcdb.PipelineStmt{
		gcdb.NewPipelineStmt(`INSERT INTO transactions (account_id, ref_trx_pgdn, ref_trx, nominal, gold_nominal,
			type, status, balance, gold_balance, description, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9 ,$10, $11) RETURNING account_id;`,
			[]string{"accID"}, trx.AccountId, trx.RefTrxPgdn, trx.RefTrx, trx.Nominal, trx.GoldNominal, trx.Type, trx.Status, trx.Balance,
			trx.GoldBalance, trx.Description, time.Now()),

		gcdb.NewPipelineStmt(`UPDATE cards c 
			set balance = $1, gold_balance = $2, updated_at = $3
			FROM accounts a 
			WHERE a.card_id = c.id
			AND a.id = {accID} RETURNING c.id;`,
			nilFilters, trx.Balance, trx.GoldBalance, time.Now()),
	}

	err := gcdb.WithTransaction(PSQLTrx.Conn, func(tx gcdb.Transaction) error {
		return gcdb.RunPipelineQueryRow(tx, stmts...)
	})

	if err != nil {
		return err
	}

	return nil
}
