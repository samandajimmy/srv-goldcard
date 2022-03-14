package repository

import (
	"database/sql"
	"srv-goldcard/internal/app/domain/productreq"
	"srv-goldcard/internal/app/model"
	"srv-goldcard/internal/pkg/logger"
	"strings"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo"
)

type psqlProductReqsRepository struct {
	Conn *sql.DB
	DBpg *pg.DB
}

// NewPsqlProductReqsRepository will create an object that represent the productreq.Repository interface
func NewPsqlProductReqsRepository(Conn *sql.DB, dbpg *pg.DB) productreq.Repository {
	return &psqlProductReqsRepository{Conn, dbpg}
}

func (prodReq *psqlProductReqsRepository) GetSertPublicHolidayDate(c echo.Context, phds []string) (string, error) {
	pubHoliDate := model.Parameter{}

	err := prodReq.DBpg.Model(&pubHoliDate).
		Where("key = ?", "PUBLIC_HOLIDAY_DATE").
		Limit(1).Select()

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return "", err
	}

	var newDataExist bool = false
	for _, data := range phds {
		// validating inputed holiday date
		_, err = time.Parse("02/01/2006", data)

		if err != nil {
			return "", model.ErrDateFormat
		}

		// if value not exist in public holiday date then append
		if !strings.Contains(pubHoliDate.Value, data) {
			pubHoliDate.Value += ", " + data
			newDataExist = true
		}
	}

	// if there is no date then return response
	if !newDataExist {
		return pubHoliDate.Value, nil
	}

	// if new data exist then do update
	if err = prodReq.updatePublicHolidayDate(c, pubHoliDate); err != nil {
		return "", err
	}

	return pubHoliDate.Value, nil
}

func (prodReq *psqlProductReqsRepository) updatePublicHolidayDate(c echo.Context, pubHoliDate model.Parameter) error {
	pubHoliDate.UpdatedAt = model.NowDbpg()
	_, err := prodReq.DBpg.Model(&pubHoliDate).Column([]string{"value", "updated_at"}...).WherePK().Update()

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return err
	}

	return nil
}
