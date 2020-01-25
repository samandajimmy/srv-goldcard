package http

import (
	"gade/srv-goldcard/activations"
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

// ActivationsHandler represent the httphandler for activations
type ActivationsHandler struct {
	response           models.Response
	respErrors         models.ResponseErrors
	activationsUseCase activations.UseCase
}

// NewActivationsHandler represent to activations gold card
func NewActivationsHandler(echoGroup models.EchoGroup, aUseCase activations.UseCase) {
	handler := &ActivationsHandler{
		activationsUseCase: aUseCase,
	}

	// End Point For External
	echoGroup.API.POST("/activations", handler.Activations)
}

// Activations a handler to activations
func (act *ActivationsHandler) Activations(c echo.Context) error {
	var pa models.PayloadActivations
	act.response, act.respErrors = models.NewResponse()

	if err := c.Bind(&pa); err != nil {
		act.respErrors.SetTitle(models.MessageUnprocessableEntity)
		act.response.SetResponse("", &act.respErrors)

		return act.response.Body(c, err)
	}

	if err := c.Validate(pa); err != nil {
		act.respErrors.SetTitle(err.Error())
		act.response.SetResponse("", &act.respErrors)

		return act.response.Body(c, err)
	}

	err := act.activationsUseCase.PostActivations(c, pa)

	if err != nil {
		act.respErrors.SetTitle(err.Error())
		act.response.SetResponse("", &act.respErrors)

		return act.response.Body(c, err)
	}

	act.response.SetResponse("", &act.respErrors)

	return act.response.Body(c, err)
}
