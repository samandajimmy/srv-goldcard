package http

import (
	"gade/srv-goldcard/models"
	"gade/srv-goldcard/update_limits"

	"github.com/labstack/echo"
)

type updateLimitHandler struct {
	response           models.Response
	respErrors         models.ResponseErrors
	updateLimitUseCase update_limits.UseCase
}

func NewUpdateLimitHandler(
	echoGroup models.EchoGroup,
	ulUseCase update_limits.UseCase) {
	handler := &updateLimitHandler{
		updateLimitUseCase: ulUseCase,
	}

	// Endpoint For Core
	echoGroup.API.POST("/update-limit/decreased-stl", handler.DecreasedSTL)

	// Endpoint For PDS
	echoGroup.API.POST("/update-limit/increase/inquiry", handler.InquiryUpdateLimit)
	echoGroup.API.POST("/update-limit/increase", handler.PostUpdateLimit)
}

func (ul *updateLimitHandler) DecreasedSTL(c echo.Context) error {
	var pcds models.PayloadCoreDecreasedSTL
	ul.response, ul.respErrors = models.NewResponse()

	if err := c.Bind(&pcds); err != nil {
		ul.respErrors.SetTitle(models.MessageUnprocessableEntity)
		ul.response.SetResponse("", &ul.respErrors)

		return ul.response.Body(c, err)
	}

	if err := c.Validate(pcds); err != nil {
		ul.respErrors.SetTitle(err.Error())
		ul.response.SetResponse("", &ul.respErrors)

		return ul.response.Body(c, err)
	}

	err := ul.updateLimitUseCase.DecreasedSTL(c, pcds)

	if err.Title != "" {
		ul.response.SetResponse("", &err)

		return ul.response.Body(c, nil)
	}

	ul.response.SetResponse("", &err)
	return ul.response.Body(c, nil)
}

func (ul *updateLimitHandler) InquiryUpdateLimit(c echo.Context) error {
	var piul models.PayloadInquiryUpdateLimit
	ul.response, ul.respErrors = models.NewResponse()

	if err := c.Bind(&piul); err != nil {
		ul.respErrors.SetTitle(models.MessageUnprocessableEntity)
		ul.response.SetResponse("", &ul.respErrors)

		return ul.response.Body(c, err)
	}

	if err := c.Validate(piul); err != nil {
		ul.respErrors.SetTitle(err.Error())
		ul.response.SetResponse("", &ul.respErrors)

		return ul.response.Body(c, err)
	}

	err := ul.updateLimitUseCase.InquiryUpdateLimit(c, piul)

	if err.Title != "" {
		ul.response.SetResponse("", &err)

		return ul.response.Body(c, nil)
	}

	ul.response.SetResponse("", &err)
	return ul.response.Body(c, nil)
}

func (ul *updateLimitHandler) PostUpdateLimit(c echo.Context) error {
	var pul models.PayloadUpdateLimit
	ul.response, ul.respErrors = models.NewResponse()

	if err := c.Bind(&pul); err != nil {
		ul.respErrors.SetTitle(models.MessageUnprocessableEntity)
		ul.response.SetResponse("", &ul.respErrors)

		return ul.response.Body(c, err)
	}

	if err := c.Validate(pul); err != nil {
		ul.respErrors.SetTitle(err.Error())
		ul.response.SetResponse("", &ul.respErrors)

		return ul.response.Body(c, err)
	}

	err := ul.updateLimitUseCase.PostUpdateLimit(c, pul)

	if err.Title != "" {
		ul.response.SetResponse("", &err)

		return ul.response.Body(c, nil)
	}

	ul.response.SetResponse("", &err)
	return ul.response.Body(c, nil)
}
