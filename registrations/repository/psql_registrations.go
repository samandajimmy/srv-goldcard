package repository

import (
	"database/sql"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/registrations"
	"time"

	"github.com/labstack/echo"
)

type psqlRegistrationsRepository struct {
	Conn *sql.DB
}

// NewPsqlRegistrationsRepository will create an object that represent the registrations.Repository interface
func NewPsqlRegistrationsRepository(Conn *sql.DB) registrations.Repository {
	return &psqlRegistrationsRepository{Conn}
}

// PostAddress representation update address to database
func (regis *psqlRegistrationsRepository) PostAddress(c echo.Context, registrations *models.Registrations) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	var lastID int64
	now := time.Now()
	query := `UPDATE personal_informations 
		set residence_address = $1,
			updated_at = $3
		WHERE phone_number = $2 RETURNING id`
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

// PostAddress representation get address from database
func (regis *psqlRegistrationsRepository) GetAddress(c echo.Context, phoneNo string) (string, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	var address string
	query := `SELECT residence_address from personal_informations WHERE phone_number = $1`
	stmt, err := regis.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return "", err
	}

	err = stmt.QueryRow(phoneNo).Scan(&address)
	if err != nil {
		requestLogger.Debug(err)

		return "", err
	}
	return address, nil
}

// PostSavingAccount representation update saving account to database
func (regis *psqlRegistrationsRepository) PostSavingAccount(c echo.Context, applications *models.Applications) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	var lastID int64
	now := time.Now()
	query := `UPDATE applications
		set saving_account = $1,
			updated_at = $3
		WHERE application_number = $2 RETURNING id`
	stmt, err := regis.Conn.Prepare(query)

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
