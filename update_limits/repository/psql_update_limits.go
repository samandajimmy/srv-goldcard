package repository

import (
	"database/sql"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/update_limits"

	"github.com/go-pg/pg/v9"
	"github.com/google/uuid"
	"github.com/labstack/echo"
)

type psqlUpdateLimitsRepository struct {
	Conn *sql.DB
	DBpg *pg.DB
}

// NewPsqlUpdateLimitsRepository will create an object that represent the transactions.Repository interface
func NewPsqlUpdateLimitsRepository(Conn *sql.DB, DBpg *pg.DB) update_limits.Repository {
	return &psqlUpdateLimitsRepository{Conn, DBpg}
}

// function to get email address by key
func (psqlUL *psqlUpdateLimitsRepository) GetEmailByKey(c echo.Context) (string, error) {
	param := models.Parameter{}
	err := psqlUL.DBpg.Model(&param).
		Where("key = ?", "UPDATE_LIMIT_EMAIL_ADDRESS").Limit(1).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return "", err
	}

	if err == pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return "", models.ErrGetParameter
	}

	return param.Value, nil
}

func (psqlUL *psqlUpdateLimitsRepository) GetLastLimitUpdate(c echo.Context, accId int64) (models.LimitUpdate, error) {
	var limitUpdate models.LimitUpdate
	statuses := []string{"pending", "applied"}
	err := psqlUL.DBpg.Model(&limitUpdate).
		Where("account_id = ? AND status in (?)", accId, pg.In(statuses)).Order("created_at DESC").
		Limit(1).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return limitUpdate, err
	}

	return limitUpdate, nil
}

func (psqlUL *psqlUpdateLimitsRepository) GetAccountBySavingAccount(c echo.Context, savingAcc string) (models.Account, error) {
	acc := models.Account{}
	err := psqlUL.DBpg.Model(&acc).Relation("Application").Relation("PersonalInformation").
		Relation("Card").Relation("Occupation").Relation("Correspondence").Relation("EmergencyContact").
		Where("application.saving_account = ? AND account.status = ?", savingAcc, models.AccStatusActive).
		Limit(1).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return acc, err
	}

	if err == pg.ErrNoRows {
		return acc, models.ErrSavingAccNotFound
	}

	return acc, nil
}

func (psqlUL *psqlUpdateLimitsRepository) InsertUpdateCardLimit(c echo.Context, limitUpdt models.LimitUpdate) (string, error) {
	refId, _ := uuid.NewRandom()
	now := models.NowDbpg()
	limitUpdt.CreatedAt = now
	limitUpdt.AppliedLimitDate = now
	limitUpdt.RefId = refId.String()
	err := psqlUL.DBpg.Insert(&limitUpdt)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return "", err
	}

	return limitUpdt.RefId, nil
}

func (psqlUL *psqlUpdateLimitsRepository) UpdateCardLimitData(c echo.Context, limitUpdt models.LimitUpdate) error {
	limitUpdt.UpdatedAt = models.NowDbpg()
	col := []string{"status", "updated_at"}
	_, err := psqlUL.DBpg.Model(&limitUpdt).Column(col...).WherePK().Update()

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (psqlUL *psqlUpdateLimitsRepository) GetLimitUpdate(c echo.Context, refId string) (models.LimitUpdate, error) {
	var limitUpdt models.LimitUpdate

	err := psqlUL.DBpg.Model(&limitUpdt).
		Where("ref_id = ? AND limit_update.status = ?", refId, models.LimitUpdateStatusInquired).
		Limit(1).Select()

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return limitUpdt, err
	}

	err = psqlUL.DBpg.Model(&limitUpdt.Account).Relation("Application").Relation("PersonalInformation").
		Where("account.id = ?", &limitUpdt.AccountID).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return limitUpdt, err
	}

	err = psqlUL.DBpg.Model(&limitUpdt.Account.Application.Documents).Where("application_id = ?", &limitUpdt.Account.ApplicationID).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return limitUpdt, err
	}

	return limitUpdt, nil
}
