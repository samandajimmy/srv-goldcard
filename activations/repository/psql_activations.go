package repository

import (
	"database/sql"
	"gade/srv-goldcard/activations"
	gcdb "gade/srv-goldcard/database"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
)

type psqlActivations struct {
	Conn *sql.DB
	DBpg *pg.DB
}

// NewPsqlActivations will create an object that represent the activations.Repository interface
func NewPsqlActivations(Conn *sql.DB, dbpg *pg.DB) activations.Repository {
	return &psqlActivations{Conn, dbpg}
}

func (pa *psqlActivations) PostActivations(c echo.Context, acc models.Account) error {
	var nilFilters []string
	app := acc.Application
	card := acc.Card
	stmts := []*gcdb.PipelineStmt{
		gcdb.NewPipelineStmt(`UPDATE accounts SET status = $1, updated_at = $2, account_number = $3
			WHERE id = $4`,
			nilFilters, acc.Status, time.Now(), acc.AccountNumber, acc.ID),
		gcdb.NewPipelineStmt(`UPDATE applications SET status = $1, updated_at = $2
			WHERE id = $3`,
			nilFilters, app.Status, time.Now(), acc.ApplicationID),
		gcdb.NewPipelineStmt(`UPDATE cards SET status = $1, card_number = $2, valid_until = $3, updated_at = $4
			WHERE id = $5`,
			nilFilters, card.Status, card.CardNumber, card.ValidUntil, time.Now(), acc.CardID),
	}

	err := gcdb.WithTransaction(pa.Conn, func(tx gcdb.Transaction) error {
		return gcdb.RunPipelineQueryRow(tx, stmts...)
	})
	if err != nil {
		logger.Make(c, nil).Debug(err)
		return err
	}
	return nil
}

func (pa *psqlActivations) UpdateGoldLimit(c echo.Context, card models.Card) error {
	card.UpdatedAt = models.NowDbpg()
	col := []string{"gold_limit", "gold_balance", "stl_limit", "stl_balance", "updated_at"}
	_, err := pa.DBpg.Model(&card).Column(col...).WherePK().Update()

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (act *psqlActivations) GetAccountByAppNumber(c echo.Context, acc *models.Account) error {
	newAcc := models.Account{}
	docs := []models.Document{}
	err := act.DBpg.Model(&newAcc).Relation("Application").Relation("Card").Relation("PersonalInformation").
		Where("application_number = ?", acc.Application.ApplicationNumber).
		Where("application.status = ?", models.AppStatusSent).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return err
	}

	if err == pg.ErrNoRows {
		return models.ErrAppNumberNotFound
	}

	err = act.DBpg.Model(&docs).Where("application_id = ?", newAcc.ApplicationID).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return err
	}

	newAcc.Application.Documents = docs

	*acc = newAcc
	return nil
}

func (act *psqlActivations) GetStoredGoldPrice(c echo.Context) (int64, error) {
	var goldPrice models.GoldPrice

	query := `select id, price, valid_date, created_at
		from gold_prices order by created_at desc limit 1;`

	_, err := act.DBpg.Query(&goldPrice, query)

	if err != nil || (goldPrice == models.GoldPrice{}) {
		logger.Make(nil, nil).Debug(err)

		return 0, err
	}

	return int64(goldPrice.Price), nil
}
