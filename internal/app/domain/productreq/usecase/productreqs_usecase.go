package usecase

import (
	"srv-goldcard/internal/app/domain/productreq"
	"srv-goldcard/internal/app/model"

	"github.com/labstack/echo"
)

type productreqsUseCase struct {
	prRepo productreq.Repository
}

// ProductReqsUseCase represent product requirements Use Case
func ProductReqsUseCase(prRepo productreq.Repository) productreq.UseCase {
	return &productreqsUseCase{prRepo}
}

// ProductRequirements represent to get all product requirements
func (prodreqs *productreqsUseCase) ProductRequirements(c echo.Context) (model.Requirements, error) {
	return model.RequirementsValue, nil
}

// InsertPublicHolidayDate represent to insert public holiday date
func (prodreqs *productreqsUseCase) InsertPublicHolidayDate(c echo.Context, phd model.PayloadInsertPublicHoliday) (model.PublicHolidayDate, error) {
	pubHoliDate, err := prodreqs.prRepo.GetSertPublicHolidayDate(c, phd.PublicHolidayDate)

	if err != nil {
		return model.PublicHolidayDate{}, err
	}

	return model.PublicHolidayDate{
		PublicHolidayDate: pubHoliDate,
	}, nil
}

// GetPublicHolidayDate represent to get public holiday date
func (prodreqs *productreqsUseCase) GetPublicHolidayDate(c echo.Context) (model.PublicHolidayDate, error) {
	pubHoliDate, err := prodreqs.prRepo.GetSertPublicHolidayDate(c, []string{})

	if err != nil {
		return model.PublicHolidayDate{}, err
	}

	return model.PublicHolidayDate{
		PublicHolidayDate: pubHoliDate,
	}, nil
}
