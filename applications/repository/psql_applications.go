package repository

import (
	"database/sql"
	"gade/srv-goldcard/applications"
	"gade/srv-goldcard/models"
	"time"

	"github.com/labstack/echo"
)

type psqlApplicationsRepository struct {
	Conn *sql.DB
}

// NewpsqlApplicationsRepository will create an object that represent the applications.Repository interface
func NewPsqlApplicationsRepository(Conn *sql.DB) applications.Repository {
	return &psqlApplicationsRepository{Conn}
}

// PostAccountNumber representation update account number to database
func (appli *psqlApplicationsRepository) PostSavingAccount(c echo.Context, applications *models.Applications) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	var lastID int64
	now := time.Now()
	query := `UPDATE applications
		set saving_account = $1,
			updated_at = $3
		WHERE application_number = $2 RETURNING id`
	stmt, err := appli.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	err = stmt.QueryRow(applications.SavingAccount, applications.ApplicationNumber, &now).Scan(&lastID)
	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	return nil
}
