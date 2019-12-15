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

// NewRegistrationsHandler represent to registration gold card
func NewRegistrationsHandler(echoGroup models.EchoGroup) {
	handler := &RegistrationsHandler{}

	// End Point For External
	echoGroup.API.POST("/registrations", handler.Registrations)
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
