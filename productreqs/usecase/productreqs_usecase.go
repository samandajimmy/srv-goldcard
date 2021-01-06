package usecase

import (
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/productreqs"
	"time"

	"github.com/labstack/echo"
)

type productreqsUseCase struct {
	prRepo productreqs.Repository
}

// ProductReqsUseCase represent product requirements Use Case
func ProductReqsUseCase(prRepo productreqs.Repository) productreqs.UseCase {
	return &productreqsUseCase{prRepo}
}

// ProductRequirements represent to get all product requirements
func (prodreqs *productreqsUseCase) ProductRequirements(c echo.Context) (models.Requirements, error) {
	return models.RequirementsValue, nil
}

// GetSertPublicHolidayDate represent to get or insert public holiday date
func (prodreqs *productreqsUseCase) GetSertPublicHolidayDate(c echo.Context, phd models.PayloadGetSertPublicHoliday) (models.PublicHolidayDate, error) {
	var err error
	for _, data := range phd.PublicHolidayDate {
		// validating inputed holiday date
		_, err = time.Parse("02/01/2006", data)

		if err != nil {
			return models.PublicHolidayDate{}, models.ErrDateFormat
		}
	}

	pubHoliDate, err := prodreqs.prRepo.GetSertPublicHolidayDate(c, phd.PublicHolidayDate)

	if err != nil {
		return models.PublicHolidayDate{}, err
	}

	return models.PublicHolidayDate{
		PublicHolidayDate: pubHoliDate,
	}, nil
}
