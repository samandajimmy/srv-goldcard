package http

import (
	"fmt"
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/registration"
	"net/http"

	"github.com/labstack/echo"
)

var response models.Response

type RegistrationHandler struct {
	RegistrationUseCase registration.UseCase
}

func NewRegistrationHandler(echoGroup models.EchoGroup) {
	handler := &RegistrationHandler{}

	// End Point For External
	echoGroup.API.POST("/registration", handler.TestRegistration)
}

// TestRegistration a handler to create a campaign
func (reg *RegistrationHandler) TestRegistration(c echo.Context) error {
	var registration models.Registration

	response = models.Response{}
	// logger := models.RequestLogger{}
	err := c.Bind(&registration)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		// return c.JSON(getStatusCode(err), response)
		return c.JSON(http.StatusOK, response)
	}
	fmt.Println("test3")
	fmt.Println(registration)
	fmt.Println("test3")

	apiRequest, err := models.NewClientRequest("https://apidigitaldev.pegadaian.co.id/v2", "application/json")

	apiRequest.ApiRequest(c, "/profile/testing_go", "POST", registration, &response)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		// return c.JSON(getStatusCode(err), response)
		return c.JSON(http.StatusOK, response)
	}

	// response.Status = models.StatusSuccess
	// response.Message = models.MessageSaveSuccess

	// requestLogger.Info("End of create a campaign.")
	fmt.Println("test2")
	fmt.Println(response)
	fmt.Println("test2")
	return c.JSON(http.StatusOK, response)
}
