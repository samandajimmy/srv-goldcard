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

func (trx *psqlTransactionsRepository) PostBRIPendingTransactions(c echo.Context, trans models.Transaction) error {
	trans.CreatedAt = time.Now()
	err := trx.DBpg.Insert(&trans)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (trx *psqlTransactionsRepository) GetAccountByBrixKey(c echo.Context, trans *models.Transaction) error {
	newAcc := models.Account{}
	err := trx.DBpg.Model(&newAcc).Relation("Application").Relation("PersonalInformation").Relation("Card").
		Where("brixkey = ?", trans.Account.BrixKey).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return err
	}

	newTrans := models.Transaction{Account: newAcc}

	if err == pg.ErrNoRows {
		return models.ErrGetAccByBrixkey
	}

	*trans = newTrans
	return nil
}
