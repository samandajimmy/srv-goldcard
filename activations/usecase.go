package activations

import "github.com/labstack/echo"

import "gade/srv-goldcard/models"

// UseCase represent the activations usecases
type UseCase interface {
	InquiryActivation(c echo.Context, pl models.PayloadAppNumber) models.ResponseErrors
	PostActivations(echo.Context, models.PayloadActivations) error
}
