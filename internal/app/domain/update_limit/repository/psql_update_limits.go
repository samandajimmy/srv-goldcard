package repository

import (
	"database/sql"
	"srv-goldcard/internal/app/domain/update_limit"
	"srv-goldcard/internal/app/model"
	"srv-goldcard/internal/pkg/logger"
	"strconv"

	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
)

type psqlUpdateLimitsRepository struct {
	Conn *sql.DB
	DBpg *pg.DB
}

// NewPsqlUpdateLimitsRepository will create an object that represent the transaction.Repository interface
func NewPsqlUpdateLimitsRepository(Conn *sql.DB, DBpg *pg.DB) update_limit.Repository {
	return &psqlUpdateLimitsRepository{Conn, DBpg}
}

// function to get email address by key
func (psqlUL *psqlUpdateLimitsRepository) GetEmailByKey(c echo.Context) (string, error) {
	param := model.Parameter{}
	err := psqlUL.DBpg.Model(&param).
		Where("key = ?", "UPDATE_LIMIT_EMAIL_ADDRESS").Limit(1).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return "", err
	}

	if err == pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return "", model.ErrGetParameter
	}

	return param.Value, nil
}

// function to get date validation to not allow update limit inquiries
func (psqlUL *psqlUpdateLimitsRepository) GetUpdateLimitInquiriesClosedDate(c echo.Context) (string, error) {
	param := model.Parameter{}
	err := psqlUL.DBpg.Model(&param).
		Where("key = ?", "UPDATE_LIMIT_INQUIRIES_CLOSED_DATE").Limit(1).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return "", err
	}

	if err == pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return "", model.ErrGetParameter
	}

	return param.Value, nil
}

func (psqlUL *psqlUpdateLimitsRepository) GetLastLimitUpdate(c echo.Context, accId int64) (model.LimitUpdate, error) {
	var limitUpdate model.LimitUpdate
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

func (psqlUL *psqlUpdateLimitsRepository) GetAccountBySavingAccount(c echo.Context, savingAcc string) (model.Account, error) {
	acc := model.Account{}
	err := psqlUL.DBpg.Model(&acc).Relation("Application").Relation("PersonalInformation").
		Relation("Card").Relation("Occupation").Relation("EmergencyContact").
		Where("application.saving_account = ? AND account.status = ?", savingAcc, model.AccStatusActive).
		Limit(1).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return acc, err
	}

	if err == pg.ErrNoRows {
		return acc, model.ErrSavingAccNotFound
	}

	return acc, nil
}

func (psqlUL *psqlUpdateLimitsRepository) InsertUpdateCardLimit(c echo.Context, limitUpdt model.LimitUpdate) error {
	now := model.NowDbpg()
	limitUpdt.CreatedAt = now
	limitUpdt.AppliedLimitDate = now
	err := psqlUL.DBpg.Insert(&limitUpdt)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (psqlUL *psqlUpdateLimitsRepository) UpdateCardLimitData(c echo.Context, limitUpdt model.LimitUpdate) error {
	limitUpdt.UpdatedAt = model.NowDbpg()
	col := []string{"status", "updated_at"}
	_, err := psqlUL.DBpg.Model(&limitUpdt).Column(col...).WherePK().Update()

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}

func (psqlUL *psqlUpdateLimitsRepository) GetLimitUpdate(c echo.Context, refId string) (model.LimitUpdate, error) {
	var limitUpdt model.LimitUpdate

	err := psqlUL.DBpg.Model(&limitUpdt).
		Where("ref_id = ? AND limit_update.status = ?", refId, model.LimitUpdateStatusInquired).
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

func (psqlUL *psqlUpdateLimitsRepository) GetsertGtePayment(c echo.Context, pl model.PayloadCoreGtePayment) (model.GtePayment, error) {
	gtePayment := model.GtePayment{}
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
		return gtePayment, model.ErrGtePaymenTrxIdExist
	}

	goldAmt, _ := strconv.ParseFloat(pl.AvailableGram, 64)
	gtePayment = model.GtePayment{
		AccountId:  acc.ID,
		TrxId:      pl.TrxId,
		GoldAmount: goldAmt,
		TrxAmount:  pl.NominalTransaction,
		CreatedAt:  model.NowDbpg(),
	}

	err = psqlUL.DBpg.Insert(&gtePayment)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return gtePayment, err
	}

	return gtePayment, nil
}

func (psqlUL *psqlUpdateLimitsRepository) UpdateGtePayment(c echo.Context, gtePayment model.GtePayment, cols []string) error {
	gtePayment.UpdatedAt = model.NowDbpg()
	cols = append(cols, "updated_at")
	_, err := psqlUL.DBpg.Model(&gtePayment).Column(cols...).WherePK().Update()

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}
