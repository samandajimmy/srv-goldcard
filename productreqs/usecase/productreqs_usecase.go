package usecase

import (
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/productreqs"

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

// InsertPublicHolidayDate represent to insert public holiday date
func (prodreqs *productreqsUseCase) InsertPublicHolidayDate(c echo.Context, phd models.PayloadInsertPublicHoliday) (models.PublicHolidayDate, error) {
	pubHoliDate, err := prodreqs.prRepo.GetSertPublicHolidayDate(c, phd.PublicHolidayDate)

	if err != nil {
		return models.PublicHolidayDate{}, err
	}

	return models.PublicHolidayDate{
		PublicHolidayDate: pubHoliDate,
	}, nil
}

// GetPublicHolidayDate represent to get public holiday date
func (prodreqs *productreqsUseCase) GetPublicHolidayDate(c echo.Context) (models.PublicHolidayDate, error) {
	pubHoliDate, err := prodreqs.prRepo.GetSertPublicHolidayDate(c, []string{})

	if err != nil {
		return models.PublicHolidayDate{}, err
	}

	return models.PublicHolidayDate{
		PublicHolidayDate: pubHoliDate,
	}, nil
}
