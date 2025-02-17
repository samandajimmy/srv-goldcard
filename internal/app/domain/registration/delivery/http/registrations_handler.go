package http

import (
	"os"
	"srv-goldcard/internal/app/domain/registration"
	"srv-goldcard/internal/app/model"

	"github.com/labstack/echo"
)

// RegistrationsHandler represent the httphandler for registrations
type RegistrationsHandler struct {
	response             model.Response
	respErrors           model.ResponseErrors
	registrationsUseCase registration.UseCase
}

// NewRegistrationsHandler represent to registration gold card
func NewRegistrationsHandler(
	echoGroup model.EchoGroup,
	regUseCase registration.UseCase) {
	handler := &RegistrationsHandler{
		registrationsUseCase: regUseCase,
	}

	// End Point For External'
	if os.Getenv("WITH_REGISTRATION") == "true" {
		echoGroup.API.POST("/registrations/address", handler.PostAddress)
		echoGroup.API.POST("/registrations/saving-account", handler.PostSavingAccount)
		echoGroup.API.POST("/registrations/personal-informations", handler.personalInfomations)
		echoGroup.API.POST("/registrations/card-limit", handler.cardLimit)
		echoGroup.API.POST("/registrations/final", handler.final)
		echoGroup.API.POST("/registrations/occupation", handler.PostOccupation)
		echoGroup.API.POST("/registrations/scheduler/final", handler.schedulerFinal)
	}

	echoGroup.API.POST("/registrations", handler.Registrations)
	echoGroup.API.GET("/registrations/address", handler.GetAddress)
	echoGroup.API.POST("/registrations/application-status", handler.applicationStatus)
	echoGroup.API.POST("/registrations/reset", handler.ResetRegistration)
	echoGroup.API.POST("/registrations/force-deliver", handler.ForceDeliver)

}

// Registrations a handler to handle goldcard registrations
func (reg *RegistrationsHandler) Registrations(c echo.Context) error {
	var pr model.PayloadRegistration
	reg.response, reg.respErrors = model.NewResponse()

	if err := c.Bind(&pr); err != nil {
		reg.respErrors.SetTitle(model.MessageUnprocessableEntity)
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	if err := c.Validate(pr); err != nil {
		reg.respErrors.SetTitle(err.Error())
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	resp, err := reg.registrationsUseCase.PostRegistration(c, pr)

	if err != nil {
		reg.respErrors.SetTitle(err.Error())
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	reg.response.SetResponse(resp, &reg.respErrors)

	return reg.response.Body(c, err)
}

// PostAddress a handler to update Address in table personal_informations
func (reg *RegistrationsHandler) PostAddress(c echo.Context) error {
	var plAddr model.PayloadAddress
	reg.response, reg.respErrors = model.NewResponse()

	if err := c.Bind(&plAddr); err != nil {
		reg.respErrors.SetTitle(model.MessageUnprocessableEntity)
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	if err := c.Validate(plAddr); err != nil {
		reg.respErrors.SetTitle(err.Error())
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	err := reg.registrationsUseCase.PostAddress(c, plAddr)

	if err != nil {
		reg.respErrors.SetTitle(err.Error())
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	reg.response.SetResponse("", &reg.respErrors)

	return reg.response.Body(c, err)
}

// GetAddress a handler to get Address in table personal_informations
func (reg *RegistrationsHandler) GetAddress(c echo.Context) error {
	var plApp model.PayloadAppNumber
	reg.response, reg.respErrors = model.NewResponse()

	if err := c.Bind(&plApp); err != nil {
		reg.respErrors.SetTitle(model.MessageUnprocessableEntity)
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	if err := c.Validate(plApp); err != nil {
		reg.respErrors.SetTitle(err.Error())
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	responseData, err := reg.registrationsUseCase.GetAddress(c, plApp)

	if err != nil {
		reg.respErrors.SetTitle(err.Error())
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	reg.response.SetResponse(responseData, &reg.respErrors)

	return reg.response.Body(c, err)
}

// PostSavingAccount a handler to update Saving Account in table applications
func (reg *RegistrationsHandler) PostSavingAccount(c echo.Context) error {
	var pl model.PayloadSavingAccount
	reg.response, reg.respErrors = model.NewResponse()

	if err := c.Bind(&pl); err != nil {
		reg.respErrors.SetTitle(model.MessageUnprocessableEntity)
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	if err := c.Validate(pl); err != nil {
		reg.respErrors.SetTitle(err.Error())
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	err := reg.registrationsUseCase.PostSavingAccount(c, pl)

	if err != nil {
		reg.respErrors.SetTitle(err.Error())
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	reg.response.SetResponse("", &reg.respErrors)
	return reg.response.Body(c, err)
}

func (reg *RegistrationsHandler) personalInfomations(c echo.Context) error {
	var ppi model.PayloadPersonalInformation
	reg.response, reg.respErrors = model.NewResponse()

	if err := c.Bind(&ppi); err != nil {
		reg.respErrors.SetTitle(model.MessageUnprocessableEntity)
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	if err := c.Validate(ppi); err != nil {
		reg.respErrors.SetTitle(err.Error())
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	err := reg.registrationsUseCase.PostPersonalInfo(c, ppi)

	if err != nil {
		reg.respErrors.SetTitle(err.Error())
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	reg.response.SetResponse("", &reg.respErrors)

	return reg.response.Body(c, err)
}

func (reg *RegistrationsHandler) cardLimit(c echo.Context) error {
	var pl model.PayloadCardLimit
	reg.response, reg.respErrors = model.NewResponse()

	if err := c.Bind(&pl); err != nil {
		reg.respErrors.SetTitle(model.MessageUnprocessableEntity)
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	if err := c.Validate(pl); err != nil {
		reg.respErrors.SetTitle(err.Error())
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	err := reg.registrationsUseCase.PostCardLimit(c, pl)

	if err != nil {
		reg.respErrors.SetTitle(err.Error())
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	reg.response.SetResponse("", &reg.respErrors)

	return reg.response.Body(c, err)
}

func (reg *RegistrationsHandler) schedulerFinal(c echo.Context) error {
	var pl model.PayloadAppNumber
	reg.response, reg.respErrors = model.NewResponse()

	if err := c.Bind(&pl); err != nil {
		reg.respErrors.SetTitle(model.MessageUnprocessableEntity)
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	if err := c.Validate(pl); err != nil {
		reg.respErrors.SetTitle(err.Error())
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	err := reg.registrationsUseCase.FinalRegistrationScheduler(c, pl)

	if err != nil {
		reg.respErrors.SetTitle(err.Error())
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	reg.response.SetResponse("", &reg.respErrors)

	return reg.response.Body(c, err)
}

func (reg *RegistrationsHandler) final(c echo.Context) error {
	var pl model.PayloadAppNumber
	reg.response, reg.respErrors = model.NewResponse()

	if err := c.Bind(&pl); err != nil {
		reg.respErrors.SetTitle(model.MessageUnprocessableEntity)
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	if err := c.Validate(pl); err != nil {
		reg.respErrors.SetTitle(err.Error())
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	err := reg.registrationsUseCase.FinalRegistrationPdsApi(c, pl)

	if err != nil {
		reg.respErrors.SetTitle(err.Error())
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	reg.response.SetResponse("", &reg.respErrors)

	return reg.response.Body(c, err)
}

func (reg *RegistrationsHandler) applicationStatus(c echo.Context) error {
	var pl model.PayloadAppNumber
	reg.response, reg.respErrors = model.NewResponse()

	if err := c.Bind(&pl); err != nil {
		reg.respErrors.SetTitle(model.MessageUnprocessableEntity)
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	if err := c.Validate(pl); err != nil {
		reg.respErrors.SetTitle(err.Error())
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	resp, err := reg.registrationsUseCase.GetAppStatus(c, pl)

	if err != nil {
		reg.respErrors.SetTitle(err.Error())
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	reg.response.SetResponse(resp, &reg.respErrors)

	return reg.response.Body(c, err)
}

// PostOccupation a handler to update occipation in table occupations
func (reg *RegistrationsHandler) PostOccupation(c echo.Context) error {
	var pl model.PayloadOccupation
	reg.response, reg.respErrors = model.NewResponse()

	if err := c.Bind(&pl); err != nil {
		reg.respErrors.SetTitle(model.MessageUnprocessableEntity)
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	if err := c.Validate(pl); err != nil {
		reg.respErrors.SetTitle(err.Error())
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	err := reg.registrationsUseCase.PostOccupation(c, pl)

	if err != nil {
		reg.respErrors.SetTitle(err.Error())
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	reg.response.SetResponse("", &reg.respErrors)
	return reg.response.Body(c, err)
}

// ResetRegistration a handler to reset registration
func (reg *RegistrationsHandler) ResetRegistration(c echo.Context) error {
	var pl model.PayloadAppNumber
	reg.response, reg.respErrors = model.NewResponse()

	if err := c.Bind(&pl); err != nil {
		reg.respErrors.SetTitle(model.MessageUnprocessableEntity)
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	if err := c.Validate(pl); err != nil {
		reg.respErrors.SetTitle(err.Error())
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	err := reg.registrationsUseCase.ResetRegistration(c, pl)

	if err != nil {
		reg.respErrors.SetTitle(err.Error())
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	reg.response.SetResponse("", &reg.respErrors)
	return reg.response.Body(c, err)
}

// ResetRegistration a handler to reset registration
func (reg *RegistrationsHandler) ForceDeliver(c echo.Context) error {
	var pl model.PayloadAppNumber
	reg.response, reg.respErrors = model.NewResponse()

	if err := c.Bind(&pl); err != nil {
		reg.respErrors.SetTitle(model.MessageUnprocessableEntity)
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	if err := c.Validate(pl); err != nil {
		reg.respErrors.SetTitle(err.Error())
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	err := reg.registrationsUseCase.ForceDeliver(c, pl)

	if err != nil {
		reg.respErrors.SetTitle(err.Error())
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	reg.response.SetResponse("", &reg.respErrors)
	return reg.response.Body(c, err)
}
