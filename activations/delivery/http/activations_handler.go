package http

import (
	"gade/srv-goldcard/activations"
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// ActivationsHandler represent the httphandler for activations
type ActivationsHandler struct {
	response   models.Response
	respErrors models.ResponseErrors
	aUsecase   activations.UseCase
}

// NewActivationsHandler represent to activations gold card
func NewActivationsHandler(echoGroup models.EchoGroup, aUseCase activations.UseCase) {
	handler := &ActivationsHandler{aUsecase: aUseCase}

	// End Point For External
	echoGroup.API.POST("/activations", handler.Activations)
	echoGroup.API.POST("/activations/inquiry", handler.ActivationsInquiry)
}

// ActivationsInquiry a handler to handle goldcard inquiry activations
func (ah *ActivationsHandler) ActivationsInquiry(c echo.Context) error {
	var pl models.PayloadAppNumber
	ah.response, ah.respErrors = models.NewResponse()

	if err := c.Bind(&pl); err != nil {
		ah.respErrors.SetTitle(models.MessageUnprocessableEntity)
		ah.response.SetResponse("", &ah.respErrors)

		return ah.response.Body(c, err)
	}

	if err := c.Validate(pl); err != nil {
		ah.respErrors.SetTitle(err.Error())
		ah.response.SetResponse("", &ah.respErrors)

		return ah.response.Body(c, err)
	}

	err := ah.aUsecase.InquiryActivation(c, pl)

	if err != nil {
		ah.respErrors.SetTitle(err.Error())
		ah.response.SetResponse("", &ah.respErrors)

		return ah.response.Body(c, err)
	}

	ah.response.SetResponse("", &ah.respErrors)

	return ah.response.Body(c, err)
}

// Activations a handler to activations
func (ah *ActivationsHandler) Activations(c echo.Context) error {
	var pa models.PayloadActivations
	ah.response, ah.respErrors = models.NewResponse()

	if err := c.Bind(&pa); err != nil {
		ah.respErrors.SetTitle(models.MessageUnprocessableEntity)
		ah.response.SetResponse("", &ah.respErrors)

		return ah.response.Body(c, err)
	}

	if err := c.Validate(pa); err != nil {
		ah.respErrors.SetTitle(err.Error())
		ah.response.SetResponse("", &ah.respErrors)

		return ah.response.Body(c, err)
	}

	err := ah.aUsecase.PostActivations(c, pa)

	if err != nil {
		ah.respErrors.SetTitle(err.Error())
		ah.response.SetResponse("", &ah.respErrors)

		return ah.response.Body(c, err)
	}

	ah.response.SetResponse("", &ah.respErrors)

	return ah.response.Body(c, err)
}
