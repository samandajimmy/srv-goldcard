package repository

import (
	"database/sql"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/update_limits"
	"strconv"

	"github.com/go-pg/pg/v9"
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

// function to get date validation to not allow update limit inquiries
func (psqlUL *psqlUpdateLimitsRepository) GetUpdateLimitInquiriesClosedDate(c echo.Context) (string, error) {
	param := models.Parameter{}
	err := psqlUL.DBpg.Model(&param).
		Where("key = ?", "UPDATE_LIMIT_INQUIRIES_CLOSED_DATE").Limit(1).Select()

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
		Relation("Card").Relation("Occupation").Relation("EmergencyContact").
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

func (psqlUL *psqlUpdateLimitsRepository) InsertUpdateCardLimit(c echo.Context, limitUpdt models.LimitUpdate) error {
	now := models.NowDbpg()
	limitUpdt.CreatedAt = now
	limitUpdt.AppliedLimitDate = now
	err := psqlUL.DBpg.Insert(&limitUpdt)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
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

func (psqlUL *psqlUpdateLimitsRepository) GetsertGtePayment(c echo.Context, pl models.PayloadCoreGtePayment) (models.GtePayment, error) {
	gtePayment := models.GtePayment{}
	acc, err := psqlUL.GetAccountBySavingAccount(c, pl.SavingAccount)

	if err != nil {
		return gtePayment, err
	}

	err = psqlUL.DBpg.Model(&gtePayment).
		Where("trx_id = ?", pl.TrxId).
		Limit(1).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return gtePayment, err
	}

	gtePayment.Account = acc
	// if the trx id existed with no errors
	if gtePayment.ID != 0 && (!gtePayment.BriUpdated || !gtePayment.PdsNotified) {
		return gtePayment, nil
	}

	// if the trx id existed but there are still errors
	if gtePayment.ID != 0 {
		return gtePayment, models.ErrGtePaymenTrxIdExist
	}

	goldAmt, _ := strconv.ParseFloat(pl.AvailableGram, 64)
	gtePayment = models.GtePayment{
		AccountId:  acc.ID,
		TrxId:      pl.TrxId,
		GoldAmount: goldAmt,
		TrxAmount:  pl.NominalTransaction,
		CreatedAt:  models.NowDbpg(),
	}

	err = psqlUL.DBpg.Insert(&gtePayment)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return gtePayment, err
	}

	return gtePayment, nil
}

func (psqlUL *psqlUpdateLimitsRepository) UpdateGtePayment(c echo.Context, gtePayment models.GtePayment, cols []string) error {
	gtePayment.UpdatedAt = models.NowDbpg()
	cols = append(cols, "updated_at")
	_, err := psqlUL.DBpg.Model(&gtePayment).Column(cols...).WherePK().Update()

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}
