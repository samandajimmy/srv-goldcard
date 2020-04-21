package repository

import (
	"database/sql"
	"gade/srv-goldcard/cards"
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

func (PSQLCard *psqlCardsRepository) UpdateCardStatus(c echo.Context, card models.Card) error {
	card.UpdatedAt = models.NowDbpg()
	col := []string{"status", "updated_at"}
	_, err := PSQLCard.DBpg.Model(&card).Column(col...).WherePK().Update()

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

// Insert Card Statuses
func (PSQLCard *psqlCardsRepository) PostCardStatuses(cs models.CardStatuses) error {
	cs.CreatedAt = time.Now()
	err := PSQLCard.DBpg.Insert(&cs)

	if err != nil {
		return err
	}

	return nil
}
