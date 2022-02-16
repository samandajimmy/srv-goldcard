package activations

import (
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// UseCase represent the activations usecases
type UseCase interface {
	ForceActivation(c echo.Context, acc models.Account) (models.RespActivations, error)
	InquiryActivation(c echo.Context, acc models.Account) (models.CardBalance, models.ResponseErrors)
	PostActivations(echo.Context, models.PayloadActivations) (models.RespActivations, error)
	PostReactivations(echo.Context, models.PayloadActivations) (models.RespActivations, error)
	ValidateActivation(c echo.Context, pa models.PayloadActivations) models.ResponseErrors
}
