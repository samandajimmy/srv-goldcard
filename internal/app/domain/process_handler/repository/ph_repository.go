package repository

import (
	"database/sql"

	"srv-goldcard/internal/app/domain/process_handler"
	"srv-goldcard/internal/app/model"

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
func (psqlPH *psqlProcHandler) GetProcessHandler(ps model.ProcessStatus) (model.ProcessStatus, error) {
	newPs := model.ProcessStatus{}
	err := psqlPH.DBpg.Model(&newPs).
		Where("process_id = ?", ps.ProcessID).
		Where("tbl_name = ?", ps.TblName).Select()

	if err != nil && err != pg.ErrNoRows {
		return newPs, err
	}

	if err == pg.ErrNoRows {
		return newPs, nil
	}

	return newPs, nil
}

// Insert Process handler
func (psqlPH *psqlProcHandler) PostProcessHandler(ps model.ProcessStatus) error {
	ps.CreatedAt = model.NowDbpg()
	err := psqlPH.DBpg.Insert(&ps)

	if err != nil {
		return err
	}

	return nil
}

// Update Process Handler
func (psqlPH *psqlProcHandler) UpdateProcessHandler(ps model.ProcessStatus, col []string) error {
	ps.UpdatedAt = model.NowDbpg()

	_, err := psqlPH.DBpg.Model(&ps).Column(col...).
		Where("process_id = ?", ps.ProcessID).
		Where("tbl_name = ?", ps.TblName).Update()

	if err != nil {
		return err
	}

	return nil
}
