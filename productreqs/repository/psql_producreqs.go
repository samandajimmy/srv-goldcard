package repository

import (
	"database/sql"
	"gade/srv-goldcard/logger"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/productreqs"
	"strings"

	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
)

type psqlProductReqsRepository struct {
	Conn *sql.DB
	DBpg *pg.DB
}

// NewPsqlProductReqsRepository will create an object that represent the productreqs.Repository interface
func NewPsqlProductReqsRepository(Conn *sql.DB, dbpg *pg.DB) productreqs.Repository {
	return &psqlProductReqsRepository{Conn, dbpg}
}

func (prodReq *psqlProductReqsRepository) GetSertPublicHolidayDate(c echo.Context, phds []string) (string, error) {
	pubHoliDate := models.Parameter{}

	err := prodReq.DBpg.Model(&pubHoliDate).
		Where("key = ?", "PUBLIC_HOLIDAY_DATE").
		Limit(1).Select()

	if err != nil && err != pg.ErrNoRows {
		logger.Make(c, nil).Debug(err)

		return "", err
	}

	// append inputted public holiday date in database with inputted request payload
	pubHoliDateVal := strings.Replace(pubHoliDate.Value, ",", "", -1)
	appendedPubHoliDateVal := models.UniquifyStringSlice(append(strings.Fields(pubHoliDateVal), phds...))

	// set public holiday date param value to string with delimited comma
	pubHoliDate.Value = strings.Join(appendedPubHoliDateVal, ", ")
	pubHoliDate.UpdatedAt = models.NowDbpg()
	_, err = prodReq.DBpg.Model(&pubHoliDate).Column([]string{"value", "updated_at"}...).WherePK().Update()

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return "", err
	}

	return pubHoliDate.Value, nil
}
