package repository

import (
	"database/sql"
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

func (PSQLTrx *psqlTransactionsRepository) PostBRIPendingTransactions(c echo.Context, trans models.Transaction) error {
	trans.CreatedAt = time.Now()
	err := PSQLTrx.DBpg.Insert(&trans)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
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

func (PSQLTrx *psqlTransactionsRepository) GetAccountByAccountNumber(c echo.Context, trx *models.Transaction) error {
	newAcc := models.Account{}
	err := PSQLTrx.DBpg.Model(&newAcc).Relation("Application").Relation("PersonalInformation").
		Relation("Card").
		Where("account_number = ?", trx.Account.AccountNumber).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return err
	}

	newTrx := models.Transaction{Account: newAcc}

	if err == pg.ErrNoRows {
		return models.ErrAppNumberNotFound
	}

	*trx = newTrx
	return nil
}
