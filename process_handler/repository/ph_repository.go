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

// Get all process handler for validation insert. if return true = insert
func (psqlPH *psqlProcHandler) GetProcessHandler(ps models.ProcessStatus) (bool, error) {

	err := psqlPH.DBpg.Model(&ps).
		Where("process_id = ?", ps.ProcessID).
		Where("tbl_name = ?", ps.TblName).Select()

	if err != nil && err != pg.ErrNoRows {
		return false, err
	}

	if err == pg.ErrNoRows {
		return true, nil
	}

	return false, nil
}

// Insert Process handler
func (psqlPH *psqlProcHandler) PostProcessHandler(ps models.ProcessStatus) error {
	err := psqlPH.DBpg.Insert(&ps)

	if err != nil {
		return err
	}

	return nil
}
