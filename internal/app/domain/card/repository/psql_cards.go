package repository

import (
	"database/sql"
	"srv-goldcard/internal/app/domain/card"
	"srv-goldcard/internal/app/model"
	gcdb "srv-goldcard/internal/pkg/database"
	"srv-goldcard/internal/pkg/logger"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
)

type psqlCardsRepository struct {
	Conn *sql.DB
	DBpg *pg.DB
}

// NewPsqlCardsRepository will create an object that represent the card.Repository interface
func NewPsqlCardsRepository(Conn *sql.DB, dbpg *pg.DB) card.Repository {
	return &psqlCardsRepository{Conn, dbpg}
}

func (PSQLCard *psqlCardsRepository) UpdateCardStatus(c echo.Context, card model.Card, cs model.CardStatuses) error {
	var nilFilters []string

	stmts := []*gcdb.PipelineStmt{
		gcdb.NewPipelineStmt(`UPDATE cards SET status = $1, updated_at = $2 WHERE id = $3`,
			nilFilters, card.Status, time.Now(), card.ID),
		gcdb.NewPipelineStmt(`INSERT INTO card_statuses (card_id, reason, reason_code, blocked_date, is_reactivated, created_at) VALUES ($1, $2, $3, $4, $5, $6)`,
			nilFilters, card.ID, cs.Reason, cs.ReasonCode, cs.BlockedDate, cs.IsReactivated, time.Now()),
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

func (PSQLCard *psqlCardsRepository) GetCardStatus(c echo.Context, card *model.Card) error {
	cardStatus := model.CardStatuses{}
	err := PSQLCard.DBpg.Model(&cardStatus).
		Where("card_id = ? AND is_reactivated = ?", card.ID, model.BoolNo).
		Order("created_at DESC").
		Limit(1).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return err
	}

	card.CardStatus = cardStatus
	return nil
}

func (PSQLCard *psqlCardsRepository) UpdateOneCardStatus(c echo.Context, cardStatus model.CardStatuses, cols []string) error {
	cardStatus.UpdatedAt = model.NowDbpg()
	cols = append(cols, "updated_at")
	_, err := PSQLCard.DBpg.Model(&cardStatus).Column(cols...).WherePK().Update()

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (PSQLCard *psqlCardsRepository) SetInactiveStatus(c echo.Context, acc model.Account) error {
	var nilFilters []string

	stmts := []*gcdb.PipelineStmt{
		// update account
		gcdb.NewPipelineStmt("UPDATE accounts SET status = $1, updated_at = $2 WHERE id = $3;",
			nilFilters, model.AccStatusInactive, time.Now(), acc.ID),
		// update application
		gcdb.NewPipelineStmt(`UPDATE applications set status = $1, updated_at = $2 WHERE id = $3`,
			nilFilters, model.AppStatusInactive, time.Now(), acc.Application.ID),
		// update cards
		gcdb.NewPipelineStmt(`UPDATE cards set status = $1, updated_at = $2 WHERE id = $3`,
			nilFilters, model.CardStatusInactive, time.Now(), acc.Card.ID),
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
