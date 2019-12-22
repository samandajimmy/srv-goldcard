package repository

import (
	"database/sql"
	"gade/srv-goldcard/models"
	"time"

	"github.com/labstack/echo"
)

type psqlRegistrationsRepository struct {
	Conn *sql.DB
}

func (regis *psqlRegistrationsRepository) PostRegistrations(c echo.Context, registrations models.Registrations) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	var lastID int64
	now := time.Now()
	query := `UPDATE personal_informations 
		set residence_address = $1,
			updated_at = $3
		WHERE phone_number = $2`
	stmt, err := regis.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	err = stmt.QueryRow(registrations.ResidenceAddress, registrations.PhoneNumber, &now).Scan(&lastID)
	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	return nil
}
