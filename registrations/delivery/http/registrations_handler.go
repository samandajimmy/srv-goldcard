package http

import (
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/registrations"

	"github.com/labstack/echo"
)

var response models.Response

// RegistrationsHandler represent the httphandler for registrations
type RegistrationsHandler struct {
	response             models.Response
	respErrors           models.ResponseErrors
	registrationsUseCase registrations.UseCase
}

// NewRegistrationsHandler represent to registration gold card
func NewRegistrationsHandler(
	echoGroup models.EchoGroup,
	regUseCase registrations.UseCase) {
	handler := &RegistrationsHandler{
		registrationsUseCase: regUseCase,
	}

	// End Point For External
	echoGroup.API.POST("/registrations", handler.Registrations)
	echoGroup.API.POST("/registrations/address", handler.PostAddress)
	echoGroup.API.POST("/registrations/saving-account", handler.PostSavingAccount)
	echoGroup.API.POST("/registrations/personal-informations", handler.personalInfomations)
	echoGroup.API.POST("/registrations/card-limit", handler.cardLimit)
	echoGroup.API.POST("/registrations/final", handler.final)
}

// Registrations a handler to create a campaign
func (reg *RegistrationsHandler) Registrations(c echo.Context) error {
	var pr models.PayloadRegistration
	reg.response, reg.respErrors = models.NewResponse()

	if err := c.Bind(&pr); err != nil {
		reg.respErrors.SetTitle(models.MessageUnprocessableEntity)
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	if err := c.Validate(pr); err != nil {
		reg.respErrors.SetTitle(err.Error())
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	appNumber, err := reg.registrationsUseCase.PostRegistration(c, pr)

	if err != nil {
		reg.respErrors.SetTitle(err.Error())
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	reg.response.SetResponse(map[string]string{"applicationNumber": appNumber}, &reg.respErrors)

	return reg.response.Body(c, err)
}

// PostAddress a handler to update Address in table personal_informations
func (reg *RegistrationsHandler) PostAddress(c echo.Context) error {
	var plAddr models.PayloadAddress
	reg.response, reg.respErrors = models.NewResponse()

	if err := c.Bind(&plAddr); err != nil {
		reg.respErrors.SetTitle(models.MessageUnprocessableEntity)
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

// PostSavingAccount a handler to update Saving Account in table applications
func (reg *RegistrationsHandler) PostSavingAccount(c echo.Context) error {
	var pl models.PayloadSavingAccount
	reg.response, reg.respErrors = models.NewResponse()

	if err := c.Bind(&pl); err != nil {
		reg.respErrors.SetTitle(models.MessageUnprocessableEntity)
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
	var ppi models.PayloadPersonalInformation
	reg.response, reg.respErrors = models.NewResponse()

	if err := c.Bind(&ppi); err != nil {
		reg.respErrors.SetTitle(models.MessageUnprocessableEntity)
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
	var pl models.PayloadCardLimit
	reg.response, reg.respErrors = models.NewResponse()

	if err := c.Bind(&pl); err != nil {
		reg.respErrors.SetTitle(models.MessageUnprocessableEntity)
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

func (reg *RegistrationsHandler) final(c echo.Context) error {
	var pl models.PayloadRegistrationFinal
	reg.response, reg.respErrors = models.NewResponse()

	if err := c.Bind(&pl); err != nil {
		reg.respErrors.SetTitle(models.MessageUnprocessableEntity)
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	if err := c.Validate(pl); err != nil {
		reg.respErrors.SetTitle(err.Error())
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	err := reg.registrationsUseCase.FinalRegistration(c, pl)

	if err != nil {
		reg.respErrors.SetTitle(err.Error())
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	reg.response.SetResponse("", &reg.respErrors)

	return reg.response.Body(c, err)
}
