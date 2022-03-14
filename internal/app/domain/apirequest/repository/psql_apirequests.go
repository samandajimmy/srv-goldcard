package repository

import (
	"database/sql"
	"srv-goldcard/internal/app/domain/apirequest"
	"srv-goldcard/internal/app/model"
	"srv-goldcard/internal/pkg/logger"

	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
)

type psqlAPIRequests struct {
	Conn *sql.DB
	DBpg *pg.DB
}

// NewPsqlAPIRequestsRepository will create an object that represent the apirequest.Repository interface
func NewPsqlAPIRequestsRepository(Conn *sql.DB, dbpg *pg.DB) apirequest.Repository {
	return &psqlAPIRequests{Conn, dbpg}
}

func (par *psqlAPIRequests) InserAPIRequest(c echo.Context, ar model.APIRequest) error {
	ar.CreatedAt = model.NowDbpg()
	err := par.DBpg.Insert(&ar)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}
