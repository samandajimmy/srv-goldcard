package repository

import (
	"database/sql"
	"gade/srv-goldcard/cards"
	gcdb "gade/srv-goldcard/database"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"

	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
)

type psqlCardsRepository struct {
	Conn *sql.DB
	DBpg *pg.DB
}

// NewPsqlCardsRepository will create an object that represent the cards.Repository interface
func NewPsqlCardsRepository(Conn *sql.DB, dbpg *pg.DB) cards.Repository {
	return &psqlCardsRepository{Conn, dbpg}
}

func (PSQLCard *psqlCardsRepository) UpdateCardStatus(c echo.Context, card models.Card, cs models.CardStatuses) error {
	time := models.NowDbpg()

	var nilFilters []string
	stmts := []*gcdb.PipelineStmt{
		gcdb.NewPipelineStmt(`UPDATE cards SET status = $1, updated_at = $2 WHERE id = $3`,
			nilFilters, card.Status, time, card.ID),
		gcdb.NewPipelineStmt(`INSERT INTO card_statuses (card_id, reason, reason_code, blocked_date, is_reactivated, created_at) VALUES ($1, $2, $3, $4, $5, $6)`,
			nilFilters, card.ID, cs.Reason, cs.ReasonCode, cs.BlockedDate, cs.IsReactivated, time),
	}

	err := gcdb.WithTransaction(PSQLCard.Conn, func(tx gcdb.Transaction) error {
		return gcdb.RunPipelineQueryRow(tx, stmts...)
	})
	if err != nil {
		logger.Make(c, nil).Debug(err)
		return err
	}
	return nil
}
