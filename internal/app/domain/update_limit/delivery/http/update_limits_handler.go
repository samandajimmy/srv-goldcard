package http

import (
	"net/http"
	"srv-goldcard/internal/app/domain/update_limit"
	"srv-goldcard/internal/app/model"

	"github.com/labstack/echo"
)

type updateLimitHandler struct {
	response           model.Response
	respErrors         model.ResponseErrors
	updateLimitUseCase update_limit.UseCase
}

func NewUpdateLimitHandler(
	echoGroup model.EchoGroup,
	ulUseCase update_limit.UseCase) {
	handler := &updateLimitHandler{
		updateLimitUseCase: ulUseCase,
	}

	// Endpoint For Core
	echoGroup.API.POST("/update-limit/decreased-stl", handler.DecreasedSTL)
	echoGroup.API.POST("/update-limit/gte-payment", handler.CoreGtePayment)

	// Endpoint For PDS
	echoGroup.API.POST("/update-limit/increase/inquiry", handler.InquiryUpdateLimit)
	echoGroup.API.POST("/update-limit/increase", handler.PostUpdateLimit)
	echoGroup.API.GET("/update-limit/account-by-accnumber", handler.GetSavingAccount)
}

func (ul *updateLimitHandler) DecreasedSTL(c echo.Context) error {
	var pcds model.PayloadCoreDecreasedSTL
	ul.response, ul.respErrors = model.NewResponse()

	if err := c.Bind(&pcds); err != nil {
		ul.respErrors.SetTitle(model.MessageUnprocessableEntity)
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
	var piul model.PayloadInquiryUpdateLimit
	ul.response, ul.respErrors = model.NewResponse()

	if err := c.Bind(&piul); err != nil {
		ul.respErrors.SetTitle(model.MessageUnprocessableEntity)
		ul.response.SetResponse("", &ul.respErrors)

		return ul.response.Body(c, err)
	}

	if err := c.Validate(piul); err != nil {
		ul.respErrors.SetTitle(err.Error())
		ul.response.SetResponse("", &ul.respErrors)

		return ul.response.Body(c, err)
	}

	resp, err := ul.updateLimitUseCase.InquiryUpdateLimit(c, piul)

	if err.Title != "" {
		ul.response.SetResponse(resp, &err)

		return ul.response.Body(c, nil)
	}

	ul.response.SetResponse(resp, &err)

	return ul.response.Body(c, nil)
}

func (ul *updateLimitHandler) CoreGtePayment(c echo.Context) error {
	var pcgp model.PayloadCoreGtePayment
	ul.response, ul.respErrors = model.NewResponse()

	if err := c.Bind(&pcgp); err != nil {
		ul.respErrors.SetTitle(model.MessageUnprocessableEntity)
		ul.response.SetResponse("", &ul.respErrors)

		return ul.response.Body(c, err)
	}

	if err := c.Validate(pcgp); err != nil {
		ul.respErrors.SetTitle(err.Error())
		ul.response.SetResponse("", &ul.respErrors)

		return ul.response.Body(c, err)
	}

	err := ul.updateLimitUseCase.CoreGtePayment(c, pcgp)

	if err.Title != "" {
		ul.response.SetResponse("", &err)

		return ul.response.Body(c, nil)
	}

	ul.response.SetResponse("", &err)
	return ul.response.Body(c, nil)
}

func (ul *updateLimitHandler) PostUpdateLimit(c echo.Context) error {
	var pul model.PayloadUpdateLimit
	ul.response, ul.respErrors = model.NewResponse()

	if err := c.Bind(&pul); err != nil {
		ul.respErrors.SetTitle(model.MessageUnprocessableEntity)
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

func (ul *updateLimitHandler) GetSavingAccount(c echo.Context) error {
	var pan model.PayloadAccNumber
	ul.response, ul.respErrors = model.NewResponse()

	if err := c.Bind(&pan); err != nil {
		ul.respErrors.SetTitle(model.MessageUnprocessableEntity)
		ul.response.SetResponse("", &ul.respErrors)

		return ul.response.Body(c, err)
	}

	if err := c.Validate(pan); err != nil {
		ul.respErrors.SetTitle(err.Error())
		ul.response.SetResponse("", &ul.respErrors)

		return ul.response.Body(c, err)
	}

	responseData, err := ul.updateLimitUseCase.GetSavingAccount(c, pan)

	if err != nil {
		ul.respErrors.SetTitle(err.Error())
		ul.response.SetResponse("", &ul.respErrors)

		return c.JSON(http.StatusBadRequest, ul.response)
	}

	ul.response.SetResponse(responseData, &ul.respErrors)

	return ul.response.Body(c, nil)
}
