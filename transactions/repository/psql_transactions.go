package repository

import (
	"database/sql"
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

func (PSQLTrx *psqlTransactionsRepository) GetTrxAccountByBrixKey(c echo.Context, brixkey string) (models.Transaction, error) {
	acc := models.Account{}
	trx := models.Transaction{}
	err := PSQLTrx.DBpg.Model(&acc).Relation("Application").Relation("PersonalInformation").Relation("Card").
		Where("brixkey = ?", brixkey).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return trx, err
	}

	if err == pg.ErrNoRows {
		return trx, models.ErrGetAccByBrixkey
	}

	trx.Account = acc
	trx.AccountId = acc.ID
	return trx, nil
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
