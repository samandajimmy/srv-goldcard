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
	// ? this was only for a template and you should replace `_`
	// ? with handler variable name
	_ = &ActivationsHandler{
		activationsUseCase: aUseCase,
	}

	// End Point For External
	// ? this was only for a template and you should replace `func(echo.Context) error { return nil }`
	// ? with handler func
	echoGroup.API.POST("/activations", func(echo.Context) error { return nil })
}
