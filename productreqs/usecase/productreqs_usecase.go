package usecase

import (
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/productreqs"

	"github.com/labstack/echo"
)

type productreqsUseCase struct{}

// ProductReqsUseCase represent product requirements Use Case
func ProductReqsUseCase() productreqs.UseCase {
	return &productreqsUseCase{}
}

// ProductRequirements represent to get all product requirements
func (prodreqs *productreqsUseCase) ProductRequirements(c echo.Context) (models.Requirements, error) {
	return models.RequirementsValue, nil
}
