package repository

import (
	"database/sql"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/update_limits"

	"github.com/go-pg/pg/v9"
)

type psqlUpdateLimitsRepository struct {
	Conn *sql.DB
	DBpg *pg.DB
}

// NewPsqlUpdateLimitsRepository will create an object that represent the transactions.Repository interface
func NewPsqlUpdateLimitsRepository(Conn *sql.DB, DBpg *pg.DB) update_limits.Repository {
	return &psqlUpdateLimitsRepository{Conn, DBpg}
}

func (psqlUL *psqlUpdateLimitsRepository) GetParameterByKey(key string) (models.Parameter, error) {
	var param models.Parameter
	query := `select id, key, value, description, created_at, updated_at
		from parameters where key = ? limit 1;`

	_, err := psqlUL.DBpg.Query(&param, query, key)

	if err != nil || (param == models.Parameter{}) {
		logger.Make(nil, nil).Debug(err)

		return param, err
	}

	return param, nil
}
