package repository

import (
	"database/sql"
	"srv-goldcard/internal/app/domain/activation"
	"srv-goldcard/internal/app/model"
	gcdb "srv-goldcard/internal/pkg/database"
	"srv-goldcard/internal/pkg/logger"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
)

type psqlActivations struct {
	Conn *sql.DB
	DBpg *pg.DB
}

// NewPsqlActivations will create an object that represent the activation.Repository interface
func NewPsqlActivations(Conn *sql.DB, dbpg *pg.DB) activation.Repository {
	return &psqlActivations{Conn, dbpg}
}

func (pa *psqlActivations) PostActivations(c echo.Context, acc model.Account) error {
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
		gcdb.NewPipelineStmt(`UPDATE cards SET status = $1, card_number = $2, valid_until = $3,
			encrypted_card_number = $4, activated_date = $5, updated_at = $6
			WHERE id = $7`,
			nilFilters, card.Status, card.CardNumber, card.ValidUntil, card.EncryptedCardNumber,
			card.ActivatedDate, time.Now(), acc.CardID),
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

func (pa *psqlActivations) UpdateGoldLimit(c echo.Context, card model.Card) error {
	card.UpdatedAt = model.NowDbpg()
	col := []string{"gold_limit", "gold_balance", "stl_limit", "stl_balance", "updated_at"}
	_, err := pa.DBpg.Model(&card).Column(col...).WherePK().Update()

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (act *psqlActivations) GetStoredGoldPrice(c echo.Context) (int64, error) {
	var goldPrice model.GoldPrice

	query := `select id, price, valid_date, created_at
		from gold_prices order by created_at desc limit 1;`

	_, err := act.DBpg.Query(&goldPrice, query)

	if err != nil || (goldPrice == model.GoldPrice{}) {
		logger.Make(nil, nil).Debug(err)

		return 0, err
	}

	return int64(goldPrice.Price), nil
}
