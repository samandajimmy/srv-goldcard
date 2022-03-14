package http

import (
	activation "srv-goldcard/internal/app/domain/activation"
	"srv-goldcard/internal/app/model"

	"github.com/labstack/echo"
)

// ActivationsHandler represent the httphandler for activation
type ActivationsHandler struct {
	response   model.Response
	respErrors model.ResponseErrors
	aUsecase   activation.UseCase
}

// NewActivationsHandler represent to activation gold card
func NewActivationsHandler(echoGroup model.EchoGroup, aUseCase activation.UseCase) {
	handler := &ActivationsHandler{aUsecase: aUseCase}

	// End Point For External
	echoGroup.API.POST("/activations", handler.Activations)
	echoGroup.API.POST("/re-activations", handler.Reactivations)
	echoGroup.API.POST("/activations/inquiry", handler.ActivationsInquiry)
	echoGroup.API.POST("/activations/force", handler.ActivationsForce)
}

// ActivationsInquiry a handler to handle goldcard inquiry activation
func (ah *ActivationsHandler) ActivationsInquiry(c echo.Context) error {
	var pl model.PayloadAppNumber
	ah.response, ah.respErrors = model.NewResponse()

	if err := c.Bind(&pl); err != nil {
		ah.respErrors.SetTitle(model.MessageUnprocessableEntity)
		ah.response.SetResponse("", &ah.respErrors)

		return ah.response.Body(c, err)
	}

	if err := c.Validate(pl); err != nil {
		ah.respErrors.SetTitle(err.Error())
		ah.response.SetResponse("", &ah.respErrors)

		return ah.response.Body(c, err)
	}

	acc := model.Account{}
	acc.Application.ApplicationNumber = pl.ApplicationNumber

	resp, err := ah.aUsecase.InquiryActivation(c, acc)

	if err.Title != "" {
		ah.response.SetResponse(resp, &err)

		return ah.response.Body(c, nil)
	}

	ah.response.SetResponse("", &err)
	return ah.response.Body(c, nil)
}

func (ah *ActivationsHandler) ActivationsForce(c echo.Context) error {
	var pl model.PayloadAppNumber

	if err := c.Bind(&pl); err != nil {
		ah.respErrors.SetTitle(model.MessageUnprocessableEntity)
		ah.response.SetResponse("", &ah.respErrors)

		return ah.response.Body(c, err)
	}

	if err := c.Validate(pl); err != nil {
		ah.respErrors.SetTitle(err.Error())
		ah.response.SetResponse("", &ah.respErrors)

		return ah.response.Body(c, err)
	}

	acc := model.Account{}
	acc.Application.ApplicationNumber = pl.ApplicationNumber

	resp, err := ah.aUsecase.ForceActivation(c, acc)

	if err != nil {
		ah.respErrors.SetTitle(err.Error())
		ah.response.SetResponse("", &ah.respErrors)

		return ah.response.Body(c, err)
	}

	ah.response.SetResponse(resp, &ah.respErrors)

	return ah.response.Body(c, err)
}

// Activations a handler to activation
func (ah *ActivationsHandler) Activations(c echo.Context) error {
	var pa model.PayloadActivations
	ah.response, ah.respErrors = model.NewResponse()

	if err := c.Bind(&pa); err != nil {
		ah.respErrors.SetTitle(model.MessageUnprocessableEntity)
		ah.response.SetResponse("", &ah.respErrors)

		return ah.response.Body(c, err)
	}

	if err := c.Validate(pa); err != nil {
		ah.respErrors.SetTitle(err.Error())
		ah.response.SetResponse("", &ah.respErrors)

		return ah.response.Body(c, err)
	}

	if err := ah.aUsecase.ValidateActivation(c, pa); err.Title != "" {
		ah.response.SetResponse("", &err)

		return ah.response.Body(c, nil)
	}

	resp, err := ah.aUsecase.PostActivations(c, pa)

	if err != nil {
		ah.respErrors.SetTitle(err.Error())
		ah.response.SetResponse("", &ah.respErrors)

		return ah.response.Body(c, err)
	}

	ah.response.SetResponse(resp, &ah.respErrors)

	return ah.response.Body(c, err)
}
func (ah *ActivationsHandler) Reactivations(c echo.Context) error {
	var pa model.PayloadActivations
	ah.response, ah.respErrors = model.NewResponse()

	if err := c.Bind(&pa); err != nil {
		ah.respErrors.SetTitle(model.MessageUnprocessableEntity)
		ah.response.SetResponse("", &ah.respErrors)

		return ah.response.Body(c, err)
	}

	if err := c.Validate(pa); err != nil {
		ah.respErrors.SetTitle(err.Error())
		ah.response.SetResponse("", &ah.respErrors)

		return ah.response.Body(c, err)
	}

	if err := ah.aUsecase.ValidateActivation(c, pa); err.Title != "" {
		ah.response.SetResponse("", &err)

		return ah.response.Body(c, nil)
	}

	resp, err := ah.aUsecase.PostReactivations(c, pa)

	if err != nil {
		ah.respErrors.SetTitle(err.Error())
		ah.response.SetResponse("", &ah.respErrors)

		return ah.response.Body(c, err)
	}

	ah.response.SetResponse(resp, &ah.respErrors)

	return ah.response.Body(c, err)
}
