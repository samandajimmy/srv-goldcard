package repository

import (
	"database/sql"
	"gade/srv-goldcard/activations"

	"github.com/go-pg/pg/v9"
)

type psqlActivations struct {
	Conn *sql.DB
	DBpg *pg.DB
}

// NewPsqlActivations will create an object that represent the activations.Repository interface
func NewPsqlActivations(Conn *sql.DB, dbpg *pg.DB) activations.Repository {
	return &psqlActivations{Conn, dbpg}
}
