package repository

import (
	"database/sql"
	"gade/srv-goldcard/logger"
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

func (regis *psqlRegistrationsRepository) CreateApplication(c echo.Context, app models.Applications,
	acc models.Account, pi models.PersonalInformation) error {
	var appID, piID int64
	tx, err := regis.Conn.Begin()

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	{
		stmt, err := tx.Prepare(`INSERT INTO applications (application_number) VALUES ($1)
			RETURNING id;`)

		if err != nil {
			logger.Make(c, nil).Debug(err)
			tx.Rollback()

			return err
		}

		defer stmt.Close()

		err = stmt.QueryRow(app.ApplicationNumber).Scan(&appID)

		if err != nil {
			logger.Make(c, nil).Debug(err)
			tx.Rollback()

			return err
		}
	}

	{
		stmt, err := tx.Prepare(`INSERT INTO personal_informations (phone_number) VALUES ($1)
			RETURNING id;`)

		if err != nil {
			logger.Make(c, nil).Debug(err)
			tx.Rollback()

			return err
		}

		defer stmt.Close()

		err = stmt.QueryRow(pi.PhoneNumber).Scan(&piID)

		if err != nil {
			logger.Make(c, nil).Debug(err)
			tx.Rollback()

			return err
		}
	}

	{
		stmt, err := tx.Prepare(`INSERT INTO accounts (cif, bank_id, application_id,
			personal_information_id) VALUES ($1, $2, $3, $4);`)

		if err != nil {
			logger.Make(c, nil).Debug(err)
			tx.Rollback()

			return err
		}

		defer stmt.Close()

		if _, err := stmt.Exec(acc.CIF, acc.BankID, appID, piID); err != nil {
			logger.Make(c, nil).Debug(err)
			tx.Rollback()

			return err
		}
	}

	if err := tx.Commit(); err != nil {
		logger.Make(c, nil).Debug(err)
		tx.Rollback()

		return err
	}

	return nil
}

func (regis *psqlRegistrationsRepository) GetBankIDByCode(c echo.Context, bankCode string) (int64, error) {
	var bankID int64
	query := `SELECT id FROM banks WHERE code = $1`
	err := regis.Conn.QueryRow(query, bankCode).Scan(&bankID)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return bankID, err
	}

	return bankID, nil
}

func (regis *psqlRegistrationsRepository) PostAddress(c echo.Context, registrations *models.Registrations) error {
	var lastID int64
	now := time.Now()
	query := `UPDATE personal_informations
		set residence_address = $1, updated_at = $3
		WHERE phone_number = $2 RETURNING id`
	stmt, err := regis.Conn.Prepare(query)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	err = stmt.QueryRow(registrations.ResidenceAddress, registrations.PhoneNumber, &now).Scan(&lastID)
	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (regis *psqlRegistrationsRepository) GetAddress(c echo.Context, phoneNo string) (string, error) {
	var address string
	query := `SELECT residence_address from personal_informations WHERE phone_number = $1`
	stmt, err := regis.Conn.Prepare(query)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return "", err
	}

	err = stmt.QueryRow(phoneNo).Scan(&address)
	if err != nil {
		logger.Make(c, nil).Debug(err)

		return "", err
	}
	return address, nil
}

func (regis *psqlRegistrationsRepository) PostSavingAccount(c echo.Context, applications *models.Applications) error {
	var lastID int64
	now := time.Now()
	query := `UPDATE applications
		set saving_account = $1,
			updated_at = $3
		WHERE application_number = $2 RETURNING id`
	stmt, err := regis.Conn.Prepare(query)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	err = stmt.QueryRow(applications.SavingAccount, applications.ApplicationNumber, &now).Scan(&lastID)
	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}
