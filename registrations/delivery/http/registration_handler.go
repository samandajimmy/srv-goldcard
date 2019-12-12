package http

import (
	"gade/srv-goldcard/models"
	"net/http"

	"github.com/labstack/echo"
)

var response models.Response

// RegistrationsHandler represent the httphandler for registrations
type RegistrationsHandler struct {
}

// NewRegistrationHandler represent to registration gold card
func NewRegistrationHandler(echoGroup models.EchoGroup) {
	handler := &RegistrationsHandler{}

	// End Point For External
	echoGroup.API.POST("/registrations", handler.Registrations)
}

// Registrations a handler to create a campaign
func (reg *RegistrationsHandler) Registrations(c echo.Context) error {
	var registration models.Registrations

	response = models.Response{}
	err := c.Bind(&registration)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(http.StatusOK, response)
	}

	apiRequest, err := models.NewClientRequest("https://apidigitaldev.pegadaian.co.id/v2", "application/x-www-form-urlencoded")

	apiRequest.ApiRequest(c, "/profile/testing_go", "POST", registration, &response)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(http.StatusOK, response)
	}
	return c.JSON(http.StatusOK, response)
}
