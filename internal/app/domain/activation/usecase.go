package activation

import (
	"srv-goldcard/internal/app/model"

	"github.com/labstack/echo"
)

// UseCase represent the activation usecases
type UseCase interface {
	ForceActivation(c echo.Context, acc model.Account) (model.RespActivations, error)
	InquiryActivation(c echo.Context, acc model.Account) (model.CardBalance, model.ResponseErrors)
	PostActivations(echo.Context, model.PayloadActivations) (model.RespActivations, error)
	PostReactivations(echo.Context, model.PayloadActivations) (model.RespActivations, error)
	ValidateActivation(c echo.Context, pa model.PayloadActivations) model.ResponseErrors
}
