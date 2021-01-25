package repository

import (
	"database/sql"
	"gade/srv-goldcard/cards"
	gcdb "gade/srv-goldcard/database"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"time"

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
	var nilFilters []string
	logger.MakeStructToJSON(cs)

	cardStatsQuery := `INSERT INTO card_statuses (card_id, reason, reason_code, blocked_date, is_reactivated, created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	cardStatsParams := []interface{}{card.ID, cs.Reason, cs.ReasonCode, cs.BlockedDate, cs.IsReactivated, time.Now()}

	if cs.ID != 0 {
		cardStatsQuery = `UPDATE card_statuses SET is_reactivated = $1, last_encrypted_card_number = $2, reactivated_date = $3, updated_at = $4 where id = $5`
		cardStatsParams = []interface{}{cs.IsReactivated, cs.LastEncryptedCardNumber, cs.ReactivatedDate, time.Now(), cs.ID}
	}

	stmts := []*gcdb.PipelineStmt{
		gcdb.NewPipelineStmt(`UPDATE cards SET status = $1, updated_at = $2 WHERE id = $3`,
			nilFilters, card.Status, time.Now(), card.ID),
		gcdb.NewPipelineStmt(cardStatsQuery,
			nilFilters, cardStatsParams...),
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

func (PSQLCard *psqlCardsRepository) GetCardStatus(c echo.Context, card *models.Card) error {
	cardStatus := models.CardStatuses{}
	err := PSQLCard.DBpg.Model(&cardStatus).
		Where("card_id = ? AND is_reactivated = ?", card.ID, models.BoolNo).
		Order("created_at DESC").
		Limit(1).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return err
	}

	card.CardStatus = cardStatus
	return nil
}
