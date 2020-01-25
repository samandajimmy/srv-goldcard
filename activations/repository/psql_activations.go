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

type psqlActivationsRepository struct {
	Conn *sql.DB
	DBpg *pg.DB
}

// NewPsqlActivations will create an object that represent the activations.Repository interface
func NewPsqlActivations(Conn *sql.DB, dbpg *pg.DB) activations.Repository {
	return &psqlActivationsRepository{Conn, dbpg}
}

func (act *psqlActivationsRepository) PostActivations(c echo.Context, acc models.Account) error {
	var nilFilters []string
	app := acc.Application
	card := acc.Card

	stmts := []*gcdb.PipelineStmt{
		gcdb.NewPipelineStmt(`UPDATE accounts SET status = $1, updated_at = $2
			WHERE id = $3`,
			nilFilters, acc.Status, time.Now(), acc.ID),
		gcdb.NewPipelineStmt(`UPDATE applications SET status = $1, update_at = $2
			WHERE id = $3`,
			nilFilters, app.Status, time.Now(), acc.ApplicationID),
		gcdb.NewPipelineStmt(`UPDATE cards SET status = $1, card_number = $2, valid_until = $3, update_at = $4
			WHERE id = $5`,
			nilFilters, card.Status, card.CardNumber, card.ValidUntil, time.Now(), acc.CardID),
	}

	err := gcdb.WithTransaction(act.Conn, func(tx gcdb.Transaction) error {
		return gcdb.RunPipelineQueryRow(tx, stmts...)
	})

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}
