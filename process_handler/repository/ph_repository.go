package repository

import (
	"database/sql"

	"gade/srv-goldcard/models"
	"gade/srv-goldcard/process_handler"

	"github.com/go-pg/pg/v9"
)

type psqlProcHandler struct {
	Conn *sql.DB
	DBpg *pg.DB
}

func NewPsqlProcHandlerRepository(Conn *sql.DB, dbpg *pg.DB) process_handler.Repository {
	return &psqlProcHandler{Conn, dbpg}
}

func (psqlPH *psqlProcHandler) GetProcessHandler(processID string) (models.ProcessStatus, error) {
	ps := models.ProcessStatus{}

	err := psqlPH.DBpg.Model(&ps).Where("process_id = ?", processID).Select()

	if err != nil {
		return ps, err
	}

	return ps, nil
}

func (psqlPH *psqlProcHandler) PostProcessHandler(ps models.ProcessStatus) error {
	err := psqlPH.DBpg.Insert(&ps)

	if err != nil {
		return err
	}

	return nil
}

func (psqlPH *psqlProcHandler) PutProcessHandler(ps models.ProcessStatus) error {
	err := psqlPH.DBpg.Update(&ps)

	if err != nil {
		return err
	}

	return nil
}
