package repository

import (
	"database/sql"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/update_limits"

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

func (psqlUL *psqlUpdateLimitsRepository) GetDocumentByTypeAndApplicationId(appId int64, docType string) (models.Document, error) {
	var listDocument models.Document
	err := psqlUL.DBpg.Model(&listDocument).
		Where("application_id = ? AND type = ?", appId, docType).Order("created_at DESC").Limit(1).Select()

	if err != nil {
		logger.Make(nil, nil).Debug(err)

		return listDocument, err
	}

	return listDocument, nil
}

func (psqlUL *psqlUpdateLimitsRepository) GetLastLimitUpdate(accId int64) (models.LimitUpdate, error) {
	var limitUpdate models.LimitUpdate
	err := psqlUL.DBpg.Model(&limitUpdate).
		Where("account_id = ?", accId).Order("created_at DESC").Limit(1).Select()

	if err != nil {
		logger.Make(nil, nil).Debug(err)

		return limitUpdate, err
	}

	return limitUpdate, nil
}
