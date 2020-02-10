package activations

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// UseCase represent the activations usecases
type UseCase interface {
	InquiryActivation(c echo.Context, pl models.PayloadAppNumber) models.ResponseErrors
	PostActivations(echo.Context, models.PayloadActivations) (models.RespActivations, error)
}
