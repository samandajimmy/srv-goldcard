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
	echoGroup.API.GET("/registrations/address", handler.GetAddress)
	echoGroup.API.POST("/registrations/saving-account", handler.PostSavingAccount)
	echoGroup.API.POST("/registrations/personal-informations", handler.personalInfomations)
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
	var registrations models.Registrations
	reg.response, reg.respErrors = models.NewResponse()

	if err := c.Bind(&registrations); err != nil {
		reg.respErrors.SetTitle(models.MessageUnprocessableEntity)
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	if err := c.Validate(registrations); err != nil {
		reg.respErrors.SetTitle(err.Error())
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	err := reg.registrationsUseCase.PostAddress(c, &registrations)

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
	var getAddress models.PayloadGetAddress
	reg.response, reg.respErrors = models.NewResponse()

	if err := c.Bind(&getAddress); err != nil {
		reg.respErrors.SetTitle(models.MessageUnprocessableEntity)
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	if err := c.Validate(getAddress); err != nil {
		reg.respErrors.SetTitle(err.Error())
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	res, err := reg.registrationsUseCase.GetAddress(c, getAddress.PhoneNumber)

	if err != nil {
		reg.respErrors.SetTitle(err.Error())
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	reg.response.SetResponse(res, &reg.respErrors)
	return reg.response.Body(c, err)
}

// PostSavingAccount a handler to update Saving Account in table applications
func (reg *RegistrationsHandler) PostSavingAccount(c echo.Context) error {
	var applications models.Applications
	reg.response, reg.respErrors = models.NewResponse()

	if err := c.Bind(&applications); err != nil {
		reg.respErrors.SetTitle(models.MessageUnprocessableEntity)
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	if err := c.Validate(applications); err != nil {
		reg.respErrors.SetTitle(err.Error())
		reg.response.SetResponse("", &reg.respErrors)

		return reg.response.Body(c, err)
	}

	err := reg.registrationsUseCase.PostSavingAccount(c, &applications)

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
