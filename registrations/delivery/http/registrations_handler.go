package http

import (
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/registrations"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

var response models.Response

// RegistrationsHandler represent the httphandler for registrations
type RegistrationsHandler struct {
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
	echoGroup.API.POST("/registrations/address/simpan", handler.PostAddress)
	echoGroup.API.GET("/registrations/address", handler.GetAddress)
}

// Registrations a handler to create a campaign
func (reg *RegistrationsHandler) Registrations(c echo.Context) error {
	var registrations models.Registrations

	response = models.Response{}
	err := c.Bind(&registrations)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(http.StatusOK, response)
	}

	api, _ := models.NewAPI("https://apidigitaldev.pegadaian.co.id/v2", "application/x-www-form-urlencoded")
	req, _ := api.Request("/profile/testing_go", "POST", registrations)
	_, _ = api.Do(req, &response)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(http.StatusOK, response)
	}

	return c.JSON(http.StatusOK, response)
}

// PostAddress a handler to update Address in table personal_informations
func (reg *RegistrationsHandler) PostAddress(c echo.Context) error {
	respErrors := &models.ResponseErrors{}
	logger := models.RequestLogger{}
	response = models.Response{}
	var registrations models.Registrations

	c.Bind(&registrations)
	logger.DataLog(c, registrations).Info("Start of Post Address")
	err := reg.registrationsUseCase.PostAddress(c, &registrations)

	if err != nil {
		respErrors.SetTitle(err.Error())
		response.SetResponse("", respErrors)
		logger.DataLog(c, response).Info("End of Post Address")

		return c.JSON(getStatusCode(err), response)
	}

	response.Code = "00"
	response.Status = models.StatusSuccess
	response.Message = models.MessageUpdateSuccess
	logger.DataLog(c, response).Info("End of Post Address")
	return c.JSON(getStatusCode(err), response)
}

// GetAddress a handler to get Address in table personal_informations
func (reg *RegistrationsHandler) GetAddress(c echo.Context) error {
	respErrors := &models.ResponseErrors{}
	logger := models.RequestLogger{}
	response = models.Response{}

	logger.DataLog(c, c.QueryParam("phoneno")).Info("Start of Get Address")
	res, err := reg.registrationsUseCase.GetAddress(c, c.QueryParam("phoneno"))

	if err != nil {
		respErrors.SetTitle(err.Error())
		response.SetResponse("", respErrors)
		logger.DataLog(c, response).Info("End of Get Address")

		return c.JSON(getStatusCode(err), res)
	}

	logger.DataLog(c, response).Info("End of Post Address")

	response.SetResponse(res, respErrors)
	return c.JSON(getStatusCode(err), response)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	if strings.Contains(err.Error(), "400") {
		return http.StatusBadRequest
	}

	switch err {
	case models.ErrInternalServerError:
		return http.StatusInternalServerError
	case models.ErrNotFound:
		return http.StatusNotFound
	case models.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusOK
	}
}
