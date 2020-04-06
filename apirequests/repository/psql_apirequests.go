package repository

import (
	"database/sql"
	"gade/srv-goldcard/apirequests"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"

	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
)

type psqlAPIRequests struct {
	Conn *sql.DB
	DBpg *pg.DB
}

// NewPsqlAPIRequestsRepository will create an object that represent the apirequests.Repository interface
func NewPsqlAPIRequestsRepository(Conn *sql.DB, dbpg *pg.DB) apirequests.Repository {
	return &psqlAPIRequests{Conn, dbpg}
}

func (par *psqlAPIRequests) InserAPIRequest(c echo.Context, ar models.APIRequest) error {
	ar.CreatedAt = models.NowDbpg()
	err := par.DBpg.Insert(&ar)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}
