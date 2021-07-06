package http

import (
	"gade/srv-goldcard/api"
	"gade/srv-goldcard/models"

	"github.com/labstack/echo"
)

type HealthsHandler struct {
	response   models.Response
	respErrors models.ResponseErrors
}

func NewHealthsHandler(ech *echo.Echo) {
	handler := &HealthsHandler{}

	ech.GET("/health-check", handler.healthCheck)
}

func (health *HealthsHandler) healthCheck(c echo.Context) error {
	var err error
	response := models.RespHealthCheck{}
	_, err = api.NewSwitchingAPI(c)
	response.SwitchingApi = (err == nil)

	_, err = api.NewBriAPI(c)
	response.BriApi = (err == nil)

	err = api.PdsHealthCheck(c)
	response.PdsApi = (err == nil)

	health.response.SetResponse(response, &health.respErrors)
	return health.response.Body(c, err)
}
